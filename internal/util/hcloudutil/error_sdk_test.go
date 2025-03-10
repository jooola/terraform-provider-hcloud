package hcloudutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestErrorToDiag(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		wantedSummary string
	}{
		{
			"basic error",
			errors.New("basic error"),
			"basic error",
		},
		{
			"hcloud generic error",
			hcloud.Error{Code: hcloud.ErrorCodeServiceError, Message: "Service Error"},
			"Service Error (service_error)",
		},
		{
			"hcloud invalid input",
			hcloud.Error{Code: hcloud.ErrorCodeInvalidInput, Message: "Invalid Input", Details: hcloud.ErrorDetailsInvalidInput{Fields: []hcloud.ErrorDetailsInvalidInputField{{Name: "ip", Messages: []string{"invalid field"}}}}},
			"Invalid Input (invalid_input): [ip => [invalid field]]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ErrorToDiag(tt.err)
			if !got.HasError() {
				t.Fatal("Expected to get errors")
			}
			assert.Equal(t, tt.wantedSummary, got[0].Summary)
		})
	}
}
