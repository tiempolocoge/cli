package shared

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/cli/cli/api"
	"github.com/cli/cli/internal/ghrepo"
)

func WithPrAndIssueQueryParams(baseURL string, state IssueMetadataState) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	if state.Title != "" {
		q.Set("title", state.Title)
	}
	if state.Body != "" {
		q.Set("body", state.Body)
	}
	if len(state.Assignees) > 0 {
		q.Set("assignees", strings.Join(state.Assignees, ","))
	}
	if len(state.Labels) > 0 {
		q.Set("labels", strings.Join(state.Labels, ","))
	}
	if len(state.Projects) > 0 {
		q.Set("projects", strings.Join(state.Projects, ","))
	}
	if len(state.Milestones) > 0 {
		q.Set("milestone", state.Milestones[0])
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// Ensure that tb.MetadataResult object exists and contains enough pre-fetched API data to be able
// to resolve all object listed in tb to GraphQL IDs.
func fillMetadata(client *api.Client, baseRepo ghrepo.Interface, tb *IssueMetadataState) error {
	resolveInput := api.RepoResolveInput{}

	if len(tb.Assignees) > 0 && (tb.MetadataResult == nil || len(tb.MetadataResult.AssignableUsers) == 0) {
		resolveInput.Assignees = tb.Assignees
	}

	if len(tb.Reviewers) > 0 && (tb.MetadataResult == nil || len(tb.MetadataResult.AssignableUsers) == 0) {
		resolveInput.Reviewers = tb.Reviewers
	}

	if len(tb.Labels) > 0 && (tb.MetadataResult == nil || len(tb.MetadataResult.Labels) == 0) {
		resolveInput.Labels = tb.Labels
	}

	if len(tb.Projects) > 0 && (tb.MetadataResult == nil || len(tb.MetadataResult.Projects) == 0) {
		resolveInput.Projects = tb.Projects
	}

	if len(tb.Milestones) > 0 && (tb.MetadataResult == nil || len(tb.MetadataResult.Milestones) == 0) {
		resolveInput.Milestones = tb.Milestones
	}

	metadataResult, err := api.RepoResolveMetadataIDs(client, baseRepo, resolveInput)
	if err != nil {
		return err
	}

	if tb.MetadataResult == nil {
		tb.MetadataResult = metadataResult
	} else {
		tb.MetadataResult.Merge(metadataResult)
	}

	return nil
}

func AddMetadataToIssueParams(client *api.Client, baseRepo ghrepo.Interface, params map[string]interface{}, tb *IssueMetadataState) error {
	if !tb.HasMetadata() {
		return nil
	}

	var err error
	if tb.Assignees, err = ReplaceAtMeLogin(tb.Assignees, client, baseRepo); err != nil {
		return err
	}

	if err := fillMetadata(client, baseRepo, tb); err != nil {
		return err
	}

	assigneeIDs, err := tb.MetadataResult.MembersToIDs(tb.Assignees)
	if err != nil {
		return fmt.Errorf("could not assign user: %w", err)
	}
	params["assigneeIds"] = assigneeIDs

	labelIDs, err := tb.MetadataResult.LabelsToIDs(tb.Labels)
	if err != nil {
		return fmt.Errorf("could not add label: %w", err)
	}
	params["labelIds"] = labelIDs

	projectIDs, err := tb.MetadataResult.ProjectsToIDs(tb.Projects)
	if err != nil {
		return fmt.Errorf("could not add to project: %w", err)
	}
	params["projectIds"] = projectIDs

	if len(tb.Milestones) > 0 {
		milestoneID, err := tb.MetadataResult.MilestoneToID(tb.Milestones[0])
		if err != nil {
			return fmt.Errorf("could not add to milestone '%s': %w", tb.Milestones[0], err)
		}
		params["milestoneId"] = milestoneID
	}

	if len(tb.Reviewers) == 0 {
		return nil
	}

	var userReviewers []string
	var teamReviewers []string
	for _, r := range tb.Reviewers {
		if strings.ContainsRune(r, '/') {
			teamReviewers = append(teamReviewers, r)
		} else {
			userReviewers = append(userReviewers, r)
		}
	}

	userReviewerIDs, err := tb.MetadataResult.MembersToIDs(userReviewers)
	if err != nil {
		return fmt.Errorf("could not request reviewer: %w", err)
	}
	params["userReviewerIds"] = userReviewerIDs

	teamReviewerIDs, err := tb.MetadataResult.TeamsToIDs(teamReviewers)
	if err != nil {
		return fmt.Errorf("could not request reviewer: %w", err)
	}
	params["teamReviewerIds"] = teamReviewerIDs

	return nil
}

type FilterOptions struct {
	Entity     string
	State      string
	Assignee   string
	Labels     []string
	Author     string
	BaseBranch string
	Mention    string
	Milestone  string
}

func ListURLWithQuery(listURL string, options FilterOptions) (string, error) {
	u, err := url.Parse(listURL)
	if err != nil {
		return "", err
	}
	query := fmt.Sprintf("is:%s ", options.Entity)
	if options.State != "all" {
		query += fmt.Sprintf("is:%s ", options.State)
	}
	if options.Assignee != "" {
		query += fmt.Sprintf("assignee:%s ", options.Assignee)
	}
	for _, label := range options.Labels {
		query += fmt.Sprintf("label:%s ", quoteValueForQuery(label))
	}
	if options.Author != "" {
		query += fmt.Sprintf("author:%s ", options.Author)
	}
	if options.BaseBranch != "" {
		query += fmt.Sprintf("base:%s ", options.BaseBranch)
	}
	if options.Mention != "" {
		query += fmt.Sprintf("mentions:%s ", options.Mention)
	}
	if options.Milestone != "" {
		query += fmt.Sprintf("milestone:%s ", quoteValueForQuery(options.Milestone))
	}
	q := u.Query()
	q.Set("q", strings.TrimSuffix(query, " "))
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func quoteValueForQuery(v string) string {
	if strings.ContainsAny(v, " \"\t\r\n") {
		return fmt.Sprintf("%q", v)
	}
	return v
}

// ReplaceAtMeLogin iterates over the list of specified login names, replacing
// any "@me" mentions with the current user LoginName.
func ReplaceAtMeLogin(logins []string, client *api.Client, repo ghrepo.Interface) ([]string, error) {
	for i, u := range logins {
		if strings.EqualFold(u, "@me") {
			login, err := api.CurrentLoginName(client, repo.RepoHost())
			if err != nil {
				return logins, fmt.Errorf("@me resolve: failed obtaining user id: %w", err)
			}
			logins[i] = login
		}
	}
	return logins, nil
}
