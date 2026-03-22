package section

import (
	"github.com/yammesicka/gci/pkg/parse"
	"github.com/yammesicka/gci/pkg/specificity"
)

type Blank struct{}

const BlankType = "blank"

func (b Blank) MatchSpecificity(spec *parse.GciImports) specificity.MatchSpecificity {
	if spec.Name == "_" {
		return specificity.NameMatch{}
	}
	return specificity.MisMatch{}
}

func (b Blank) String() string {
	return BlankType
}

func (b Blank) Type() string {
	return BlankType
}
