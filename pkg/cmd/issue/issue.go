package issue

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/internal/ghrepo"
	cmdClose "github.com/cli/cli/pkg/cmd/issue/close"
	cmdCreate "github.com/cli/cli/pkg/cmd/issue/create"
	cmdList "github.com/cli/cli/pkg/cmd/issue/list"
	cmdReopen "github.com/cli/cli/pkg/cmd/issue/reopen"
	cmdStatus "github.com/cli/cli/pkg/cmd/issue/status"
	cmdView "github.com/cli/cli/pkg/cmd/issue/view"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdIssue(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue <command>",
		Short: "Manage issues",
		Long:  `Work with GitHub issues`,
		Example: heredoc.Doc(`
			$ gh issue list
			$ gh issue create --label bug
			$ gh issue view --web
		`),
		Annotations: map[string]string{
			"IsCore": "true",
			"help:arguments": heredoc.Doc(`
				An issue can be supplied as argument in any of the following formats:
				- by number, e.g. "123"; or
				- by URL, e.g. "https://github.com/OWNER/REPO/issues/123".
			`),
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if repo, _ := cmd.Flags().GetString("repo"); repo != "" {
				// NOTE: this mutates the factory
				f.BaseRepo = func() (ghrepo.Interface, error) {
					return ghrepo.FromFullName(repo)
				}
			}
		},
	}

	cmd.PersistentFlags().StringP("repo", "R", "", "Select another repository using the `OWNER/REPO` format")

	cmd.AddCommand(cmdClose.NewCmdClose(f, nil))
	cmd.AddCommand(cmdCreate.NewCmdCreate(f, nil))
	cmd.AddCommand(cmdList.NewCmdList(f, nil))
	cmd.AddCommand(cmdReopen.NewCmdReopen(f, nil))
	cmd.AddCommand(cmdStatus.NewCmdStatus(f, nil))
	cmd.AddCommand(cmdView.NewCmdView(f, nil))

	return cmd
}
