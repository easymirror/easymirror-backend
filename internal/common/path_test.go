package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilenameFromURI(t *testing.T) {
	tests := []struct {
		URL, Expected string
	}{
		{URL: "https://easymirror.s3.us-east-1.amazonaws.com/43bb847b-d89c-4db4-803a-522a0148b8d0/Matthew_Lugo_Resume.pdf?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA4LVGZIJF45Z3KRFR%2F20240303%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240303T025927Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&x-id=GetObject&X-Amz-Signature=31720faf643cd52ef0385574ec1563615bc1f7213c18777bd14e6c259582f704", Expected: "Matthew_Lugo_Resume.pdf"},
	}

	for testNum, test := range tests {
		t.Run(fmt.Sprintf("Test #%v", testNum), func(t *testing.T) {
			result, err := FilenameFromURI(test.URL)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, test.Expected, result)
		})
	}

}
