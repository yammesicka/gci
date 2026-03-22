package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yammesicka/gci/pkg/section"
)

// the custom sections sort alphabetically as default.
func TestParseOrder(t *testing.T) {
	cfg := YamlConfig{
		SectionStrings: []string{"default", "prefix(github/yammesicka/gci)", "prefix(github/yammesicka/gai)"},
	}
	gciCfg, err := cfg.Parse()
	assert.NoError(t, err)
	assert.Equal(t, section.SectionList{section.Default{}, section.Custom{Prefix: "github/yammesicka/gai"}, section.Custom{Prefix: "github/yammesicka/gci"}}, gciCfg.Sections)
}

func TestParseCustomOrder(t *testing.T) {
	cfg := YamlConfig{
		SectionStrings: []string{"default", "prefix(github/yammesicka/gci)", "prefix(github/yammesicka/gai)"},
		Cfg: BoolConfig{
			CustomOrder: true,
		},
	}
	gciCfg, err := cfg.Parse()
	assert.NoError(t, err)
	assert.Equal(t, section.SectionList{section.Default{}, section.Custom{Prefix: "github/yammesicka/gci"}, section.Custom{Prefix: "github/yammesicka/gai"}}, gciCfg.Sections)
}

func TestParseNoLexOrder(t *testing.T) {
	cfg := YamlConfig{
		SectionStrings: []string{"prefix(github/yammesicka/gci)", "prefix(github/yammesicka/gai)", "default"},
		Cfg: BoolConfig{
			NoLexOrder: true,
		},
	}

	gciCfg, err := cfg.Parse()
	assert.NoError(t, err)
	assert.Equal(t, section.SectionList{section.Default{}, section.Custom{Prefix: "github/yammesicka/gci"}, section.Custom{Prefix: "github/yammesicka/gai"}}, gciCfg.Sections)
}
