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
	PrefixItems []*jSchema          `json:"prefixItems"`
	Properties  map[string]*jSchema `json:"properties,omitempty"`
	Required    []string            `json:"required,omitempty"`
}

func schema(d any, dd, mo bool) *jSchema {
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
			js.Properties[k] = schema(v, dd, mo)
		}
		return js
	case []any:
		s := d.([]any)
		if len(s) == 0 {
			return &jSchema{Type: "array"}
		}

		jss := make([]*jSchema, len(s))
		for i, v := range s {
			jss[i] = schema(v, dd, mo)
		}

		if dd {
			djss := []*jSchema{jss[0]}
		outer:
			for _, v := range jss {
				for _, dv := range djss {
					if reflect.DeepEqual(v, dv) {
						continue outer
					}
				}
				djss = append(djss, v)
			}
			jss = djss
		}

		if mo {
			jss = mergeArray(jss)
		}

		if len(jss) == 1 {
			return &jSchema{
				Type:  "array",
				Items: jss[0],
			}
		}
		return &jSchema{
			Type:        "array",
			PrefixItems: jss,
		}

	default:
		fmt.Fprintf(os.Stderr, "failed to identify type %T", d)
		return nil
	}
}

func identify(d any, dd, mo bool) any {
	switch d.(type) {
	case string:
		return "string"
	case float64, int:
		return "number"
	case bool:
		return "boolean"
	case nil:
		return "null"
	case map[string]any:
		m := d.(map[string]any)
		for k, v := range m {
			m[k] = identify(v, dd, mo) // TODO(fg) 🤨
		}
		return m
	case []any:
		s := d.([]any)
		if len(s) == 0 {
			return s
		}
		for i, v := range s {
			s[i] = identify(v, dd, mo)
		}
		if dd {
			s = dedupe(s)
		}
		if mo {
			s = mergeArray(s)
		}
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
		switch any(v).(type) {
		case map[string]any:
			if fi == -1 {
				fi = i
				mo = any(v).(map[string]any)
				mvv = append(mvv, v)
				continue
			}
			mo = merge(mo, any(v).(map[string]any))

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
