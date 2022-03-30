package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclwrite"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/minamijoyo/tfmigrate/config"
	"github.com/minamijoyo/tfmigrate/tfmigrate"
)

func main() {
	if err := run(time.Now().Format("20060102150405_tfmigration"), os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(name string, r io.Reader, w io.Writer) error {
	var tfplan tfjson.Plan
	if err := json.NewDecoder(r).Decode(&tfplan); err != nil {
		return err
	}

	migration := config.MigrationFile{
		Migration: config.MigrationBlock{
			Type: "state",
			Name: name,
		},
	}
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(&migration, f.Body())

	gohcl.EncodeIntoBody(&tfmigrate.StateMigratorConfig{
		Actions: makeMvActions(tfplan.ResourceChanges),
	}, f.Body().Blocks()[0].Body())

	_, err := f.WriteTo(w)
	return err
}

func makeMvActions(changes []*tfjson.ResourceChange) []string {
	var creates, deletes []*tfjson.ResourceChange
	for _, rc := range changes {
		switch {
		case rc.Change.Actions.Create() && rc.Change.Before == nil:
			creates = append(creates, rc)
		case rc.Change.Actions.Delete():
			deletes = append(deletes, rc)
		}
	}

	var actions []string
	for _, c := range creates {
		for _, d := range deletes {
			if eq(c, d) {
				actions = append(actions, fmt.Sprintf("mv %q %q", d.Address, c.Address))
			}
		}
	}
	return actions
}

func eq(c, d *tfjson.ResourceChange) bool {
	if c.Type != d.Type {
		return false
	}
	if c.ProviderName != d.ProviderName {
		return false
	}
	before := d.Change.Before.(map[string]interface{})
	after := c.Change.After.(map[string]interface{})
	for k, a := range after {
		if b, ok := before[k]; ok && reflect.DeepEqual(b, a) {
			return true
		}
	}
	return false
}
