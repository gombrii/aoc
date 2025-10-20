package internal_test

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/fs"
)

func assertDirEqual(t *testing.T, wantDir, gotDir string) {
	t.Helper()
	expected := fs.ManifestFromDir(t, wantDir)
	assert.Assert(t, fs.Equal(gotDir, expected))
}

func prepare(t *testing.T) (tmpRoot, tmpCache, wd string) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fail()
	}
	tmpRoot = t.TempDir()
	tmpCache = t.TempDir()
	t.Chdir(tmpRoot)
	t.Setenv("AOC_CACHE", tmpCache)

	return tmpRoot, tmpCache, wd
}
