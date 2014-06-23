// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package newmath

import (
	"flag"
	"math"
	"testing"
)

var err = flag.Float64("err", 1e-6, "absolute error tolerance")

var sqrtTests = []struct {
	In  float64
	Out float64
}{
	{1, 1},
	{2, 1.4142135623730951},
	{3, 1.7320508075688771},
	{100, 10},
	{200, 14.142135623730951},
}

func TestSqrt(t *testing.T) {
	for _, tt := range sqrtTests {
		if x := Sqrt(tt.In); math.Abs(x-tt.Out) > *err {
			t.Errorf("Sqrt(%v) = %v, want %v", tt.In, x, tt.Out)
		}
	}
}
