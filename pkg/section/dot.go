package section

import (
	"github.com/yammesicka/gci/pkg/parse"
	"github.com/yammesicka/gci/pkg/specificity"
)

type Dot struct{}

const DotType = "dot"

func (d Dot) MatchSpecificity(spec *parse.GciImports) specificity.MatchSpecificity {
	if spec.Name == "." {
		return specificity.NameMatch{}
	}
	return specificity.MisMatch{}
}

func (d Dot) String() string {
	return DotType
}

func (d Dot) Type() string {
	return DotType
}
