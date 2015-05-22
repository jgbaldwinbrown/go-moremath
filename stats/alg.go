// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stats

// Miscellaneous helper algorithms

import (
	"fmt"
	"math"
)

// sign returns the sign of x: -1 if x < 0, 0 if x == 0, 1 if x > 0.
// If x is NaN, it returns NaN.
func sign(x float64) float64 {
	if x == 0 {
		return 0
	} else if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	}
	return nan
}

func maxint(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minint(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sumint(xs []int) int {
	sum := 0
	for _, x := range xs {
		sum += x
	}
	return sum
}

// lchoose returns math.Log(choose(n, k)).
func lchoose(n, k int) float64 {
	a, _ := math.Lgamma(float64(n + 1))
	b, _ := math.Lgamma(float64(k + 1))
	c, _ := math.Lgamma(float64(n - k + 1))
	return a - b - c
}

const smallFactLimit = 20 // 20! => 62 bits
var smallFact [smallFactLimit + 1]int64

func init() {
	smallFact[0] = 1
	fact := int64(1)
	for n := int64(1); n <= smallFactLimit; n++ {
		fact *= n
		smallFact[n] = fact
	}
}

// choose returns the binomial coefficient of n and k.
func choose(n, k int) int {
	if k == 0 || k == n {
		return 1
	}
	if k < 0 || n < k {
		return 0
	}
	if n <= smallFactLimit { // Implies k <= smallFactLimit
		// It's faster to do several integer multiplications
		// than it is to do an extra integer division.
		// Remarkably, this is also faster than pre-computing
		// Pascal's triangle (presumably because this is very
		// cache efficient).
		numer := int64(1)
		for n1 := int64(n - (k - 1)); n1 <= int64(n); n1++ {
			numer *= n1
		}
		denom := smallFact[k]
		return int(numer / denom)
	}

	return int(math.Exp(lchoose(n, k)) + 0.5)
}

// atEach returns f(x) for each x in xs.
func atEach(f func(float64) float64, xs []float64) []float64 {
	// TODO(austin) Parallelize
	res := make([]float64, len(xs))
	for i, x := range xs {
		res[i] = f(x)
	}
	return res
}

// bisect returns an x in [low, high] such that |f(x)| <= tolerance
// using the bisection method.
//
// f(low) and f(high) must have opposite signs.
//
// If f does not have a root in this interval (e.g., it is
// discontiguous), this returns the X of the apparent discontinuity
// and false.
func bisect(f func(float64) float64, low, high, tolerance float64) (float64, bool) {
	flow, fhigh := f(low), f(high)
	if -tolerance <= flow && flow <= tolerance {
		return low, true
	}
	if -tolerance <= fhigh && fhigh <= tolerance {
		return high, true
	}
	if sign(flow) == sign(fhigh) {
		panic(fmt.Sprintf("root of f is not bracketed by [low, high]; f(%g)=%g f(%g)=%g", low, flow, high, fhigh))
	}
	for {
		mid := (high + low) / 2
		fmid := f(mid)
		if -tolerance <= fmid && fmid <= tolerance {
			return mid, true
		}
		if mid == high || mid == low {
			return mid, false
		}
		if sign(fmid) == sign(flow) {
			low = mid
			flow = fmid
		} else {
			high = mid
			fhigh = fmid
		}
	}
}

// bisectBool implements the bisection method on a boolean function.
// It returns x1, x2 ∈ [low, high], x1 < x2 such that f(x1) != f(x2)
// and x2 - x1 <= xtol.
//
// If f(low) == f(high), it panics.
func bisectBool(f func(float64) bool, low, high, xtol float64) (x1, x2 float64) {
	flow, fhigh := f(low), f(high)
	if flow == fhigh {
		panic(fmt.Sprintf("root of f is not bracketed by [low, high]; f(%g)=%v f(%g)=%v", low, flow, high, fhigh))
	}
	for {
		if high-low <= xtol {
			return low, high
		}
		mid := (high + low) / 2
		if mid == high || mid == low {
			return low, high
		}
		fmid := f(mid)
		if fmid == flow {
			low = mid
			flow = fmid
		} else {
			high = mid
			fhigh = fmid
		}
	}
}

// series returns the sum of the series f(0), f(1), ...
//
// This implementation is fast, but subject to round-off error.
func series(f func(float64) float64) float64 {
	y, yp := 0.0, 1.0
	for n := 0.0; y != yp; n++ {
		yp = y
		y += f(n)
	}
	return y
}
