package internal_test

import (
	"testing"
)

func TestCheck(t *testing.T) {
	testRoot, testCache, wd := prepare(t)

	t.Log(testRoot, testCache, wd)
}

//TODO: Add tests for all failing cases and check error and wether or not they generate anything
