package domain

import (
	"testing"
)

func TestGetResponseType(t *testing.T) {
	testcases := []struct {
		name          string
		input         string
		wantErr       bool
		expectedError string
	}{
		{
			name:          "empty string",
			input:         "",
			wantErr:       true,
			expectedError: "unsupported response_type: ",
		},
		{
			name:    "supported response_type code",
			input:   "code",
			wantErr: false,
		},
		{
			name:          "unsupported response_type",
			input:         "token",
			wantErr:       true,
			expectedError: "unsupported response_type: token",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetResponseType(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for input '%s', got nil", tt.input)
				} else if err.Error() != tt.expectedError {
					t.Errorf("expected error message '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error for input '%s', got '%s'", tt.input, err.Error())
				}
			}
		})
	}
}
