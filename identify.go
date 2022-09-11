package main

import (
	"fmt"
	"os"
	"reflect"
)

func identify(d any, dd, mo bool) any {
	switch d.(type) {
	case string:
		return "string"
	case float64, int:
		return "number"
	case bool:
		return "bool"
	case map[string]any:
		m := d.(map[string]any)
		for k, v := range m {
			m[k] = identify(v, dd, mo)
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
	case nil:
		return nil

	default:
		fmt.Fprintf(os.Stderr, "failed to identify type %T", d)
		return d
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
