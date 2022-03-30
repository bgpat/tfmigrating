package main

import (
	"bytes"
	_ "embed"
	"testing"
)

var (
	//go:embed test-fixtures/tfplan.json
	tfplan []byte

	//go:embed test-fixtures/migration.hcl
	want string
)

func TestRun(t *testing.T) {
	var buf bytes.Buffer
	if err := run("20220331032518_tfmigration", bytes.NewReader(tfplan), &buf); err != nil {
		t.Fatalf("%s", err)
	}
	if got := buf.String(); got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}
