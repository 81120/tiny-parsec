package ini_test

import (
	"testing"

	"github.com/81120/tiny-parsec/ini"
	"github.com/stretchr/testify/assert"
)

func TestISectionName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		err      bool
	}{
		{"valid section", "[database]", "database", false},
		{"with spaces", "[  redis  ]", "redis", false},
		{"missing open bracket", "database]", "", true},
		{"missing close bracket", "[database", "", true},
		{"empty section", "[]", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ini.ISectionName().Parse(tt.input)
			if tt.err {
				assert.True(t, result.IsNothing())
			} else {
				assert.True(t, result.IsJust())
				assert.Equal(t, tt.expected, result.Get().First)
			}
		})
	}
}

func TestIniParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ini.Ini
		err      bool
	}{
		{
			"basic structure",
			"[section]\nkey=value",
			ini.Ini{
				Sections: []ini.Section{{
					Name:    "section",
					Entries: []ini.Entry{{Key: "key", Value: "value"}},
				}},
			},
			false,
		},
		{
			"multiple sections",
			"[db]\nhost=localhost\n[cache]\nport=6379",
			ini.Ini{
				Sections: []ini.Section{
					{Name: "db", Entries: []ini.Entry{
						{Key: "host", Value: "localhost"},
					}},
					{Name: "cache", Entries: []ini.Entry{
						{Key: "port", Value: "6379"},
					}},
				},
			},
			false,
		},
		{
			"ignore comments",
			"; comment\n[section]\n# another comment\nkey=value",
			ini.Ini{
				Sections: []ini.Section{{
					Name:    "section",
					Entries: []ini.Entry{{Key: "key", Value: "value"}},
				}},
			},
			false,
		},
		{
			"invalid key format",
			"[section]\nkeyvalue",
			ini.Ini{
				Sections: []ini.Section{{
					Name:    "section",
					Entries: []ini.Entry{},
				}},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ini.IniParse().Parse(tt.input)
			if tt.err {
				assert.True(t, result.IsNothing())
			} else {
				assert.True(t, result.IsJust())
				assert.Equal(t, tt.expected, result.Get().First)
			}
		})
	}
}
