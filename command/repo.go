package command

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/api"
	"github.com/cli/cli/git"
	"github.com/cli/cli/internal/ghrepo"
	"github.com/cli/cli/internal/run"
	"github.com/cli/cli/utils"
	"github.com/spf13/cobra"
)

func init() {
	repoCmd.AddCommand(repoCreateCmd)
	repoCreateCmd.Flags().StringP("description", "d", "", "Description of repository")
	repoCreateCmd.Flags().StringP("homepage", "h", "", "Repository home page URL")
	repoCreateCmd.Flags().StringP("team", "t", "", "The name of the organization team to be granted access")
	repoCreateCmd.Flags().Bool("enable-issues", true, "Enable issues in the new repository")
	repoCreateCmd.Flags().Bool("enable-wiki", true, "Enable wiki in the new repository")
	repoCreateCmd.Flags().Bool("public", false, "Make the new repository public (default: private)")

	repoCmd.AddCommand(repoForkCmd)
	repoForkCmd.Flags().String("clone", "prompt", "Clone fork: {true|false|prompt}")
	repoForkCmd.Flags().String("remote", "prompt", "Add remote for fork: {true|false|prompt}")
	repoForkCmd.Flags().Lookup("clone").NoOptDefVal = "true"
	repoForkCmd.Flags().Lookup("remote").NoOptDefVal = "true"

	repoCmd.AddCommand(repoCreditsCmd)
	repoCreditsCmd.Flags().BoolP("static", "s", false, "Print a static version of the credits")
}

var repoCmd = &cobra.Command{
	Use:   "repo <command>",
	Short: "Create, clone, fork, and view repositories",
	Long:  `Work with GitHub repositories`,
	Example: heredoc.Doc(`
	$ gh repo create
	$ gh repo clone cli/cli
	$ gh repo view --web
	`),
	Annotations: map[string]string{
		"IsCore": "true",
		"help:arguments": `
A repository can be supplied as an argument in any of the following formats:
- "OWNER/REPO"
- by URL, e.g. "https://github.com/OWNER/REPO"`},
}

var repoCreateCmd = &cobra.Command{
	Use:   "create [<name>]",
	Short: "Create a new repository",
	Long:  `Create a new GitHub repository.`,
	Example: heredoc.Doc(`
	# create a repository under your account using the current directory name
	$ gh repo create

	# create a repository with a specific name
	$ gh repo create my-project

	# create a repository in an organization
	$ gh repo create cli/my-project
	`),
	Annotations: map[string]string{"help:arguments": `A repository can be supplied as an argument in any of the following formats:
- <OWNER/REPO>
- by URL, e.g. "https://github.com/OWNER/REPO"`},
	RunE: repoCreate,
}

var repoForkCmd = &cobra.Command{
	Use:   "fork [<repository>]",
	Short: "Create a fork of a repository",
	Long: `Create a fork of a repository.

With no argument, creates a fork of the current repository. Otherwise, forks the specified repository.`,
	RunE: repoFork,
}

var repoCreditsCmd = &cobra.Command{
	Use:   "credits [<repository>]",
	Short: "View credits for a repository",
	Example: heredoc.Doc(`
	# view credits for the current repository
	$ gh repo credits
	
	# view credits for a specific repository
	$ gh repo credits cool/repo

	# print a non-animated thank you
	$ gh repo credits -s
	
	# pipe to just print the contributors, one per line
	$ gh repo credits | cat
	`),
	Args:   cobra.MaximumNArgs(1),
	RunE:   repoCredits,
	Hidden: true,
}

func addUpstreamRemote(cmd *cobra.Command, parentRepo ghrepo.Interface, cloneDir string) error {
	upstreamURL := formatRemoteURL(cmd, parentRepo)

	cloneCmd := git.GitCommand("-C", cloneDir, "remote", "add", "-f", "upstream", upstreamURL)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	return run.PrepareCmd(cloneCmd).Run()
}

func repoCreate(cmd *cobra.Command, args []string) error {
	projectDir, projectDirErr := git.ToplevelDir()

	orgName := ""
	teamSlug, err := cmd.Flags().GetString("team")
	if err != nil {
		return err
	}

	var name string
	if len(args) > 0 {
		name = args[0]
		if strings.Contains(name, "/") {
			newRepo, err := ghrepo.FromFullName(name)
			if err != nil {
				return fmt.Errorf("argument error: %w", err)
			}
			orgName = newRepo.RepoOwner()
			name = newRepo.RepoName()
		}
	} else {
		if projectDirErr != nil {
			return projectDirErr
		}
		name = path.Base(projectDir)
	}

	isPublic, err := cmd.Flags().GetBool("public")
	if err != nil {
		return err
	}
	hasIssuesEnabled, err := cmd.Flags().GetBool("enable-issues")
	if err != nil {
		return err
	}
	hasWikiEnabled, err := cmd.Flags().GetBool("enable-wiki")
	if err != nil {
		return err
	}
	description, err := cmd.Flags().GetString("description")
	if err != nil {
		return err
	}
	homepage, err := cmd.Flags().GetString("homepage")
	if err != nil {
		return err
	}

	// TODO: move this into constant within `api`
	visibility := "PRIVATE"
	if isPublic {
		visibility = "PUBLIC"
	}

	input := api.RepoCreateInput{
		Name:             name,
		Visibility:       visibility,
		OwnerID:          orgName,
		TeamID:           teamSlug,
		Description:      description,
		HomepageURL:      homepage,
		HasIssuesEnabled: hasIssuesEnabled,
		HasWikiEnabled:   hasWikiEnabled,
	}

	ctx := contextForCommand(cmd)
	client, err := apiClientForContext(ctx)
	if err != nil {
		return err
	}

	repo, err := api.RepoCreate(client, input)
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	greenCheck := utils.Green("✓")
	isTTY := false
	if outFile, isFile := out.(*os.File); isFile {
		isTTY = utils.IsTerminal(outFile)
		if isTTY {
			// FIXME: duplicates colorableOut
			out = utils.NewColorable(outFile)
		}
	}

	if isTTY {
		fmt.Fprintf(out, "%s Created repository %s on GitHub\n", greenCheck, ghrepo.FullName(repo))
	} else {
		fmt.Fprintln(out, repo.URL)
	}

	remoteURL := formatRemoteURL(cmd, repo)

	if projectDirErr == nil {
		_, err = git.AddRemote("origin", remoteURL)
		if err != nil {
			return err
		}
		if isTTY {
			fmt.Fprintf(out, "%s Added remote %s\n", greenCheck, remoteURL)
		}
	} else if isTTY {
		doSetup := false
		err := Confirm(fmt.Sprintf("Create a local project directory for %s?", ghrepo.FullName(repo)), &doSetup)
		if err != nil {
			return err
		}

		if doSetup {
			path := repo.Name

			gitInit := git.GitCommand("init", path)
			gitInit.Stdout = os.Stdout
			gitInit.Stderr = os.Stderr
			err = run.PrepareCmd(gitInit).Run()
			if err != nil {
				return err
			}
			gitRemoteAdd := git.GitCommand("-C", path, "remote", "add", "origin", remoteURL)
			gitRemoteAdd.Stdout = os.Stdout
			gitRemoteAdd.Stderr = os.Stderr
			err = run.PrepareCmd(gitRemoteAdd).Run()
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "%s Initialized repository in './%s/'\n", greenCheck, path)
		}
	}

	return nil
}

var Since = func(t time.Time) time.Duration {
	return time.Since(t)
}

func repoFork(cmd *cobra.Command, args []string) error {
	ctx := contextForCommand(cmd)

	clonePref, err := cmd.Flags().GetString("clone")
	if err != nil {
		return err
	}
	remotePref, err := cmd.Flags().GetString("remote")
	if err != nil {
		return err
	}

	apiClient, err := apiClientForContext(ctx)
	if err != nil {
		return fmt.Errorf("unable to create client: %w", err)
	}

	var repoToFork ghrepo.Interface
	inParent := false // whether or not we're forking the repo we're currently "in"
	if len(args) == 0 {
		baseRepo, err := determineBaseRepo(apiClient, cmd, ctx)
		if err != nil {
			return fmt.Errorf("unable to determine base repository: %w", err)
		}
		inParent = true
		repoToFork = baseRepo
	} else {
		repoArg := args[0]

		if utils.IsURL(repoArg) {
			parsedURL, err := url.Parse(repoArg)
			if err != nil {
				return fmt.Errorf("did not understand argument: %w", err)
			}

			repoToFork, err = ghrepo.FromURL(parsedURL)
			if err != nil {
				return fmt.Errorf("did not understand argument: %w", err)
			}

		} else if strings.HasPrefix(repoArg, "git@") {
			parsedURL, err := git.ParseURL(repoArg)
			if err != nil {
				return fmt.Errorf("did not understand argument: %w", err)
			}
			repoToFork, err = ghrepo.FromURL(parsedURL)
			if err != nil {
				return fmt.Errorf("did not understand argument: %w", err)
			}
		} else {
			repoToFork, err = ghrepo.FromFullName(repoArg)
			if err != nil {
				return fmt.Errorf("argument error: %w", err)
			}
		}
	}

	if !connectedToTerminal(cmd) {
		if (inParent && remotePref == "prompt") || (!inParent && clonePref == "prompt") {
			return errors.New("--remote or --clone must be explicitly set when not attached to tty")
		}
	}

	greenCheck := utils.Green("✓")
	stderr := colorableErr(cmd)
	s := utils.Spinner(stderr)
	stopSpinner := func() {}

	if connectedToTerminal(cmd) {
		loading := utils.Gray("Forking ") + utils.Bold(utils.Gray(ghrepo.FullName(repoToFork))) + utils.Gray("...")
		s.Suffix = " " + loading
		s.FinalMSG = utils.Gray(fmt.Sprintf("- %s\n", loading))
		utils.StartSpinner(s)
		stopSpinner = func() {
			utils.StopSpinner(s)

		}
	}

	forkedRepo, err := api.ForkRepo(apiClient, repoToFork)
	if err != nil {
		stopSpinner()
		return fmt.Errorf("failed to fork: %w", err)
	}

	stopSpinner()

	// This is weird. There is not an efficient way to determine via the GitHub API whether or not a
	// given user has forked a given repo. We noticed, also, that the create fork API endpoint just
	// returns the fork repo data even if it already exists -- with no change in status code or
	// anything. We thus check the created time to see if the repo is brand new or not; if it's not,
	// we assume the fork already existed and report an error.
	createdAgo := Since(forkedRepo.CreatedAt)
	if createdAgo > time.Minute {
		if connectedToTerminal(cmd) {
			fmt.Fprintf(stderr, "%s %s %s\n",
				utils.Yellow("!"),
				utils.Bold(ghrepo.FullName(forkedRepo)),
				"already exists")
		} else {
			fmt.Fprintf(stderr, "%s already exists", ghrepo.FullName(forkedRepo))
			return nil
		}
	} else {
		if connectedToTerminal(cmd) {
			fmt.Fprintf(stderr, "%s Created fork %s\n", greenCheck, utils.Bold(ghrepo.FullName(forkedRepo)))
		}
	}

	if (inParent && remotePref == "false") || (!inParent && clonePref == "false") {
		return nil
	}

	if inParent {
		remotes, err := ctx.Remotes()
		if err != nil {
			return err
		}
		if remote, err := remotes.FindByRepo(forkedRepo.RepoOwner(), forkedRepo.RepoName()); err == nil {
			if connectedToTerminal(cmd) {
				fmt.Fprintf(stderr, "%s Using existing remote %s\n", greenCheck, utils.Bold(remote.Name))
			}
			return nil
		}

		remoteDesired := remotePref == "true"
		if remotePref == "prompt" {
			err = Confirm("Would you like to add a remote for the fork?", &remoteDesired)
			if err != nil {
				return fmt.Errorf("failed to prompt: %w", err)
			}
		}
		if remoteDesired {
			remoteName := "origin"

			remotes, err := ctx.Remotes()
			if err != nil {
				return err
			}
			if _, err := remotes.FindByName(remoteName); err == nil {
				renameTarget := "upstream"
				renameCmd := git.GitCommand("remote", "rename", remoteName, renameTarget)
				err = run.PrepareCmd(renameCmd).Run()
				if err != nil {
					return err
				}
				if connectedToTerminal(cmd) {
					fmt.Fprintf(stderr, "%s Renamed %s remote to %s\n", greenCheck, utils.Bold(remoteName), utils.Bold(renameTarget))
				}
			}

			forkedRepoCloneURL := formatRemoteURL(cmd, forkedRepo)

			_, err = git.AddRemote(remoteName, forkedRepoCloneURL)
			if err != nil {
				return fmt.Errorf("failed to add remote: %w", err)
			}

			if connectedToTerminal(cmd) {
				fmt.Fprintf(stderr, "%s Added remote %s\n", greenCheck, utils.Bold(remoteName))
			}
		}
	} else {
		cloneDesired := clonePref == "true"
		if clonePref == "prompt" {
			err = Confirm("Would you like to clone the fork?", &cloneDesired)
			if err != nil {
				return fmt.Errorf("failed to prompt: %w", err)
			}
		}
		if cloneDesired {
			forkedRepoCloneURL := formatRemoteURL(cmd, forkedRepo)
			cloneDir, err := git.RunClone(forkedRepoCloneURL, []string{})
			if err != nil {
				return fmt.Errorf("failed to clone fork: %w", err)
			}

			// TODO This is overly wordy and I'd like to streamline this.
			cfg, err := ctx.Config()
			if err != nil {
				return err
			}
			protocol, err := cfg.Get("", "git_protocol")
			if err != nil {
				return err
			}

			upstreamURL := ghrepo.FormatRemoteURL(repoToFork, protocol)

			err = git.AddUpstreamRemote(upstreamURL, cloneDir)
			if err != nil {
				return err
			}

			if connectedToTerminal(cmd) {
				fmt.Fprintf(stderr, "%s Cloned fork\n", greenCheck)
			}
		}
	}

	return nil
}

var Confirm = func(prompt string, result *bool) error {
	p := &survey.Confirm{
		Message: prompt,
		Default: true,
	}
	return survey.AskOne(p, result)
}

func repoCredits(cmd *cobra.Command, args []string) error {
	return credits(cmd, args)
}
