package section

import (
	"testing"

	"github.com/yammesicka/gci/pkg/specificity"
)

func TestNewLineSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{`""`, NewLine{}, specificity.MisMatch{}},
		{`"x"`, NewLine{}, specificity.MisMatch{}},
		{`"\n"`, NewLine{}, specificity.MisMatch{}},
	}
	testSpecificity(t, testCases)
}

// func TestNewLineToString(t *testing.T) {
// 	testSectionToString(t, NewLine{})
// }
