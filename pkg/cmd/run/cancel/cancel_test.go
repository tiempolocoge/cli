package cancel

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/cli/cli/internal/ghrepo"
	"github.com/cli/cli/pkg/cmdutil"
	"github.com/cli/cli/pkg/httpmock"
	"github.com/cli/cli/pkg/iostreams"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdCancel(t *testing.T) {
	tests := []struct {
		name     string
		cli      string
		wants    CancelOptions
		wantsErr bool
	}{
		{
			name:     "blank",
			wantsErr: true,
		},
		{
			name: "with arg",
			cli:  "1234",
			wants: CancelOptions{
				RunID: "1234",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io, _, _, _ := iostreams.Test()
			io.SetStdinTTY(true)
			io.SetStdoutTTY(true)

			f := &cmdutil.Factory{
				IOStreams: io,
			}

			argv, err := shlex.Split(tt.cli)
			assert.NoError(t, err)

			var gotOpts *CancelOptions
			cmd := NewCmdCancel(f, func(opts *CancelOptions) error {
				gotOpts = opts
				return nil
			})

			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(ioutil.Discard)
			cmd.SetErr(ioutil.Discard)

			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, tt.wants.RunID, gotOpts.RunID)
		})
	}
}

func TestRunCancel(t *testing.T) {
	tests := []struct {
		name      string
		httpStubs func(*httpmock.Registry)
		opts      *CancelOptions
		wantErr   bool
		wantOut   string
		errMsg    string
	}{
		{
			name: "cancel run",
			opts: &CancelOptions{
				RunID: "1234",
			},
			wantErr: false,
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("POST", "repos/OWNER/REPO/actions/runs/1234/cancel"),
					httpmock.StatusStringResponse(202, "{}"),
				)
			},
			wantOut: "✓ You have successfully requested the workflow to be canceled.",
		},
		{
			name: "not found",
			opts: &CancelOptions{
				RunID: "1234",
			},
			wantErr: true,
			errMsg:  "Could not find any workflow run with ID 1234",
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("POST", "repos/OWNER/REPO/actions/runs/1234/cancel"),
					httpmock.StatusStringResponse(404, ""),
				)
			},
		},
		{
			name: "completed",
			opts: &CancelOptions{
				RunID: "1234",
			},
			wantErr: true,
			errMsg:  "Cannot cancel a workflow run that is completed",
			httpStubs: func(reg *httpmock.Registry) {
				reg.Register(
					httpmock.REST("POST", "repos/OWNER/REPO/actions/runs/1234/cancel"),
					httpmock.StatusStringResponse(409, ""),
				)
			},
		},
	}

	for _, tt := range tests {
		reg := &httpmock.Registry{}
		tt.httpStubs(reg)
		tt.opts.HttpClient = func() (*http.Client, error) {
			return &http.Client{Transport: reg}, nil
		}

		io, _, stdout, _ := iostreams.Test()
		io.SetStdoutTTY(true)
		tt.opts.IO = io
		tt.opts.BaseRepo = func() (ghrepo.Interface, error) {
			return ghrepo.FromFullName("OWNER/REPO")
		}

		t.Run(tt.name, func(t *testing.T) {
			err := runCancel(tt.opts)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Equal(t, tt.errMsg, err.Error())
				}
			}
			assert.Equal(t, tt.wantOut, stdout.String())
			reg.Verify(t)
		})
	}
}
