package serport

import "testing"

var mergeData = [][]string{
	[]string{"pe D0 po", "b115200 l8 pe s1 D0 po"},
	[]string{"pe W0 b38400", "b115200 l8 pe s1 W0 b38400"},
	[]string{"b38400 D0 pe", "b38400 l8 pn s1 D0 pe"},
}

func TestMergeWithDefault(t *testing.T) {
	for i, d := range mergeData {
		if r := MergeCtlCmds(d[0]); r != d[1] {
			t.Fatalf("mergeData[%d]: expected: %q, got: %q", i, d[1], r)
		}
	}
}
