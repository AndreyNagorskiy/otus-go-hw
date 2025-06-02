package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     strings.Repeat("1", 36),
				Name:   "Test",
				Age:    14,
				Email:  "example@example",
				Role:   "invalid",
				Phones: []string{"78000000001", "780"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   ErrValidateMin,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrValidateRegexp,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrValidateIn,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrValidateLen,
				},
			},
		},
		{
			in: App{
				Version: "1.2.3.4",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrValidateLen,
				},
			},
		},
		{
			in: Response{
				Code: 301,
				Body: "test",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrValidateIn,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("negative case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}

	tests = []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     strings.Repeat("1", 36),
				Name:   "Test",
				Age:    20,
				Email:  "example@example.com",
				Role:   "admin",
				Phones: []string{"78000000001", "78000000002"},
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "1.2.3",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
				Body: "test",
			},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte("1"),
				Payload:   []byte("2"),
				Signature: []byte("3"),
			},
			expectedErr: nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("positive case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
		})
	}
}
