package publisher

import (
	"fmt"
	"strings"

	"github.com/dbarney/professor/types"
)

func wrapWithMarkdown(body string, status types.Status) string {
	switch status {
	case types.Success:
		return addSuccessMarkdown(body)
	case types.Failure:
		return addFailureMarkdown(body)
	case types.Error:
		return addErrorMarkdown(body)
	}
	return ""
}

func addSuccessMarkdown(body string) string {
	return fmt.Sprintf(`# SUCCESS
### complete output
total time %v seconds

%v
%v
%v`, 1, "```", body, "```")
}

func addFailureMarkdown(body string) string {
	shortBody := failMatch.FindAllString(body, -1)
	return fmt.Sprintf(`# FAILED
tldr;
%v
# grep 'fail'
%v
%v

### complete output
total time %v seconds

%v
%v
%v`, "```", strings.Join(shortBody, "\n"), "```", 1, "```", body, "```")
}

func addErrorMarkdown(body string) string {
	return fmt.Sprintf(`# ERRORED OUT
tldr;

something unexpected happened

### complete output
total time %v seconds

%v
%v
%v`, 1, "```", body, "```")
}
