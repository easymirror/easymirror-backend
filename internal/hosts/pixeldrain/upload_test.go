package pixeldrain

import (
	"context"
	"fmt"
	"testing"
)

// go test -v -timeout 30s -run ^TestUpload$ github.com/easymirror/easymirror-backend/internal/hosts/pixeldrain
func TestUpload(t *testing.T) {
	id, err := Upload(context.Background(), "https://easymirror.s3.us-east-1.amazonaws.com/48751a99-dc86-4e08-a885-ef0f75337779/mattaio.png?response-content-disposition=inline&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEFAaCXVzLWVhc3QtMSJIMEYCIQC%2FJTUzb2R8wRHugBa9rgWnjAMr3m8Z3Xx4aRa8KgIHRgIhAPBpuQY1%2FcbkjiJA83Jm2N1PrBwj6kROie28df2Wsl3yKuQCCEkQABoMODQ5Njc5NTY1Mzg3IgxowpOfNjZbYjBXMXEqwQIVw31XVc3%2Bvn0d4QwO60%2FVOOoFHY02%2BOiDp8ngCoRzUlJYUXCfvb56tb8wcpCb8q4mCFvwLW5t6xi8ZWVABDpCRXkAc%2FdiH%2FlhnWLpifIHNwjKUYmQj9lIeJUpXcMKKC4l1CFYt5dF%2F7Bco%2FEHUT4s8QgzwjeDlfvoU9G0u2G6qDjxyHaG%2BoYLfkOrV99yn9kdXpXeKkgLX7h0L5qvs9uy1UoXk291wkCMDAiEC2VmBUhzlEEjoQJNCFZAAU95if3i0wMlNppLzFyFqLoRGEr%2B6A8vzXZgNoBupIGG4HfOAvP5c7zpuerVak6guYdg4ToUmh7uYdXVckGozm4WMZQt5btGfJD9H1JGrOR%2B%2FJWZu9TUqKApvv6NwnapgKKzaDSdKdIvoH%2FW9J63VgVrSnxvJZ9GCxZHmOkTF1VdpSf1f00w15GNrwY6sgKfYjZoTGcLsXC627mH3VdaPSZdyU%2BPszqmoBK5XtJ0JbOZz3pVontXbWKFXNNry%2F6oYvQCNOPvL7qA6XIFdAMJ83STj7LQGJHb3pJQbfcwE%2B9zDR9oGGBta%2FYsJF3MBTbwpyeEuOXAuu8%2Fe6ELrJddv9JNq6pEbc9imvx0OeLwZBy8vnbHFaFMUs84qhWQAVkrEa%2BYtiaVMdZgSXTxzYDxv%2Ff5LkWvEi%2Firrj67tXf%2FIL3kg%2FEvc9xfIYHp5CilSpAPBb2d%2BI4kity52s1y45oRTxOZy4cKNCnIeYzHV6c%2BJqgXb4ON1roPHd1DKdFoGciEFXMj161MKcnqHfu%2BCU2kS5rsE%2Bf9Op0xLMWyLSMfqtDzbO7Du2FfyRuoocW6YBHOGxZ5ACSg5ZHNP1HX56Kt80%3D&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Date=20240302T212653Z&X-Amz-SignedHeaders=host&X-Amz-Expires=300&X-Amz-Credential=ASIA4LVGZIJFR5HEZ3FF%2F20240302%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Signature=be6c9d446d0e20bf6730f125363462d04332bca3fa8ed2f6bc7b9477a98d62e6")
	if err != nil {
		t.Fatalf("Error uploading: %v", err)
	}
	fmt.Println(id)

}
