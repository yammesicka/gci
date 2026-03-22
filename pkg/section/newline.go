package section

import (
	"github.com/yammesicka/gci/pkg/parse"
	"github.com/yammesicka/gci/pkg/specificity"
)

const newLineName = "newline"

type NewLine struct{}

func (n NewLine) MatchSpecificity(spec *parse.GciImports) specificity.MatchSpecificity {
	return specificity.MisMatch{}
}

func (n NewLine) String() string {
	return newLineName
}

func (n NewLine) Type() string {
	return newLineName
}
