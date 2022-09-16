package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_identify_schema(t *testing.T) {
	tcc := map[string]struct {
		in       any
		merge    bool
		expected *jSchema
	}{
		"number": {
			in:       23.42,
			expected: &jSchema{Type: newSet("number")},
		},
		"int": {
			in:       42,
			expected: &jSchema{Type: newSet("integer")},
		},
		"string": {
			in:       "peter",
			expected: &jSchema{Type: newSet("string")},
		},
		"empty string": {
			in:       "",
			expected: &jSchema{Type: newSet("string")},
		},
		"true": {
			in:       true,
			expected: &jSchema{Type: newSet("boolean")},
		},
		"false": {
			in:       false,
			expected: &jSchema{Type: newSet("boolean")},
		},
		"null": {
			in:       nil,
			expected: &jSchema{Type: newSet("null")},
		},
		"empty object": {
			in: map[string]any{},
			expected: &jSchema{
				Type:       newSet("object"),
				Properties: map[string]*jSchema{},
			},
		},
		"object with string values": {
			in: map[string]any{"hans": "peter"},
			expected: &jSchema{
				Type: newSet("object"),
				Properties: map[string]*jSchema{
					"hans": {Type: newSet("string")},
				},
			},
		},
		"object with different values": {
			in: map[string]any{
				"name":     "peter",
				"age":      23,
				"money":    42.11,
				"verified": false,
			},
			expected: &jSchema{
				Type: newSet("object"),
				Properties: map[string]*jSchema{
					"name":     {Type: newSet("string")},
					"age":      {Type: newSet("integer")},
					"money":    {Type: newSet("number")},
					"verified": {Type: newSet("boolean")},
				},
			},
		},
		"empty array": {
			in: []any{},
			expected: &jSchema{
				Type: newSet("array"),
			},
		},
		"array of strings": {
			in: []any{"red", "green", "blue"},
			expected: &jSchema{
				Type: newSet("array"),
				PrefixItems: []*jSchema{
					{Type: newSet("string")},
					{Type: newSet("string")},
					{Type: newSet("string")},
				},
			},
		},
		"array of strings - merge": {
			in:    []any{"red", "green", "blue"},
			merge: true,
			expected: &jSchema{
				Type:     newSet("array"),
				Contains: &jSchema{Type: newSet("string")},
			},
		},
		"array - merge": {
			in:    []any{"red", true, nil},
			merge: true,
			expected: &jSchema{
				Type:     newSet("array"),
				Contains: &jSchema{Type: newSet("string", "boolean", "null")},
			},
		},
		"array - nested": {
			in: []any{[]any{"a"}, []any{}, []any{"b"}},
			expected: &jSchema{
				Type: newSet("array"),
				PrefixItems: []*jSchema{
					{
						Type:        newSet("array"),
						PrefixItems: []*jSchema{{Type: newSet("string")}},
					},
					{
						Type: newSet("array"),
					},
					{
						Type:        newSet("array"),
						PrefixItems: []*jSchema{{Type: newSet("string")}},
					},
				},
			},
		},
		"array - nested - merge": {
			in:    []any{[]any{"a"}, []any{}, []any{"b"}},
			merge: true,
			expected: &jSchema{
				Type: newSet("array"),
				Contains: &jSchema{
					Type:     newSet("array"),
					Contains: &jSchema{Type: newSet("string")},
				},
			},
		},
		"array - nested mixed - merge": {
			in:    []any{[]any{"a", 42}, []any{nil}, []any{24.31}},
			merge: true,
			expected: &jSchema{
				Type: newSet("array"),
				Contains: &jSchema{
					Type:     newSet("array"),
					Contains: &jSchema{Type: newSet("string", "integer", "null", "number")},
				},
			},
		},
		"array - nested objects - merge": {
			in:    []any{[]any{map[string]any{"name": "hans"}}, []any{nil}, []any{map[string]any{"age": 42}}},
			merge: true,
			expected: &jSchema{
				Type: newSet("array"),
				Contains: &jSchema{
					Type: newSet("array"),
					Contains: &jSchema{
						Type: newSet("object", "null"),
						Properties: map[string]*jSchema{
							"name": {Type: newSet("string")},
							"age":  {Type: newSet("integer")},
						},
					},
				},
			},
		},
		"array - nested nested objects - merge": {
			in: []any{
				[]any{
					map[string]any{
						"name": "hans",
						"address": map[string]any{
							"street": "main st",
						}},
				},
				[]any{nil},
				[]any{
					map[string]any{
						"age": 42,
						"address": map[string]any{
							"town": "london",
							"tags": []any{"a", "b"},
						}},
				},
			},
			merge: true,
			expected: &jSchema{
				Type: newSet("array"),
				Contains: &jSchema{
					Type: newSet("array"),
					Contains: &jSchema{
						Type: newSet("object", "null"),
						Properties: map[string]*jSchema{
							"name": {Type: newSet("string")},
							"age":  {Type: newSet("integer")},
							"address": {
								Type: newSet("object"),
								Properties: map[string]*jSchema{
									"street": {Type: newSet("string")},
									"town":   {Type: newSet("string")},
									"tags": {
										Type:     newSet("array"),
										Contains: &jSchema{Type: newSet("string")},
									},
								},
							},
						},
					},
				},
			},
		},
		"array of objects": {
			in: []any{
				map[string]any{"name": "hans"},
				map[string]any{"name": "peter"},
			},
			expected: &jSchema{
				Type: newSet("array"),
				PrefixItems: []*jSchema{
					{
						Type:       newSet("object"),
						Properties: map[string]*jSchema{"name": {Type: newSet("string")}},
					},
					{
						Type:       newSet("object"),
						Properties: map[string]*jSchema{"name": {Type: newSet("string")}},
					},
				},
			},
		},
		"array of any": {
			in: []any{
				"color",
				map[string]any{"name": "hans"},
				42,
				map[string]any{"name": "peter"},
				true,
			},
			expected: &jSchema{
				Type: newSet("array"),
				PrefixItems: []*jSchema{
					{Type: newSet("string")},
					{Type: newSet("object"), Properties: map[string]*jSchema{"name": {Type: newSet("string")}}},
					{Type: newSet("integer")},
					{Type: newSet("object"), Properties: map[string]*jSchema{"name": {Type: newSet("string")}}},
					{Type: newSet("boolean")},
				},
			},
		},
		"array of different objects": {
			in: []any{
				map[string]any{"name": "hans"},
				map[string]any{"age": 42},
			},
			expected: &jSchema{
				Type: newSet("array"),
				PrefixItems: []*jSchema{
					{
						Type: newSet("object"),
						Properties: map[string]*jSchema{
							"name": {Type: newSet("string")},
						},
					},
					{
						Type: newSet("object"),
						Properties: map[string]*jSchema{
							"age": {Type: newSet("integer")},
						},
					},
				},
			},
		},
		"array of different objects - merge": {
			in: []any{
				map[string]any{"name": "hans"},
				map[string]any{"age": 42},
			},
			merge: true,
			expected: &jSchema{
				Type: newSet("array"),
				Contains: &jSchema{
					Type: newSet("object"),
					Properties: map[string]*jSchema{
						"name": {Type: newSet("string")},
						"age":  {Type: newSet("integer")},
					},
				},
			},
		},
		"array of different objects - merge more complex": {
			in: []any{
				map[string]any{
					"first-name": "hans",
					"address":    map[string]any{"street": "main str"},
					"last":       false,
				},
				true,
				"hans",
				map[string]any{
					"age":       42,
					"last-name": "schmitt",
					"verified":  true,
					"address":   map[string]any{"city": "metropolis"},
					"last":      42.11,
				},
			},
			merge: true,
			expected: &jSchema{
				Type: newSet("array"),
				Contains: &jSchema{
					Type: newSet("object", "boolean", "string"),
					Properties: map[string]*jSchema{
						"first-name": {Type: newSet("string")},
						"address": {
							Type: newSet("object"),
							Properties: map[string]*jSchema{
								"street": {Type: newSet("string")},
								"city":   {Type: newSet("string")},
							},
						},
						"last-name": {Type: newSet("string")},
						"verified":  {Type: newSet("boolean")},
						"age":       {Type: newSet("integer")},
						"last":      {Type: newSet("boolean", "number")},
					},
				},
			},
		},
	}

	for tn, tc := range tcc {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			actual := schema(tc.in, tc.merge)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func Test_schema_simple(t *testing.T) {
	tcc := map[string]struct {
		in       any
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
			expected: "boolean",
		},
		"false": {
			in:       false,
			expected: "boolean",
		},
		"nil": {
			in:       nil,
			expected: "null",
		},
		"empty array": {
			in:       []any{},
			expected: []any{},
		},
		"array of numbers": {
			in:       []any{1, 1.2, -3, 4},
			expected: []any{"integer", "number", "integer", "integer"},
		},
		"array of strings": {
			in:       []any{"hans", "", "peter"},
			expected: []any{"string", "string", "string"},
		},
		"array of any": {
			in:       []any{"hans", nil, -23.22, nil, 2, true, []any{"peter"}, map[string]any{"key": 42.222}},
			expected: []any{"string", "null", "number", "null", "integer", "boolean", []any{"string"}, map[string]any{"key": "number"}},
		},
		"array nested - merge": {
			in: []any{
				[]any{nil},
				[]any{nil},
			},
			merge: true,
			expected: []any{
				[]any{"null"},
			},
		},
		"array nested": {
			in: []any{
				[]any{nil},
				[]any{nil},
			},
			expected: []any{
				[]any{"null"},
				[]any{"null"},
			},
		},
		"array nested mixed": {
			in: []any{
				23,
				[]any{nil},
			},
			expected: []any{
				"integer",
				[]any{"null"},
			},
		},
		"array nested mixed empty - merge": {
			in: []any{
				23,
				[]any{},
			},
			merge: true,
			expected: []any{
				"integer",
				[]any{},
			},
		},
		"array nested mixed - merge": {
			in: []any{
				23,
				[]any{nil},
			},
			merge: true,
			expected: []any{
				"integer",
				[]any{"null"},
			},
		},
		"array nested mixed nested - merge": {
			in: []any{
				23,
				[]any{
					nil,
					1.1,
					[]any{true},
				},
			},
			merge: true,
			expected: []any{
				"integer",
				[]any{
					"null",
					"number",
					[]any{
						"boolean",
					},
				},
			},
		},
		"array nested mixed nested": {
			in: []any{
				23,
				[]any{
					nil,
					1.1,
					[]any{true},
				},
			},
			expected: []any{
				"integer",
				[]any{
					"null",
					"number",
					[]any{
						"boolean",
					},
				},
			},
		},
		"array of any - merge": {
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
			merge: true,
			expected: []any{
				"integer",
				map[string]any{
					"key1": "any",
					"key2": "boolean",
					"key3": map[string]any{"a": "boolean", "b": "integer"},
				},
				[]any{"null"},
				"boolean",
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
				"age":      "integer",
				"colors":   []any{"string", "string"},
				"verified": "boolean",
				"email":    "null",
			},
		},
	}

	for tn, tc := range tcc {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			s := schema(tc.in, tc.merge)
			assert.Equal(t, tc.expected, s.toSimple())
		})
	}
}
