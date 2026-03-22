package config

import (
	"reflect"
	"testing"

	"github.com/yammesicka/gci/v2/pkg/section"
)

func TestParseOrder(t *testing.T) {
	cfg := YamlConfig{
		SectionStrings: []string{"default", "prefix(github/yammesicka/gci)", "prefix(github/yammesicka/gai)"},
	}
	gciCfg, err := cfg.Parse()
	if err != nil {
		t.Fatal(err)
	}
	want := section.SectionList{
		section.Default{},
		section.Custom{Prefix: "github/yammesicka/gai"},
		section.Custom{Prefix: "github/yammesicka/gci"},
	}
	if !reflect.DeepEqual(want, gciCfg.Sections) {
		t.Fatalf("unexpected sections: got=%v want=%v", gciCfg.Sections, want)
	}
}

func TestParseCustomOrder(t *testing.T) {
	cfg := YamlConfig{
		SectionStrings: []string{"default", "prefix(github/yammesicka/gci)", "prefix(github/yammesicka/gai)"},
		Cfg: BoolConfig{
			CustomOrder: true,
		},
	}
	gciCfg, err := cfg.Parse()
	if err != nil {
		t.Fatal(err)
	}
	want := section.SectionList{
		section.Default{},
		section.Custom{Prefix: "github/yammesicka/gci"},
		section.Custom{Prefix: "github/yammesicka/gai"},
	}
	if !reflect.DeepEqual(want, gciCfg.Sections) {
		t.Fatalf("unexpected sections: got=%v want=%v", gciCfg.Sections, want)
	}
}

func TestParseNoLexOrder(t *testing.T) {
	cfg := YamlConfig{
		SectionStrings: []string{"prefix(github/yammesicka/gci)", "prefix(github/yammesicka/gai)", "default"},
		Cfg: BoolConfig{
			NoLexOrder: true,
		},
	}
	gciCfg, err := cfg.Parse()
	if err != nil {
		t.Fatal(err)
	}
	want := section.SectionList{
		section.Default{},
		section.Custom{Prefix: "github/yammesicka/gci"},
		section.Custom{Prefix: "github/yammesicka/gai"},
	}
	if !reflect.DeepEqual(want, gciCfg.Sections) {
		t.Fatalf("unexpected sections: got=%v want=%v", gciCfg.Sections, want)
	}
}
