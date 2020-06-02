package command

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cli/cli/internal/config"
	"github.com/cli/cli/test"
)

func TestAliasSet_existing_alias(t *testing.T) {
	cfg := `---
hosts:
  github.com:
    user: OWNER
    oauth_token: token123
aliases:
  co: pr checkout
`
	initBlankContext(cfg, "OWNER/REPO", "trunk")

	_, err := RunCommand("alias set co pr checkout")

	if err == nil {
		t.Fatal("expected error")
	}

	eq(t, err.Error(), "alias co already exists")
}

func TestAliasSet_arg_processing(t *testing.T) {
	initBlankContext("", "OWNER/REPO", "trunk")
	cases := []struct {
		Cmd                string
		ExpectedOutputLine string
		ExpectedConfigLine string
	}{
		{"alias set co pr checkout", "- Adding alias for co: pr checkout", "co: pr checkout"},

		{`alias set il "issue list"`, "- Adding alias for il: issue list", "il: issue list"},

		{`alias set iz 'issue list'`, "- Adding alias for iz: issue list", "iz: issue list"},

		{`alias set iy issue list --author=\$1 --label=\$2`,
			`- Adding alias for iy: issue list --author=\$1 --label=\$2`,
			`iy: issue list --author=\$1 --label=\$2`},

		{`alias set ii 'issue list --author="$1" --label="$2"'`,
			`- Adding alias for ii: issue list --author="\$1" --label="\$2"`,
			`ii: issue list --author="\$1" --label="\$2"`},

		{`alias set ix issue list --author='$1' --label='$2'`,
			`- Adding alias for ix: issue list --author=\$1 --label=\$2`,
			`ix: issue list --author=\$1 --label=\$2`},
	}

	var buf *bytes.Buffer
	for _, c := range cases {
		buf = bytes.NewBufferString("")
		defer config.StubWriteConfig(buf)()
		output, err := RunCommand(c.Cmd)
		if err != nil {
			t.Fatalf("got unexpected error running %s: %s", c.Cmd, err)
		}

		test.ExpectLines(t, output.String(), c.ExpectedOutputLine)
		test.ExpectLines(t, buf.String(), c.ExpectedConfigLine)
	}
}

func TestAliasSet_init_alias_cfg(t *testing.T) {
	cfg := `---
hosts:
  github.com:
    user: OWNER
    oauth_token: token123
`
	initBlankContext(cfg, "OWNER/REPO", "trunk")

	buf := bytes.NewBufferString("")
	defer config.StubWriteConfig(buf)()

	output, err := RunCommand("alias set diff pr diff")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	expected := `hosts:
    github.com:
        user: OWNER
        oauth_token: token123
aliases:
    diff: pr diff
`

	test.ExpectLines(t, output.String(), "Adding alias for diff: pr diff", "Added alias.")
	eq(t, buf.String(), expected)
}

func TestAliasSet_existing_aliases(t *testing.T) {
	cfg := `---
hosts:
  github.com:
    user: OWNER
    oauth_token: token123
aliases:
    foo: bar
`
	initBlankContext(cfg, "OWNER/REPO", "trunk")

	buf := bytes.NewBufferString("")
	defer config.StubWriteConfig(buf)()

	output, err := RunCommand("alias set view pr view")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	expected := `hosts:
    github.com:
        user: OWNER
        oauth_token: token123
aliases:
    foo: bar
    view: pr view
`

	test.ExpectLines(t, output.String(), "Adding alias for view: pr view", "Added alias.")
	eq(t, buf.String(), expected)

}

func TestExpandAlias(t *testing.T) {
	cfg := `---
hosts:
  github.com:
    user: OWNER
    oauth_token: token123
aliases:
  co: pr checkout
  il: issue list --author="$1" --label="$2"
`
	initBlankContext(cfg, "OWNER/REPO", "trunk")
	for _, c := range []struct {
		Args         string
		ExpectedArgs []string
		Err          string
	}{
		{"gh co", []string{"pr", "checkout"}, ""},
		{"gh il", nil, `not enough arguments for alias: issue list --author="$1" --label="$2"`},
		{"gh il vilmibm", nil, `not enough arguments for alias: issue list --author="vilmibm" --label="$2"`},
		{"gh il vilmibm epic", []string{"issue", "list", `--author=vilmibm`, `--label=epic`}, ""},
		{"gh pr status", []string{"pr", "status"}, ""},
		{"gh dne", []string{"dne"}, ""},
		{"gh", []string{}, ""},
	} {
		out, err := ExpandAlias(strings.Split(c.Args, " "))
		if err == nil && c.Err != "" {
			t.Errorf("expected error %s for %s", c.Err, c.Args)
			continue
		}

		if err != nil {
			eq(t, err.Error(), c.Err)
			continue
		}

		eq(t, out, c.ExpectedArgs)
	}
}

func TestAliasSet_invalid_command(t *testing.T) {
	initBlankContext("", "OWNER/REPO", "trunk")
	_, err := RunCommand("alias set co pe checkout")
	if err == nil {
		t.Fatal("expected error")
	}

	eq(t, err.Error(), "could not create alias: pe checkout does not correspond to a gh command")
}
