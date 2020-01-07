package utils

import (
	"testing"
)

type PairsToCompare struct {
	Wants string
	Test  string
}

var testsData = []PairsToCompare{
	PairsToCompare{Wants: "Hey", Test: "Hey"}, //no changes
	PairsToCompare{Wants: "DISH Network", Test: "DISH Network\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000"},
	PairsToCompare{Wants: "GoChain 4", Test: "GoChain 4\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000\u0000"},
}

func TestCleanUpText(t *testing.T) {
	for _, element := range testsData {
		if element.Wants != CleanUpText(element.Test) {
			t.Error("Wants:", element.Wants, "Got:", CleanUpText(element.Wants))
		}
	}
}
