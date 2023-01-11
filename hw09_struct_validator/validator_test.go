package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
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
				ID:  "123", // valid
				Age: 100,   // invalid
				Phones: []string{
					"12345",        // valid
					"123456789012", // invalid
				},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   errors.New("validation failed on 'max:50'"),
				},
				ValidationError{
					Field: "Email",
					Err:   errors.New("validation failed on 'regexp:^\\w+@\\w+\\.\\w+$'"),
				},
				ValidationError{
					Field: "Phones",
					Err:   errors.New("validation failed on 'len:11'"),
				},
			},
		},
		{
			in: User{
				ID:    "123",           // valid
				Age:   10,              // invalid
				Email: "test@mail.com", // valid
				Phones: []string{
					"12345", // valid
				},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Age",
					Err:   errors.New("validation failed on 'min:18'"),
				},
			},
		},
		{
			in: App{
				Version: "12345", // valid
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: App{
				Version: "123456", // invalid
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   errors.New("validation failed on 'len:5'"),
				},
			},
		},
		{
			in: Response{
				Code: 200, // valid
			},
			expectedErr: ValidationErrors{},
		},
		{
			in: Response{
				Code: 300, // invalid
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   errors.New("validation failed on 'in:200,404,500'"),
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			require.Equal(t, tt.expectedErr, Validate(tt.in))
		})
	}
}
