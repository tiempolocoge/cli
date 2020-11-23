package shared

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cli/cli/pkg/iostreams"
)

func dumpPath(random int64) string {
	r := fmt.Sprintf("%x", random)
	r = r[len(r)-5:]
	dumpFilename := fmt.Sprintf("gh%s.json", r)
	return filepath.Join(os.TempDir(), dumpFilename)
}

func PreserveInput(io *iostreams.IOStreams, state *IssueMetadataState, createErr *error) func() {
	return func() {
		if !state.IsDirty() {
			return
		}

		if *createErr == nil {
			return
		}

		out := io.ErrOut

		// this extra newline guards against appending to the end of a survey line
		fmt.Fprintln(out)

		data, err := json.Marshal(state)
		if err != nil {
			fmt.Fprintf(out, "failed to save input to file: %s\n", err)
			fmt.Fprintln(out, "would have saved:")
			fmt.Fprintf(out, "%v\n", state)
			return
		}

		tmpfile, err := io.TempFile(os.TempDir(), "gh*.json")
		if err != nil {
			fmt.Fprintf(out, "failed to save input to file: %s\n", err)
			fmt.Fprintln(out, "would have saved:")
			fmt.Fprintf(out, "%v\n", state)
			return
		}

		tmpfilePath := filepath.Join(os.TempDir(), tmpfile.Name())

		_, err = tmpfile.Write(data)
		if err != nil {
			fmt.Fprintf(out, "failed to save input to file: %s\n", err)
			fmt.Fprintln(out, "would have saved:")
			fmt.Fprintln(out, string(data))
			return
		}

		cs := io.ColorScheme()

		issueType := "pr"
		if state.Type == IssueMetadata {
			issueType = "issue"
		}

		fmt.Fprintf(out, "%s operation failed. input saved to: %s\n", cs.FailureIcon(), tmpfilePath)
		fmt.Fprintf(out, "resubmit with: gh %s create -j@%s\n", issueType, tmpfilePath)

		// some whitespace before the actual error
		fmt.Fprintln(out)
	}
}
