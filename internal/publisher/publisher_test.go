package publisher

import (
	"testing"
)

func TestSummaryRegex(t *testing.T) {
	short := failMatch.FindAllString("failure: something went wrong", -1)
	if len(short) == 0 || short[0] == "" {
		t.Fatalf("nothing was matched on a single line")
	}
	short = failMatch.FindAllString("fAiLure: something went wrong", -1)
	if len(short) == 0 || short[0] == "" {
		t.Fatalf("case insensitive failed")
	}
}
