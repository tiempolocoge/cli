package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/camelcase"
	"github.com/github/ghcs/api"
	"github.com/github/ghcs/cmd/ghcs/output"
	"github.com/github/ghcs/internal/codespaces"
	"github.com/spf13/cobra"
)

type createOptions struct {
	repo       string
	branch     string
	machine    string
	showStatus bool
}

func newCreateCmd() *cobra.Command {
	opts := &createOptions{}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a Codespace",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(opts)
		},
	}

	createCmd.Flags().StringVarP(&opts.repo, "repo", "r", "", "repository name with owner: user/repo")
	createCmd.Flags().StringVarP(&opts.branch, "branch", "b", "", "repository branch")
	createCmd.Flags().StringVarP(&opts.machine, "machine", "m", "", "hardware specifications for the VM")
	createCmd.Flags().BoolVarP(&opts.showStatus, "status", "s", false, "show status of post-create command and dotfiles")

	return createCmd
}

func init() {
	rootCmd.AddCommand(newCreateCmd())
}

func create(opts *createOptions) error {
	ctx := context.Background()
	apiClient := api.New(os.Getenv("GITHUB_TOKEN"))
	locationCh := getLocation(ctx, apiClient)
	userCh := getUser(ctx, apiClient)
	log := output.NewLogger(os.Stdout, os.Stderr, false)

	repo, err := getRepoName(opts.repo)
	if err != nil {
		return fmt.Errorf("error getting repository name: %v", err)
	}
	branch, err := getBranchName(opts.branch)
	if err != nil {
		return fmt.Errorf("error getting branch name: %v", err)
	}

	repository, err := apiClient.GetRepository(ctx, repo)
	if err != nil {
		return fmt.Errorf("error getting repository: %v", err)
	}

	locationResult := <-locationCh
	if locationResult.Err != nil {
		return fmt.Errorf("error getting codespace region location: %v", locationResult.Err)
	}

	userResult := <-userCh
	if userResult.Err != nil {
		return fmt.Errorf("error getting codespace user: %v", userResult.Err)
	}

	machine, err := getMachineName(ctx, opts.machine, userResult.User, repository, locationResult.Location, apiClient)
	if err != nil {
		return fmt.Errorf("error getting machine type: %v", err)
	}
	if machine == "" {
		return errors.New("There are no available machine types for this repository")
	}

	log.Println("Creating your codespace...")

	codespace, err := apiClient.CreateCodespace(ctx, userResult.User, repository, machine, branch, locationResult.Location)
	if err != nil {
		return fmt.Errorf("error creating codespace: %v", err)
	}

	if opts.showStatus {
		if err := showStatus(ctx, log, apiClient, userResult.User, codespace); err != nil {
			return fmt.Errorf("show status: %w", err)
		}
	}

	log.Printf("Codespace created: ")

	fmt.Fprintln(os.Stdout, codespace.Name)

	return nil
}

func showStatus(ctx context.Context, log *output.Logger, apiClient *api.API, user *api.User, codespace *api.Codespace) error {
	var lastState codespaces.PostCreateState
	var breakNextState bool

	finishedStates := make(map[string]bool)
	ctx, stopPolling := context.WithCancel(ctx)

	poller := func(states []codespaces.PostCreateState) {
		var inProgress bool
		for _, state := range states {
			if _, found := finishedStates[state.Name]; found {
				continue // skip this state as we've processed it already
			}

			if state.Name != lastState.Name {
				log.Print(state.Name)

				if state.Status == codespaces.PostCreateStateRunning {
					inProgress = true
					lastState = state
					log.Print("...")
					break
				}

				finishedStates[state.Name] = true
				log.Println("..." + state.Status)
			} else {
				if state.Status == codespaces.PostCreateStateRunning {
					inProgress = true
					log.Print(".")
					break
				}

				finishedStates[state.Name] = true
				log.Println(state.Status)
				lastState = codespaces.PostCreateState{} // reset the value
			}
		}

		if !inProgress {
			if breakNextState {
				stopPolling()
				return
			}
			breakNextState = true
		}
	}

	if err := codespaces.PollPostCreateStates(ctx, log, apiClient, user, codespace, poller); err != nil {
		return fmt.Errorf("failed to poll state changes from codespace: %v", err)
	}

	return nil
}

type getUserResult struct {
	User *api.User
	Err  error
}

// getUser fetches the user record associated with the GITHUB_TOKEN
func getUser(ctx context.Context, apiClient *api.API) <-chan getUserResult {
	ch := make(chan getUserResult)
	go func() {
		user, err := apiClient.GetUser(ctx)
		ch <- getUserResult{user, err}
	}()
	return ch
}

type locationResult struct {
	Location string
	Err      error
}

// getLocation fetches the closest Codespace datacenter region/location to the user.
func getLocation(ctx context.Context, apiClient *api.API) <-chan locationResult {
	ch := make(chan locationResult)
	go func() {
		location, err := apiClient.GetCodespaceRegionLocation(ctx)
		ch <- locationResult{location, err}
	}()
	return ch
}

// getRepoName prompts the user for the name of the repository, or returns the repository if non-empty.
func getRepoName(repo string) (string, error) {
	if repo != "" {
		return repo, nil
	}

	repoSurvey := []*survey.Question{
		{
			Name:     "repository",
			Prompt:   &survey.Input{Message: "Repository"},
			Validate: survey.Required,
		},
	}
	err := survey.Ask(repoSurvey, &repo)
	return repo, err
}

// getBranchName prompts the user for the name of the branch, or returns the branch if non-empty.
func getBranchName(branch string) (string, error) {
	if branch != "" {
		return branch, nil
	}

	branchSurvey := []*survey.Question{
		{
			Name:     "branch",
			Prompt:   &survey.Input{Message: "Branch"},
			Validate: survey.Required,
		},
	}
	err := survey.Ask(branchSurvey, &branch)
	return branch, err
}

// getMachineName prompts the user to select the machine type, or validates the machine if non-empty.
func getMachineName(ctx context.Context, machine string, user *api.User, repo *api.Repository, location string, apiClient *api.API) (string, error) {
	skus, err := apiClient.GetCodespacesSkus(ctx, user, repo, location)
	if err != nil {
		return "", fmt.Errorf("error getting codespace skus: %v", err)
	}

	// if user supplied a machine type, it must be valid
	// if no machine type was supplied, we don't error if there are no machine types for the current repo
	if machine != "" {
		for _, sku := range skus {
			if machine == sku.Name {
				return machine, nil
			}
		}

		availableSkus := make([]string, len(skus))
		for i := 0; i < len(skus); i++ {
			availableSkus[i] = skus[i].Name
		}

		return "", fmt.Errorf("there is no such machine for the repository: %s\nAvailable machines: %v", machine, availableSkus)
	} else if len(skus) == 0 {
		return "", nil
	}

	skuNames := make([]string, 0, len(skus))
	skuByName := make(map[string]*api.Sku)
	for _, sku := range skus {
		nameParts := camelcase.Split(sku.Name)
		machineName := strings.Title(strings.ToLower(nameParts[0]))
		skuName := fmt.Sprintf("%s - %s", machineName, sku.DisplayName)
		skuNames = append(skuNames, skuName)
		skuByName[skuName] = sku
	}

	skuSurvey := []*survey.Question{
		{
			Name: "sku",
			Prompt: &survey.Select{
				Message: "Choose Machine Type:",
				Options: skuNames,
				Default: skuNames[0],
			},
			Validate: survey.Required,
		},
	}

	var skuAnswers struct{ SKU string }
	if err := survey.Ask(skuSurvey, &skuAnswers); err != nil {
		return "", fmt.Errorf("error getting SKU: %v", err)
	}

	sku := skuByName[skuAnswers.SKU]
	machine = sku.Name

	return machine, nil
}
