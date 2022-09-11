package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_identify(t *testing.T) {
	tcc := map[string]struct {
		in       any
		dedupe   bool
		merge    bool
		expected any
	}{
		"number": {
			in:       23.42,
			expected: "number",
		},
		"string": {
			in:       "hans",
			expected: "string",
		},
		"empty string": {
			in:       "",
			expected: "string",
		},
		"true": {
			in:       true,
			expected: "bool",
		},
		"false": {
			in:       false,
			expected: "bool",
		},
		"nil": {
			in:       nil,
			expected: nil,
		},
		"empty array": {
			in:       []any{},
			expected: []any{},
		},
		"array of numbers": {
			in:       []any{1, 1.2, -3, 4},
			expected: []any{"number", "number", "number", "number"},
		},
		"array of strings": {
			in:       []any{"hans", "", "peter"},
			expected: []any{"string", "string", "string"},
		},
		"array of any": {
			in:       []any{"hans", nil, -23.22, nil, 2, true, []any{"peter"}, map[string]any{"key": 42.222}},
			expected: []any{"string", nil, "number", nil, "number", "bool", []any{"string"}, map[string]any{"key": "number"}},
		},
		"array of any - dedupe": {
			in:       []any{"hans", nil, -23.22, nil, 2, true, []any{"peter"}, "hans", "peter", map[string]any{"key": 42.222}},
			dedupe:   true,
			expected: []any{"string", nil, "number", "bool", []any{"string"}, map[string]any{"key": "number"}},
		},
		"array of any - merge objects": {
			in: []any{
				23,
				map[string]any{"key1": 42},
				map[string]any{"key1": "hans", "key2": true},
				[]any{nil},
				map[string]any{"key3": map[string]any{"a": false}},
				true,
				46,
				map[string]any{"key3": map[string]any{"b": 23}},
			},
			merge: true,
			expected: []any{
				"number",
				map[string]any{
					"key1": "any",
					"key2": "bool",
					"key3": map[string]any{"a": "bool", "b": "number"},
				},
				[]any{nil},
				"bool",
				"number",
			},
		},
		"array of any - merge and dedupe": {
			in: []any{
				23,
				map[string]any{"key1": 42},
				map[string]any{"key1": "hans", "key2": true},
				[]any{nil},
				map[string]any{"key3": map[string]any{"a": false}},
				true,
				46,
				map[string]any{"key3": map[string]any{"b": 23}},
				false,
			},
			merge:  true,
			dedupe: true,
			expected: []any{
				"number",
				map[string]any{
					"key1": "any",
					"key2": "bool",
					"key3": map[string]any{"a": "bool", "b": "number"},
				},
				[]any{nil},
				"bool",
			},
		},
		"empty map": {
			in:       map[string]any{},
			expected: map[string]any{},
		},
		"map to string": {
			in:       map[string]any{"hans": "peter"},
			expected: map[string]any{"hans": "string"},
		},
		"map to any": {
			in: map[string]any{
				"name":     "peter",
				"age":      23,
				"colors":   []any{"red", "blue"},
				"verified": true,
				"email":    nil,
			},
			expected: map[string]any{
				"name":     "string",
				"age":      "number",
				"colors":   []any{"string", "string"},
				"verified": "bool",
				"email":    nil,
			},
		},
	}

	for tn, tc := range tcc {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			actual := identify(tc.in, tc.dedupe, tc.merge)
			assert.Equal(t, tc.expected, actual)
		})
	}

}
