package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -timeout 30s -run ^TestGetPageOffset$ github.com/easymirror/easymirror-backend/internal/common
func TestGetPageOffset(t *testing.T) {
	tests := []struct {
		Limit, PageNum int
		Expected       int
	}{
		{Limit: 25, PageNum: 0, Expected: 0},
		{Limit: 25, PageNum: 1, Expected: 0},
		{Limit: 25, PageNum: 2, Expected: 25},
		{Limit: 25, PageNum: 3, Expected: 50},
		{Limit: 10, PageNum: 0, Expected: 0},
		{Limit: 10, PageNum: 1, Expected: 0},
		{Limit: 10, PageNum: 2, Expected: 10},
		{Limit: 10, PageNum: 3, Expected: 20},
	}

	for testNum, test := range tests {
		t.Run(fmt.Sprintf("Test #%v", testNum), func(t *testing.T) {
			result := GetPageOffset(test.Limit, test.PageNum)
			assert.Equal(t, test.Expected, result)
		})
	}
}
