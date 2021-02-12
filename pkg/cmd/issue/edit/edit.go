package edit

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/api"
	"github.com/cli/cli/internal/ghrepo"
	shared "github.com/cli/cli/pkg/cmd/issue/shared"
	prShared "github.com/cli/cli/pkg/cmd/pr/shared"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
)

type EditOptions struct {
	HttpClient func() (*http.Client, error)
	IO         *iostreams.IOStreams
	BaseRepo   func() (ghrepo.Interface, error)

	DetermineEditor    func() (string, error)
	FieldsToEditSurvey func(*prShared.EditableOptions) error
	EditableSurvey     func(string, *prShared.EditableOptions) error
	FetchOptions       func(*api.Client, ghrepo.Interface, *prShared.EditableOptions) error

	SelectorArg string
	Interactive bool

	prShared.EditableOptions
}

func NewCmdEdit(f *cmdutil.Factory, runF func(*EditOptions) error) *cobra.Command {
	opts := &EditOptions{
		IO:                 f.IOStreams,
		HttpClient:         f.HttpClient,
		DetermineEditor:    func() (string, error) { return cmdutil.DetermineEditor(f.Config) },
		FieldsToEditSurvey: prShared.FieldsToEditSurvey,
		EditableSurvey:     prShared.EditableSurvey,
		FetchOptions:       prShared.FetchOptions,
	}

	cmd := &cobra.Command{
		Use:   "edit {<number> | <url>}",
		Short: "Edit an issue",
		Example: heredoc.Doc(`
			$ gh issue edit 23 --title "I found a bug" --body "Nothing works"
			$ gh issue edit 23 --label "bug,help wanted"
			$ gh issue edit 23 --label bug --label "help wanted"
			$ gh issue edit 23 --assignee monalisa,hubot
			$ gh issue edit 23 --assignee @me
			$ gh issue edit 23 --project "Roadmap"
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// support `-R, --repo` override
			opts.BaseRepo = f.BaseRepo

			opts.SelectorArg = args[0]

			flags := cmd.Flags()
			if flags.Changed("title") {
				opts.EditableOptions.TitleEdited = true
			}
			if flags.Changed("body") {
				opts.EditableOptions.BodyEdited = true
			}
			if flags.Changed("assignee") {
				opts.EditableOptions.AssigneesEdited = true
			}
			if flags.Changed("label") {
				opts.EditableOptions.LabelsEdited = true
			}
			if flags.Changed("project") {
				opts.EditableOptions.ProjectsEdited = true
			}
			if flags.Changed("milestone") {
				opts.EditableOptions.MilestoneEdited = true
			}

			if !opts.EditableOptions.Dirty() {
				opts.Interactive = true
			}

			if opts.Interactive && !opts.IO.CanPrompt() {
				return &cmdutil.FlagError{Err: errors.New("--tile, --body, --assignee, --label, --project, or --milestone required when not running interactively")}
			}

			if runF != nil {
				return runF(opts)
			}

			return editRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.EditableOptions.Title, "title", "t", "", "Revise the issue title.")
	cmd.Flags().StringVarP(&opts.EditableOptions.Body, "body", "b", "", "Revise the issue body.")
	cmd.Flags().StringSliceVarP(&opts.EditableOptions.Assignees, "assignee", "a", nil, "Set assigned people by their `login`. Use \"@me\" to self-assign.")
	cmd.Flags().StringSliceVarP(&opts.EditableOptions.Labels, "label", "l", nil, "Set the issue labels by `name`")
	cmd.Flags().StringSliceVarP(&opts.EditableOptions.Projects, "project", "p", nil, "Set the projects the issue belongs to by `name`")
	cmd.Flags().StringVarP(&opts.EditableOptions.Milestone, "milestone", "m", "", "Set the milestone the issue belongs to by `name`")

	return cmd
}

func editRun(opts *EditOptions) error {
	httpClient, err := opts.HttpClient()
	if err != nil {
		return err
	}
	apiClient := api.NewClientFromHTTP(httpClient)

	issue, repo, err := shared.IssueFromArg(apiClient, opts.BaseRepo, opts.SelectorArg)
	if err != nil {
		return err
	}

	editOptions := opts.EditableOptions
	editOptions.TitleDefault = issue.Title
	editOptions.BodyDefault = issue.Body
	editOptions.AssigneesDefault = issue.Assignees
	editOptions.LabelsDefault = issue.Labels
	editOptions.ProjectsDefault = issue.ProjectCards
	editOptions.MilestoneDefault = issue.Milestone

	if opts.Interactive {
		err = opts.FieldsToEditSurvey(&editOptions)
		if err != nil {
			return err
		}
	}

	opts.IO.StartProgressIndicator()
	err = opts.FetchOptions(apiClient, repo, &editOptions)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	if opts.Interactive {
		editorCommand, err := opts.DetermineEditor()
		if err != nil {
			return err
		}
		err = opts.EditableSurvey(editorCommand, &editOptions)
		if err != nil {
			return err
		}
	}

	opts.IO.StartProgressIndicator()
	err = updateIssue(apiClient, repo, issue.ID, editOptions)
	opts.IO.StopProgressIndicator()
	if err != nil {
		return err
	}

	fmt.Fprintln(opts.IO.Out, issue.URL)

	return nil
}

func updateIssue(client *api.Client, repo ghrepo.Interface, id string, options prShared.EditableOptions) error {
	params := githubv4.UpdateIssueInput{ID: id}
	if options.TitleEdited {
		title := githubv4.String(options.Title)
		params.Title = &title
	}
	if options.BodyEdited {
		body := githubv4.String(options.Body)
		params.Body = &body
	}
	if options.AssigneesEdited {
		meReplacer := prShared.NewMeReplacer(client, repo.RepoHost())
		assignees, err := meReplacer.ReplaceSlice(options.Assignees)
		if err != nil {
			return err
		}
		ids, err := options.Metadata.MembersToIDs(assignees)
		if err != nil {
			return err
		}
		assigneeIDs := make([]githubv4.ID, len(ids))
		for i, v := range ids {
			assigneeIDs[i] = v
		}
		params.AssigneeIDs = &assigneeIDs
	}
	if options.LabelsEdited {
		ids, err := options.Metadata.LabelsToIDs(options.Labels)
		if err != nil {
			return err
		}
		labelIDs := make([]githubv4.ID, len(ids))
		for i, v := range ids {
			labelIDs[i] = v
		}
		params.LabelIDs = &labelIDs
	}
	if options.ProjectsEdited {
		ids, err := options.Metadata.ProjectsToIDs(options.Projects)
		if err != nil {
			return err
		}
		projectIDs := make([]githubv4.ID, len(ids))
		for i, v := range ids {
			projectIDs[i] = v
		}
		params.ProjectIDs = &projectIDs
	}
	if options.MilestoneEdited {
		id, err := options.Metadata.MilestoneToID(options.Milestone)
		if err != nil {
			return err
		}
		milestoneID := githubv4.ID(id)
		params.MilestoneID = &milestoneID
	}
	return api.IssueUpdate(client, repo, params)
}
