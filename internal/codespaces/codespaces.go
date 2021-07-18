package codespaces

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/github/ghcs/api"
	"github.com/github/go-liveshare"
)

var (
	ErrNoCodespaces = errors.New("You have no codespaces.")
)

func ChooseCodespace(ctx context.Context, apiClient *api.API, user *api.User) (*api.Codespace, error) {
	codespaces, err := apiClient.ListCodespaces(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error getting codespaces: %v", err)
	}

	if len(codespaces) == 0 {
		return nil, ErrNoCodespaces
	}

	codespaces.SortByCreatedAt()

	codespacesByName := make(map[string]*api.Codespace)
	codespacesNames := make([]string, 0, len(codespaces))
	for _, codespace := range codespaces {
		codespacesByName[codespace.Name] = codespace
		codespacesNames = append(codespacesNames, codespace.Name)
	}

	sshSurvey := []*survey.Question{
		{
			Name: "codespace",
			Prompt: &survey.Select{
				Message: "Choose Codespace:",
				Options: codespacesNames,
				Default: codespacesNames[0],
			},
			Validate: survey.Required,
		},
	}

	answers := struct {
		Codespace string
	}{}
	if err := survey.Ask(sshSurvey, &answers); err != nil {
		return nil, fmt.Errorf("error getting answers: %v", err)
	}

	codespace := codespacesByName[answers.Codespace]
	return codespace, nil
}

func ConnectToLiveshare(ctx context.Context, apiClient *api.API, token string, codespace *api.Codespace) (client *liveshare.Client, err error) {
	if codespace.Environment.State != api.CodespaceEnvironmentStateAvailable {
		fmt.Println("Starting your codespace...") // TODO(josebalius): better way of notifying of events
		if err := apiClient.StartCodespace(ctx, token, codespace); err != nil {
			return nil, fmt.Errorf("error starting codespace: %v", err)
		}
	}

	retries := 0
	for codespace.Environment.Connection.SessionID == "" || codespace.Environment.State != api.CodespaceEnvironmentStateAvailable {
		if retries > 1 {
			if retries%2 == 0 {
				fmt.Print(".")
			}

			time.Sleep(1 * time.Second)
		}

		if retries == 30 {
			return nil, errors.New("timed out while waiting for the codespace to start")
		}

		codespace, err = apiClient.GetCodespace(ctx, token, codespace.OwnerLogin, codespace.Name)
		if err != nil {
			return nil, fmt.Errorf("error getting codespace: %v", err)
		}

		retries += 1
	}

	if retries >= 2 {
		fmt.Print("\n")
	}

	fmt.Println("Connecting to your codespace...")

	liveShare, err := liveshare.New(
		liveshare.WithWorkspaceID(codespace.Environment.Connection.SessionID),
		liveshare.WithToken(codespace.Environment.Connection.SessionToken),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating live share: %v", err)
	}

	liveShareClient := liveShare.NewClient()
	if err := liveShareClient.Join(ctx); err != nil {
		return nil, fmt.Errorf("error joining liveshare client: %v", err)
	}

	return liveShareClient, nil
}
