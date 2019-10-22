package publisher

import (
	"bytes"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/dbarney/professor/types"
)

var (
	funcs = template.FuncMap{
		"Summarize": summarize,
	}

	success = template.Must(template.New("success.md").Funcs(funcs).Parse(string(MustAsset("success.md"))))
	failure = template.Must(template.New("failure.md").Funcs(funcs).Parse(string(MustAsset("failure.md"))))
	errored = template.Must(template.New("error.md").Funcs(funcs).Parse(string(MustAsset("error.md"))))

	failMatch = regexp.MustCompile("(?i)(.*fail.*)")
)

// Result is used to pass data to the templates
type Result struct {
	Output   string
	Duration time.Duration
}

// WrapWithMarkdown converts the output into a markdown page
func WrapWithMarkdown(res Result, status types.Status) (string, error) {
	buf := &bytes.Buffer{}
	var err error
	switch status {
	case types.Success:
		err = success.Execute(buf, res)
	case types.Failure:
		err = failure.Execute(buf, res)
	case types.Error:
		err = errored.Execute(buf, res)
	}
	return string(buf.Bytes()), err
}

// it would be really nice if this was a bit smarter.
// like actually detect which lines were important using
// a smart ranking function of some sort
func summarize(s string) string {
	summarized := failMatch.FindAllString(s, -1)
	return strings.Join(summarized, "\n")
}
