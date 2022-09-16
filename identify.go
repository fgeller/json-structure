package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type jSchema struct {
	Schema      string              `json:"$schema,omitempty"`
	Type        *orderedSet[string] `json:"type,omitempty"`
	PrefixItems []*jSchema          `json:"prefixItems,omitempty"`
	Contains    *jSchema            `json:"contains,omitempty"`
	Properties  map[string]*jSchema `json:"properties,omitempty"`
}

func (s *jSchema) MarshalJSON() ([]byte, error) {
	if s.Type.len() == 1 {
		v := struct {
			Schema      string              `json:"$schema,omitempty"`
			Type        string              `json:"type,omitempty"`
			PrefixItems []*jSchema          `json:"prefixItems,omitempty"`
			Contains    *jSchema            `json:"contains,omitempty"`
			Properties  map[string]*jSchema `json:"properties,omitempty"`
		}{
			Schema:      s.Schema,
			Type:        *s.Type.takeFirst(),
			PrefixItems: s.PrefixItems,
			Contains:    s.Contains,
			Properties:  s.Properties,
		}
		return json.Marshal(v)
	}

	ts := []string{}
	for _, t := range s.Type.data {
		ts = append(ts, t)
	}

	v := struct {
		Schema      string              `json:"$schema,omitempty"`
		Type        []string            `json:"type,omitempty"`
		PrefixItems []*jSchema          `json:"prefixItems,omitempty"`
		Contains    *jSchema            `json:"contains,omitempty"`
		Properties  map[string]*jSchema `json:"properties,omitempty"`
	}{
		Schema:      s.Schema,
		Type:        ts,
		PrefixItems: s.PrefixItems,
		Contains:    s.Contains,
		Properties:  s.Properties,
	}
	return json.Marshal(v)
}

func (s *jSchema) isType(t string) bool {
	return s.Type.len() == 1 && *s.Type.takeFirst() == t
}

func (s *jSchema) mergeProperties(a *jSchema) {
	for ak, av := range a.Properties {
		if s.Properties == nil {
			s.Properties = map[string]*jSchema{}
		}

		sv, ok := s.Properties[ak]
		// key not present
		if !ok {
			s.Properties[ak] = av
			continue
		}

		// nested objects
		if av.isType("object") && sv.isType("object") {
			sv.mergeProperties(av)
			continue
		}

		// mixed/any type
		if !reflect.DeepEqual(av, sv) {
			s.Properties[ak].Type.merge(av.Type)
			continue
		}
	}
}

func simpleContains(s *jSchema) []any {
	if s == nil {
		return []any{}
	}

	as := []any{}
	for _, t := range s.Type.data {
		switch t {
		case "string", "number", "integer", "boolean", "null":
			as = append(as, t)
		case "array":
			na := simpleContains(s.Contains)
			as = append(as, na)
		case "object":
			obj := map[string]any{}
			for n, a := range s.Properties {
				if !a.isType("array") && a.Type.len() > 1 {
					obj[n] = "any"
					continue
				}
				obj[n] = a.toSimple()
			}
			as = append(as, obj)
		default:
			panic(fmt.Sprintf("unsupported type %#v", t))
		}
	}
	return as
}

func (s *jSchema) toSimple() any {
	switch {
	case s.isType("string"):
		return "string"
	case s.isType("number"):
		return "number"
	case s.isType("integer"):
		return "integer"
	case s.isType("boolean"):
		return "boolean"
	case s.isType("null"):
		return "null"
	case s.isType("array"):
		switch {
		case s.Contains != nil:
			return simpleContains(s.Contains)
		case s.PrefixItems != nil:
			as := []any{}
			for _, s := range s.PrefixItems {
				as = append(as, s.toSimple())
			}
			return as
		default:
			return []any{}
		}
	case s.isType("object"):
		obj := map[string]any{}
		for n, a := range s.Properties {
			as := a.toSimple()
			obj[n] = as
			_, isArr := as.([]any)
			if !a.isType("array") && isArr {
				obj[n] = "any"
			}
		}
		return obj
	default:
		panic(fmt.Sprintf("unsupported type %s", s.Type))
	}
}

func schema(d any, merge bool) *jSchema {
	switch d.(type) {
	case string:
		return &jSchema{Type: newSet("string")}
	case float64:
		return &jSchema{Type: newSet("number")}
	case int:
		return &jSchema{Type: newSet("integer")}
	case bool:
		return &jSchema{Type: newSet("boolean")}
	case nil:
		return &jSchema{Type: newSet("null")}
	case map[string]any:
		m := d.(map[string]any)
		js := &jSchema{Type: newSet("object"), Properties: map[string]*jSchema{}}
		for k, v := range m {
			js.Properties[k] = schema(v, merge)
		}
		return js
	case []any:
		s := d.([]any)
		if len(s) == 0 {
			return &jSchema{Type: newSet("array")}
		}

		jss := make([]*jSchema, len(s))
		for i, v := range s {
			jss[i] = schema(v, merge)
		}

		if merge {
			return &jSchema{
				Type:     newSet("array"),
				Contains: flatten(jss),
			}

		}
		return &jSchema{
			Type:        newSet("array"),
			PrefixItems: jss,
		}

	default:
		fmt.Fprintf(os.Stderr, "failed to identify type %T", d)
		return nil
	}
}

func flatten(jss []*jSchema) *jSchema {
	switch len(jss) {
	case 0:
		return nil
	case 1:
		return jss[0]
	}

	ms := jss[0]
	for _, v := range jss {
		ms.Type.merge(v.Type)
		ms.mergeProperties(v)
		ms.mergeContains(v)
	}
	return ms
}

func (s *jSchema) mergeContains(a *jSchema) {
	if s.Contains == nil {
		s.Contains = a.Contains
		return
	}
	if s.Contains != nil && a.Contains != nil {
		s.Contains.Type.merge(a.Contains.Type)
		s.Contains.mergeProperties(a.Contains)
	}
}
