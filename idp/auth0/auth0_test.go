package auth0

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_extractID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "abc",
			args: args{id: "abc"},
			want: "abc",
		},
		{
			name: "abc|123",
			args: args{id: "abc|123"},
			want: "123",
		},
		{
			name: "abc|123|456",
			args: args{id: "abc|123|456"},
			want: "456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractID(tt.args.id); got != tt.want {
				t.Errorf("extractID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMapValue(t *testing.T) {
	values := map[string]interface{}{
		"name":                           "nameA",
		"https://openocean.studio/scope": "scpA",
		"invalidType":                    123,
	}

	tests := []struct {
		name          string
		field         string
		expectedValue string
	}{
		{
			name:          "valid case",
			field:         "name",
			expectedValue: "nameA",
		},
		{
			name:          "valid scope case",
			field:         "https://openocean.studio/scope",
			expectedValue: "scpA",
		},
		{
			name:          "invalid type case",
			field:         "invalidType",
			expectedValue: "",
		},
		{
			name:          "not found field case",
			field:         "missingField",
			expectedValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualValue := getMapValue(tt.field, values)

			assert.Equal(t, tt.expectedValue, actualValue)
		})
	}
}

func Test_findID(t *testing.T) {
	tests := []struct {
		name          string
		values        map[string]interface{}
		expectedValue string
	}{
		{
			name: "federateduserid case",
			values: map[string]interface{}{
				"sub": "waad|abc",
				"https://openocean.studio/federateduserid": "123",
			},
			expectedValue: "123",
		},
		{
			name: "sub case",
			values: map[string]interface{}{
				"sub": "waad|abc",
			},
			expectedValue: "abc",
		},
		{
			name:          "no ID case",
			values:        map[string]interface{}{},
			expectedValue: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualValue := findID(tt.values)

			assert.Equal(t, tt.expectedValue, actualValue)
		})
	}
}
