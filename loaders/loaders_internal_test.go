package loaders

import (
	"github.com/google/go-cmp/cmp"
	"io"
	"strings"
	"testing"
)

func TestLinesFromReader(t *testing.T) {
	tests := map[string]struct {
		input io.Reader
		want  []string
	}{
		"IPaddr":       {input: strings.NewReader("8.8.8.8"), want: []string{"8.8.8.8"}},
		"http":         {input: strings.NewReader("http://host/url"), want: []string{"host"}},
		"https":        {input: strings.NewReader("https://host/url"), want: []string{"host"}},
		"host:port":    {input: strings.NewReader("google.com:80"), want: []string{"google.com:80"}},
		"commentafter": {input: strings.NewReader("google.com # some comment"), want: []string{"google.com"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, _ := linesFromReader(tc.input)
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
