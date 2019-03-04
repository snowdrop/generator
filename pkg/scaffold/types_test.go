package scaffold

import (
	"reflect"
	"testing"
)

func TestIsModuleAvailable(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		module   Module
		expected bool
	}{
		{
			name:    "Module should be available if there is no availability",
			version: "2.1.2",
			module: Module{
				Name: "foo",
			},
			expected: true,
		}, {
			name:    "Module should be available if there is no availability, with .RELEASE version",
			version: "2.1.2.RELEASE",
			module: Module{
				Name: "foo",
			},
			expected: true,
		}, {
			name:    "Module should not be available if given an invalid version",
			version: "2.1EASE",
			module: Module{
				Name: "foo",
			},
			expected: false,
		}, {
			name:    "Module should not be available if availability is incorrect",
			version: "2.1.2",
			module: Module{
				Name:         "foo",
				Availability: "bar",
			},
			expected: false,
		}, {
			name:    "Module should not be available if availability is incorrect, invalid input version",
			version: "2.1EASE",
			module: Module{
				Name:         "foo",
				Availability: "2.1.2.RELEASE",
			},
			expected: false,
		}, {
			name:    "Availability requires full version or placeholder",
			version: "2.1.2",
			module: Module{
				Name:         "foo",
				Availability: ">=2.1",
			},
			expected: false,
		}, {
			name:    "Module should be available if availability is correct and version is within range",
			version: "2.1.2",
			module: Module{
				Name:         "foo",
				Availability: ">=2.1.x",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		if tt.expected != tt.module.IsAvailableFor(tt.version) {
			t.Errorf("'%s' failed", tt.name)
		}
	}
}

func TestKeepModulesCompatibleWith(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		input    []Module
		expected []Module
	}{
		{
			name:    "No defined availability, .RELEASE version",
			version: "2.1.2.RELEASE",
			input: []Module{
				{
					Name: "foo",
				},
				{
					Name: "bar",
				},
			},
			expected: []Module{
				{
					Name: "foo",
				},
				{
					Name: "bar",
				},
			},
		}, {
			name:    "No defined availability, simple version",
			version: "2.1.2",
			input: []Module{
				{
					Name: "foo",
				},
				{
					Name: "bar",
				},
			},
			expected: []Module{
				{
					Name: "foo",
				},
				{
					Name: "bar",
				},
			},
		}, {
			name:    "Exact availability, simple version",
			version: "2.1.2",
			input: []Module{
				{
					Name:         "foo",
					Availability: "2.1.2",
				},
				{
					Name: "bar",
				},
			},
			expected: []Module{
				{
					Name:         "foo",
					Availability: "2.1.2",
				},
				{
					Name: "bar",
				},
			},
		}, {
			name:    "Wrong availability, module should be ignored",
			version: "2.1.2",
			input: []Module{
				{
					Name:         "foo",
					Availability: "2.1.2.RELEASE",
				},
				{
					Name: "bar",
				},
			},
			expected: []Module{
				{
					Name: "bar",
				},
			},
		}, {
			name:    "Wrong version, no modules should be available",
			version: "foo",
			input: []Module{
				{
					Name:         "foo",
					Availability: "2.1.2.RELEASE",
				},
				{
					Name: "bar",
				},
			},
			expected: []Module{},
		},
	}

	for _, tt := range tests {
		actual := keepModulesCompatibleWith(tt.input, tt.version)
		if !reflect.DeepEqual(tt.expected, actual) {
			t.Errorf("'%s' failed: expected %v, got %v", tt.name, tt.expected, actual)
		}
	}
}
