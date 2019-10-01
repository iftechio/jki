package utils

import (
	"reflect"
	"strings"
	"testing"
)

func TestExtractBaseImages(t *testing.T) {
	t.Parallel()
	dockerfile := `
FROM w1/w2/w3/foo:v1.2.3 AS builder # qweasd
EXPOSE 1000
RUN echo 1

  FROM alpine:3.8 #foo
COPY --from=builder`
	expected := []string{"w1/w2/w3/foo:v1.2.3", "alpine:3.8"}

	got, err := ExtractBaseImages(strings.NewReader(dockerfile))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("expected: %#v, got: %#v", expected, got)
	}
}
