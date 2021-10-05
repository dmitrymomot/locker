package locker

import (
	"testing"
	"time"
)

func Test_genValue(t *testing.T) {
	count := 10000
	res := map[string]struct{}{}
	for i := 0; i < count; i++ {
		v := defaultGenValue()
		if _, ok := res[string(v)]; ok {
			t.Fatalf("value already exists: %v", v)
		}
		res[string(v)] = struct{}{}
	}
	if len(res) != count {
		t.Fatalf("result values: go=%d, want=%d", len(res), count)
	}
}

func Test_defaultDelayFunc(t *testing.T) {
	d := defaultDelayFunc(1)
	w := time.Duration(100 * time.Millisecond)
	if d != w {
		t.Fatalf("defaultDelayFunc(1) = %v, want = %v", d, w)
	}

	d2 := defaultDelayFunc(0)
	w2 := time.Duration(100 * time.Millisecond)
	if d2 != w2 {
		t.Fatalf("defaultDelayFunc(0) = %v, want = %v", d2, w2)
	}

	d3 := defaultDelayFunc(10)
	w3 := time.Duration(100 * time.Millisecond)
	if d3 < w3 {
		t.Fatalf("defaultDelayFunc(0) = %v < %v", d3, w3)
	}
}
