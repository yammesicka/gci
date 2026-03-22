package section

import (
	"testing"

	"github.com/yammesicka/gci/v2/pkg/specificity"
)

func TestDefaultSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{`""`, Default{}, specificity.DefaultMatch{}},
		{`"x"`, Default{}, specificity.DefaultMatch{}},
	}
	testSpecificity(t, testCases)
}
