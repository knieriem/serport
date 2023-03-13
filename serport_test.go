package serport

import "testing"

var mergeData = [][]string{
	[]string{"r0 D0 r1", "b115200 l8 pn r0 s1 D0 r1"},
	[]string{"r0 W0 b38400", "b115200 l8 pn r0 s1 W0 b38400"},
	[]string{"b38400 D0 r0", "b38400 l8 pn r1 s1 D0 r0"},
}

func TestMergeWithDefault(t *testing.T) {
	for i, d := range mergeData {
		if r := MergeCtlCmds(d[0]); r != d[1] {
			t.Fatalf("mergeData[%d]: expected: %q, got: %q", i, d[1], r)
		}
	}
}
