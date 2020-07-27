package command

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/briandowns/spinner"

	"github.com/cli/cli/context"
	"github.com/cli/cli/internal/run"
	"github.com/cli/cli/pkg/httpmock"
	"github.com/cli/cli/test"
	"github.com/cli/cli/utils"
	"github.com/stretchr/testify/assert"
)

func stubSpinner() {
	// not bothering with teardown since we never want spinners when doing tests
	utils.StartSpinner = func(_ *spinner.Spinner) {
	}
	utils.StopSpinner = func(_ *spinner.Spinner) {
	}
}

func TestRepoFork_nontty_insufficient_flags(t *testing.T) {
	initBlankContext("", "OWNER/REPO", "master")
	defer stubSince(2 * time.Second)()
	http := initFakeHTTP()
	http.StubRepoResponse("OWNER", "REPO")
	defer stubTerminal(false)()

	_, err := RunCommand("repo fork")
	if err == nil {
		t.Fatal("expected error")
	}

	assert.Equal(t, "--remote or --clone must be explicitly set when not attached to tty", err.Error())
}

func TestRepoFork_in_parent_nontty(t *testing.T) {
	initBlankContext("", "OWNER/REPO", "master")
	defer stubSince(2 * time.Second)()
	http := initFakeHTTP()
	http.StubRepoResponse("OWNER", "REPO")
	defer http.StubWithFixture(200, "forkResult.json")()
	defer stubTerminal(false)()

	cs, restore := test.InitCmdStubber()
	defer restore()

	cs.Stub("") // git remote rename
	cs.Stub("") // git remote add

	output, err := RunCommand("repo fork --remote")
	if err != nil {
		t.Fatalf("error running command `repo fork`: %v", err)
	}

	eq(t, strings.Join(cs.Calls[0].Args, " "), "git remote rename origin upstream")
	eq(t, strings.Join(cs.Calls[1].Args, " "), "git remote add -f origin https://github.com/someone/REPO.git")

	assert.Equal(t, "", output.String())
	assert.Equal(t, "", output.Stderr())
}

func TestRepoFork_outside_parent_nontty(t *testing.T) {
	defer stubSince(2 * time.Second)()
	http := initFakeHTTP()
	defer http.StubWithFixture(200, "forkResult.json")()
	defer stubTerminal(false)()

	cs, restore := test.InitCmdStubber()
	defer restore()

	cs.Stub("") // git clone
	cs.Stub("") // git remote add

	output, err := RunCommand("repo fork --clone OWNER/REPO")
	if err != nil {
		t.Errorf("error running command `repo fork`: %v", err)
	}

	eq(t, output.String(), "")

	eq(t, strings.Join(cs.Calls[0].Args, " "), "git clone https://github.com/someone/REPO.git")
	eq(t, strings.Join(cs.Calls[1].Args, " "), "git -C REPO remote add -f upstream https://github.com/OWNER/REPO.git")

	assert.Equal(t, output.Stderr(), "")
}

func TestRepoFork_already_forked(t *testing.T) {
	stubSpinner()
	initContext = func() context.Context {
		ctx := context.NewBlank()
		ctx.SetBaseRepo("OWNER/REPO")
		ctx.SetBranch("master")
		ctx.SetRemotes(map[string]string{
			"origin": "OWNER/REPO",
		})
		return ctx
	}
	http := initFakeHTTP()
	http.StubRepoResponse("OWNER", "REPO")
	defer http.StubWithFixture(200, "forkResult.json")()
	defer stubTerminal(true)()

	output, err := RunCommand("repo fork --remote=false")
	if err != nil {
		t.Errorf("got unexpected error: %v", err)
	}
	r := regexp.MustCompile(`someone/REPO.*already exists`)
	if !r.MatchString(output.Stderr()) {
		t.Errorf("output did not match regexp /%s/\n> output\n%s\n", r, output.Stderr())
		return
	}
}

func TestRepoFork_reuseRemote(t *testing.T) {
	stubSpinner()
	initContext = func() context.Context {
		ctx := context.NewBlank()
		ctx.SetBaseRepo("OWNER/REPO")
		ctx.SetBranch("master")
		ctx.SetRemotes(map[string]string{
			"upstream": "OWNER/REPO",
			"origin":   "someone/REPO",
		})
		return ctx
	}
	http := initFakeHTTP()
	http.StubRepoResponse("OWNER", "REPO")
	defer http.StubWithFixture(200, "forkResult.json")()
	defer stubTerminal(true)()

	output, err := RunCommand("repo fork")
	if err != nil {
		t.Errorf("got unexpected error: %v", err)
	}
	r := regexp.MustCompile(`Using existing remote.*origin`)
	if !r.MatchString(output.Stderr()) {
		t.Errorf("output did not match: %q", output.Stderr())
		return
	}
}

func stubSince(d time.Duration) func() {
	originalSince := Since
	Since = func(t time.Time) time.Duration {
		return d
	}
	return func() {
		Since = originalSince
	}
}

func TestRepoFork_in_parent(t *testing.T) {
	stubSpinner()
	initBlankContext("", "OWNER/REPO", "master")
	defer stubSince(2 * time.Second)()
	http := initFakeHTTP()
	http.StubRepoResponse("OWNER", "REPO")
	defer http.StubWithFixture(200, "forkResult.json")()
	defer stubTerminal(true)()

	output, err := RunCommand("repo fork --remote=false")
	if err != nil {
		t.Errorf("error running command `repo fork`: %v", err)
	}

	eq(t, output.String(), "")

	r := regexp.MustCompile(`Created fork.*someone/REPO`)
	if !r.MatchString(output.Stderr()) {
		t.Errorf("output did not match regexp /%s/\n> output\n%s\n", r, output)
		return
	}
}

func TestRepoFork_outside(t *testing.T) {
	stubSpinner()
	tests := []struct {
		name string
		args string
	}{
		{
			name: "url arg",
			args: "repo fork --clone=false http://github.com/OWNER/REPO.git",
		},
		{
			name: "full name arg",
			args: "repo fork --clone=false OWNER/REPO",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer stubSince(2 * time.Second)()
			http := initFakeHTTP()
			defer http.StubWithFixture(200, "forkResult.json")()
			defer stubTerminal(true)()

			output, err := RunCommand(tt.args)
			if err != nil {
				t.Errorf("error running command `repo fork`: %v", err)
			}

			eq(t, output.String(), "")

			r := regexp.MustCompile(`Created fork.*someone/REPO`)
			if !r.MatchString(output.Stderr()) {
				t.Errorf("output did not match regexp /%s/\n> output\n%s\n", r, output)
				return
			}
		})
	}
}

func TestRepoFork_in_parent_yes(t *testing.T) {
	stubSpinner()
	initBlankContext("", "OWNER/REPO", "master")
	defer stubSince(2 * time.Second)()
	http := initFakeHTTP()
	http.StubRepoResponse("OWNER", "REPO")
	defer http.StubWithFixture(200, "forkResult.json")()
	defer stubTerminal(true)()

	var seenCmds []*exec.Cmd
	defer run.SetPrepareCmd(func(cmd *exec.Cmd) run.Runnable {
		seenCmds = append(seenCmds, cmd)
		return &test.OutputStub{}
	})()

	output, err := RunCommand("repo fork --remote")
	if err != nil {
		t.Errorf("error running command `repo fork`: %v", err)
	}

	expectedCmds := []string{
		"git remote rename origin upstream",
		"git remote add -f origin https://github.com/someone/REPO.git",
	}

	for x, cmd := range seenCmds {
		eq(t, strings.Join(cmd.Args, " "), expectedCmds[x])
	}

	eq(t, output.String(), "")

	test.ExpectLines(t, output.Stderr(),
		"Created fork.*someone/REPO",
		"Added remote.*origin")
}

func TestRepoFork_outside_yes(t *testing.T) {
	stubSpinner()
	defer stubSince(2 * time.Second)()
	http := initFakeHTTP()
	defer http.StubWithFixture(200, "forkResult.json")()
	defer stubTerminal(true)()

	cs, restore := test.InitCmdStubber()
	defer restore()

	cs.Stub("") // git clone
	cs.Stub("") // git remote add

	output, err := RunCommand("repo fork --clone OWNER/REPO")
	if err != nil {
		t.Errorf("error running command `repo fork`: %v", err)
	}

	eq(t, output.String(), "")

	eq(t, strings.Join(cs.Calls[0].Args, " "), "git clone https://github.com/someone/REPO.git")
	eq(t, strings.Join(cs.Calls[1].Args, " "), "git -C REPO remote add -f upstream https://github.com/OWNER/REPO.git")

	test.ExpectLines(t, output.Stderr(),
		"Created fork.*someone/REPO",
		"Cloned fork")
}

func TestRepoFork_outside_survey_yes(t *testing.T) {
	stubSpinner()
	defer stubSince(2 * time.Second)()
	http := initFakeHTTP()
	defer http.StubWithFixture(200, "forkResult.json")()
	defer stubTerminal(true)()

	cs, restore := test.InitCmdStubber()
	defer restore()

	cs.Stub("") // git clone
	cs.Stub("") // git remote add

	oldConfirm := Confirm
	Confirm = func(_ string, result *bool) error {
		*result = true
		return nil
	}
	defer func() { Confirm = oldConfirm }()

	output, err := RunCommand("repo fork OWNER/REPO")
	if err != nil {
		t.Errorf("error running command `repo fork`: %v", err)
	}

	eq(t, output.String(), "")

	eq(t, strings.Join(cs.Calls[0].Args, " "), "git clone https://github.com/someone/REPO.git")
	eq(t, strings.Join(cs.Calls[1].Args, " "), "git -C REPO remote add -f upstream https://github.com/OWNER/REPO.git")

	test.ExpectLines(t, output.Stderr(),
		"Created fork.*someone/REPO",
		"Cloned fork")
}

func TestRepoFork_outside_survey_no(t *testing.T) {
	stubSpinner()
	defer stubSince(2 * time.Second)()
	http := initFakeHTTP()
	defer http.StubWithFixture(200, "forkResult.json")()
	defer stubTerminal(true)()

	cmdRun := false
	defer run.SetPrepareCmd(func(cmd *exec.Cmd) run.Runnable {
		cmdRun = true
		return &test.OutputStub{}
	})()

	oldConfirm := Confirm
	Confirm = func(_ string, result *bool) error {
		*result = false
		return nil
	}
	defer func() { Confirm = oldConfirm }()

	output, err := RunCommand("repo fork OWNER/REPO")
	if err != nil {
		t.Errorf("error running command `repo fork`: %v", err)
	}

	eq(t, output.String(), "")

	eq(t, cmdRun, false)

	r := regexp.MustCompile(`Created fork.*someone/REPO`)
	if !r.MatchString(output.Stderr()) {
		t.Errorf("output did not match regexp /%s/\n> output\n%s\n", r, output)
		return
	}
}

func TestRepoFork_in_parent_survey_yes(t *testing.T) {
	stubSpinner()
	initBlankContext("", "OWNER/REPO", "master")
	defer stubSince(2 * time.Second)()
	http := initFakeHTTP()
	http.StubRepoResponse("OWNER", "REPO")
	defer http.StubWithFixture(200, "forkResult.json")()
	defer stubTerminal(true)()

	var seenCmds []*exec.Cmd
	defer run.SetPrepareCmd(func(cmd *exec.Cmd) run.Runnable {
		seenCmds = append(seenCmds, cmd)
		return &test.OutputStub{}
	})()

	oldConfirm := Confirm
	Confirm = func(_ string, result *bool) error {
		*result = true
		return nil
	}
	defer func() { Confirm = oldConfirm }()

	output, err := RunCommand("repo fork")
	if err != nil {
		t.Errorf("error running command `repo fork`: %v", err)
	}

	expectedCmds := []string{
		"git remote rename origin upstream",
		"git remote add -f origin https://github.com/someone/REPO.git",
	}

	for x, cmd := range seenCmds {
		eq(t, strings.Join(cmd.Args, " "), expectedCmds[x])
	}

	eq(t, output.String(), "")

	test.ExpectLines(t, output.Stderr(),
		"Created fork.*someone/REPO",
		"Renamed.*origin.*remote to.*upstream",
		"Added remote.*origin")
}

func TestRepoFork_in_parent_survey_no(t *testing.T) {
	stubSpinner()
	initBlankContext("", "OWNER/REPO", "master")
	defer stubSince(2 * time.Second)()
	http := initFakeHTTP()
	http.StubRepoResponse("OWNER", "REPO")
	defer http.StubWithFixture(200, "forkResult.json")()
	defer stubTerminal(true)()

	cmdRun := false
	defer run.SetPrepareCmd(func(cmd *exec.Cmd) run.Runnable {
		cmdRun = true
		return &test.OutputStub{}
	})()

	oldConfirm := Confirm
	Confirm = func(_ string, result *bool) error {
		*result = false
		return nil
	}
	defer func() { Confirm = oldConfirm }()

	output, err := RunCommand("repo fork")
	if err != nil {
		t.Errorf("error running command `repo fork`: %v", err)
	}

	eq(t, output.String(), "")

	eq(t, cmdRun, false)

	r := regexp.MustCompile(`Created fork.*someone/REPO`)
	if !r.MatchString(output.Stderr()) {
		t.Errorf("output did not match regexp /%s/\n> output\n%s\n", r, output)
		return
	}
}

func TestRepoCreate(t *testing.T) {
	ctx := context.NewBlank()
	ctx.SetBranch("master")
	initContext = func() context.Context {
		return ctx
	}

	http := initFakeHTTP()
	http.Register(
		httpmock.GraphQL(`mutation RepositoryCreate\b`),
		httpmock.StringResponse(`
		{ "data": { "createRepository": {
			"repository": {
				"id": "REPOID",
				"url": "https://github.com/OWNER/REPO",
				"name": "REPO",
				"owner": {
					"login": "OWNER"
				}
			}
		} } }`))

	var seenCmd *exec.Cmd
	restoreCmd := run.SetPrepareCmd(func(cmd *exec.Cmd) run.Runnable {
		seenCmd = cmd
		return &test.OutputStub{}
	})
	defer restoreCmd()

	output, err := RunCommand("repo create REPO")
	if err != nil {
		t.Errorf("error running command `repo create`: %v", err)
	}

	eq(t, output.String(), "https://github.com/OWNER/REPO\n")
	eq(t, output.Stderr(), "")

	if seenCmd == nil {
		t.Fatal("expected a command to run")
	}
	eq(t, strings.Join(seenCmd.Args, " "), "git remote add -f origin https://github.com/OWNER/REPO.git")

	var reqBody struct {
		Query     string
		Variables struct {
			Input map[string]interface{}
		}
	}

	if len(http.Requests) != 1 {
		t.Fatalf("expected 1 HTTP request, got %d", len(http.Requests))
	}

	bodyBytes, _ := ioutil.ReadAll(http.Requests[0].Body)
	_ = json.Unmarshal(bodyBytes, &reqBody)
	if repoName := reqBody.Variables.Input["name"].(string); repoName != "REPO" {
		t.Errorf("expected %q, got %q", "REPO", repoName)
	}
	if repoVisibility := reqBody.Variables.Input["visibility"].(string); repoVisibility != "PRIVATE" {
		t.Errorf("expected %q, got %q", "PRIVATE", repoVisibility)
	}
	if _, ownerSet := reqBody.Variables.Input["ownerId"]; ownerSet {
		t.Error("expected ownerId not to be set")
	}
}

func TestRepoCreate_org(t *testing.T) {
	ctx := context.NewBlank()
	ctx.SetBranch("master")
	initContext = func() context.Context {
		return ctx
	}

	http := initFakeHTTP()
	http.Register(
		httpmock.REST("GET", "users/ORG"),
		httpmock.StringResponse(`
		{ "node_id": "ORGID"
		}`))
	http.Register(
		httpmock.GraphQL(`mutation RepositoryCreate\b`),
		httpmock.StringResponse(`
		{ "data": { "createRepository": {
			"repository": {
				"id": "REPOID",
				"url": "https://github.com/ORG/REPO",
				"name": "REPO",
				"owner": {
					"login": "ORG"
				}
			}
		} } }`))

	var seenCmd *exec.Cmd
	restoreCmd := run.SetPrepareCmd(func(cmd *exec.Cmd) run.Runnable {
		seenCmd = cmd
		return &test.OutputStub{}
	})
	defer restoreCmd()

	output, err := RunCommand("repo create ORG/REPO")
	if err != nil {
		t.Errorf("error running command `repo create`: %v", err)
	}

	eq(t, output.String(), "https://github.com/ORG/REPO\n")
	eq(t, output.Stderr(), "")

	if seenCmd == nil {
		t.Fatal("expected a command to run")
	}
	eq(t, strings.Join(seenCmd.Args, " "), "git remote add -f origin https://github.com/ORG/REPO.git")

	var reqBody struct {
		Query     string
		Variables struct {
			Input map[string]interface{}
		}
	}

	if len(http.Requests) != 2 {
		t.Fatalf("expected 2 HTTP requests, got %d", len(http.Requests))
	}

	eq(t, http.Requests[0].URL.Path, "/users/ORG")

	bodyBytes, _ := ioutil.ReadAll(http.Requests[1].Body)
	_ = json.Unmarshal(bodyBytes, &reqBody)
	if orgID := reqBody.Variables.Input["ownerId"].(string); orgID != "ORGID" {
		t.Errorf("expected %q, got %q", "ORGID", orgID)
	}
	if _, teamSet := reqBody.Variables.Input["teamId"]; teamSet {
		t.Error("expected teamId not to be set")
	}
}

func TestRepoCreate_orgWithTeam(t *testing.T) {
	ctx := context.NewBlank()
	ctx.SetBranch("master")
	initContext = func() context.Context {
		return ctx
	}

	http := initFakeHTTP()
	http.Register(
		httpmock.REST("GET", "orgs/ORG/teams/monkeys"),
		httpmock.StringResponse(`
		{ "node_id": "TEAMID",
			"organization": { "node_id": "ORGID" }
		}`))
	http.Register(
		httpmock.GraphQL(`mutation RepositoryCreate\b`),
		httpmock.StringResponse(`
		{ "data": { "createRepository": {
			"repository": {
				"id": "REPOID",
				"url": "https://github.com/ORG/REPO",
				"name": "REPO",
				"owner": {
					"login": "ORG"
				}
			}
		} } }`))

	var seenCmd *exec.Cmd
	restoreCmd := run.SetPrepareCmd(func(cmd *exec.Cmd) run.Runnable {
		seenCmd = cmd
		return &test.OutputStub{}
	})
	defer restoreCmd()

	output, err := RunCommand("repo create ORG/REPO --team monkeys")
	if err != nil {
		t.Errorf("error running command `repo create`: %v", err)
	}

	eq(t, output.String(), "https://github.com/ORG/REPO\n")
	eq(t, output.Stderr(), "")

	if seenCmd == nil {
		t.Fatal("expected a command to run")
	}
	eq(t, strings.Join(seenCmd.Args, " "), "git remote add -f origin https://github.com/ORG/REPO.git")

	var reqBody struct {
		Query     string
		Variables struct {
			Input map[string]interface{}
		}
	}

	if len(http.Requests) != 2 {
		t.Fatalf("expected 2 HTTP requests, got %d", len(http.Requests))
	}

	eq(t, http.Requests[0].URL.Path, "/orgs/ORG/teams/monkeys")

	bodyBytes, _ := ioutil.ReadAll(http.Requests[1].Body)
	_ = json.Unmarshal(bodyBytes, &reqBody)
	if orgID := reqBody.Variables.Input["ownerId"].(string); orgID != "ORGID" {
		t.Errorf("expected %q, got %q", "ORGID", orgID)
	}
	if teamID := reqBody.Variables.Input["teamId"].(string); teamID != "TEAMID" {
		t.Errorf("expected %q, got %q", "TEAMID", teamID)
	}
}
