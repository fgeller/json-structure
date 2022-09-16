package main

import (
	"fmt"
	"os"
	"reflect"
)

type jSchema struct {
	Schema      string              `json:"$schema,omitempty"`
	ID          string              `json:"$id,omitempty"`
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	Type        string              `json:"type"`
	Items       *jSchema            `json:"items"`
	Properties  map[string]*jSchema `json:"properties,omitempty"`
	Required    []string            `json:"required,omitempty"`
}

func schema(d any) *jSchema {
	switch d.(type) {
	case string:
		return &jSchema{Type: "string"}
	case float64:
		return &jSchema{Type: "number"}
	case int:
		return &jSchema{Type: "integer"}
	case bool:
		return &jSchema{Type: "boolean"}
	case nil:
		return &jSchema{Type: "null"}
	case map[string]any:
		m := d.(map[string]any)
		js := &jSchema{Type: "object", Properties: map[string]*jSchema{}}
		for k, v := range m {
			js.Properties[k] = schema(v)
		}
		return nil
	case []any:
		s := d.([]any)
		js := &jSchema{Type: "array"}
		jss := make([]*jSchema, len(s))
		if len(s) == 0 {
			return js
		}
		for i, v := range s {
			jss[i] = schema(v)
		}
		jss = dedupe(jss)
	}
	// TODO merge?
		return s

	default:
		fmt.Fprintf(os.Stderr, "failed to identify type %T", d)
		return nil
	}
}

func dedupe(a []any) []any {
	mvv := []any{a[0]}
outer:
	for _, v := range a {
		for _, mv := range mvv {
			if reflect.DeepEqual(v, mv) {
				continue outer
			}
		}
		mvv = append(mvv, v)
	}

	return mvv
}

func mergeArray(a []any) []any {
	mvv := []any{}
	mo := map[string]any{}
	fi := -1
	for i, v := range a {
		switch v.(type) {
		case map[string]any:
			if fi == -1 {
				fi = i
				mo = v.(map[string]any)
				mvv = append(mvv, v)
				continue
			}
			mo = merge(mo, v.(map[string]any))

		default:
			mvv = append(mvv, v)
		}
	}
	if fi != -1 {
		mvv[fi] = mo
	}
	return mvv
}

func merge(a, b map[string]any) map[string]any {
	mo := a
	for bk, bv := range b {
		av, ok := a[bk]
		if !ok {
			mo[bk] = bv
			continue
		}

		if isObject(av) && isObject(bv) {
			mo[bk] = merge(av.(map[string]any), bv.(map[string]any))
			continue
		}

		if !reflect.DeepEqual(bv, av) {
			mo[bk] = "any"
			continue
		}
	}

	return mo
}

func isObject(a any) bool {
	_, ok := a.(map[string]any)
	return ok
}
