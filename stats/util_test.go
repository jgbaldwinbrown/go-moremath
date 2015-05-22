// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stats

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"testing"
)

func aeq(expect, got float64) bool {
	if expect < 0 && got < 0 {
		expect, got = -expect, -got
	}
	return expect*0.99999999 <= got && got*0.99999999 <= expect
}

func testFunc(t *testing.T, name string, f func(float64) float64, vals map[float64]float64) {
	xs := make([]float64, 0, len(vals))
	for x := range vals {
		xs = append(xs, x)
	}
	sort.Float64s(xs)

	for _, x := range xs {
		want, got := vals[x], f(x)
		if math.IsNaN(want) && math.IsNaN(got) || aeq(want, got) {
			continue
		}
		var label string
		if strings.Contains(name, "%") {
			label = fmt.Sprintf(name, x)
		} else {
			label = fmt.Sprintf("%s(%v)", name, x)
		}
		t.Errorf("want %s=%v, got %v", label, want, got)
	}
}

func testDiscreteCDF(t *testing.T, name string, dist DiscreteDist) {
	// Build the expected CDF out of the PMF.
	l, h := dist.Bounds()
	s := dist.Step()
	want := map[float64]float64{l - 0.1: 0, h: 1}
	sum := 0.0
	for x := l; x < h; x += s {
		sum += dist.PMF(x)
		want[x] = sum
		want[x+s/2] = sum
	}

	testFunc(t, name, dist.CDF, want)
}
