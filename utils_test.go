package jotframe

import (
	"testing"
)

func Test_Utils_VisualLength(t *testing.T) {
	ansiString := "\x1b[2mHello, World!\x1b[0m"
	normalLen := len(ansiString)
	if normalLen != 21 {
		t.Error("TestVisualLength: Test harness not working! Expected 15 got", normalLen)
	}

	visualLen := VisualLength(ansiString)
	if visualLen != 13 {
		t.Error("Expected 13 got", visualLen)
	}
}

func Test_Utils_TrimToVisualLength(t *testing.T) {
	normalString := "Hello, World!"
	ansiString := "\x1b[2mHel\x1b[3mlo, Wor\x1b[0mld!\x1b[0m"

	for idx := 0; idx < len(normalString); idx++ {
		trimString := TrimToVisualLength(ansiString, idx)
		if VisualLength(trimString) != idx {
			t.Error("TestTrimToVisualLength: Expected", idx, "got", VisualLength(trimString), ". Trim:", trimString)
		}
	}
}
