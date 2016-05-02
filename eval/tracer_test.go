// Copyright 2016 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package eval

import (
	"fmt"
	"testing"
)

type mockTracer struct {
	buf []string
}

func (t *mockTracer) Enabled() bool { return true }

func (t *mockTracer) Trace(ctx *TopDownContext, f string, a ...interface{}) {
	t.buf = append(t.buf, fmt.Sprintf(f, a))
}

func TestTracer(t *testing.T) {

	rules := parseRules([]string{
		"p[x] :- q[x] = y",
		"q[i] = j :- a[i] = j",
	})

	store := newStorage(loadSmallTestData(), rules)

	tracer := &mockTracer{[]string{}}

	params := &TopDownQueryParams{
		Store:  store,
		Tracer: tracer,
		Path:   []string{"p"}}

	result, err := TopDownQuery(params)
	if err != nil {
		t.Errorf("Unexpected error in tracing test: %v", err)
		return
	}

	expected := []interface{}{float64(0), float64(1), float64(2), float64(3)}

	if Compare(result, expected) != 0 {
		t.Errorf("Unexpected result in tracing test: %v", result)
		return
	}

	// ((try success finish * 2) * 4) + 2
	// 2 rules
	// 4 elements in a
	if len(tracer.buf) != 26 {
		t.Errorf("Unexpected output from tracer:\n%v", tracer.buf)
	}
}