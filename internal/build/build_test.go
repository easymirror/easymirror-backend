package build

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReturnIfNotEmpty(t *testing.T) {
	tests := []struct {
		Expected string
		Input    string
	}{
		{Expected: "appVersion", Input: "appVersion"},
		{Expected: "N/A", Input: ""},
		{Expected: "N/A", Input: "   "},
		{Expected: "N/A", Input: "       "},
	}
	for testNum, test := range tests {
		t.Run(fmt.Sprintf("Test #%v", testNum), func(t *testing.T) {
			result := returnIfNotEmpty(test.Input)
			assert.Equal(t, test.Expected, result)
		})
	}

}
