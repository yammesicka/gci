package section

import (
	"github.com/yammesicka/gci/pkg/parse"
	"github.com/yammesicka/gci/pkg/specificity"
)

const DefaultType = "default"

type Default struct{}

func (d Default) MatchSpecificity(spec *parse.GciImports) specificity.MatchSpecificity {
	return specificity.Default{}
}

func (d Default) String() string {
	return DefaultType
}

func (d Default) Type() string {
	return DefaultType
}
