package main

import (
	"fmt"
	"strings"
)

type orderedSet[T comparable] struct {
	data []T
}

func (s *orderedSet[T]) String() string {
	vs := []string{}
	for _, v := range s.data {
		vs = append(vs, fmt.Sprintf("%v", v))
	}
	return "orderedSet{" + strings.Join(vs, ",") + "}"
}

func newSet[T comparable](ts ...T) *orderedSet[T] {
	res := &orderedSet[T]{data: make([]T, len(ts))}
	for i, t := range ts {
		res.data[i] = t
	}
	return res
}

func (s *orderedSet[T]) add(ts ...T) {
	toAdd := make([]T, 0, len(ts))
outer:
	for _, t := range ts {
		for _, d := range s.data {
			if t == d {
				continue outer
			}
		}
		toAdd = append(toAdd, t)
	}
	s.data = append(s.data, toAdd...)
}

func (s *orderedSet[T]) merge(a *orderedSet[T]) {
	s.add(a.data...)
}

func (s *orderedSet[T]) takeFirst() *T {
	for _, v := range s.data {
		return &v
	}
	return nil
}

func (s *orderedSet[T]) len() int {
	return len(s.data)
}

func (s *orderedSet[T]) contains(t T) bool {
	for _, v := range s.data {
		if v == t {
			return true
		}
	}
	return false
}
