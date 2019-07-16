package publisher

import (
	"testing"

	"github.com/dbarney/professor/types"
)

func TestMarkdownRendering(t *testing.T) {
	res, err := wrapWithMarkdown(Result{Output: "this is the body"}, types.Success)
	if err != nil {
		t.Fatalf("unable to render body %v", err)
	}
	if len(res) == 0 {
		t.Fatalf("wrong body was created")
	}

	res, err = wrapWithMarkdown(Result{Output: "this is the body"}, types.Failure)
	if err != nil {
		t.Fatalf("unable to render body %v", err)
	}
	if len(res) == 0 {
		t.Fatalf("wrong body was created")
	}

	res, err = wrapWithMarkdown(Result{Output: "this is the body"}, types.Error)
	if err != nil {
		t.Fatalf("unable to render body %v", err)
	}
	if len(res) == 0 {
		t.Fatalf("wrong body was created")
	}
}
