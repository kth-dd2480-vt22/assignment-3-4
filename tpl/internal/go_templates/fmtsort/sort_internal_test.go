package fmtsort

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	exitCode := t.Run()

	fmt.Printf("Branch coverage of compare():\n")

	covered := 0
	for k, v := range CompareBC {
		if !v.Reached {
			fmt.Printf("Condition %d never reached.\n", k)
		}

		if !v.True {
			fmt.Printf("Condition %d never reached TRUE.\n", k)
		} else {
			covered++
		}

		if !v.False {
			fmt.Printf("Condition %d never reached FALSE.\n", k)
		} else {
			covered++
		}
	}

	nBranch := len(CompareBC) * 2
	fmt.Printf("Coverage: %.2f%% (%d/%d)\n", float64(covered)/float64(nBranch)*100, covered, nBranch)

	os.Exit(exitCode)
}
