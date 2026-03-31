package test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Run all tests first, cleanup runs afterward for shared resources.
	exitCode := m.Run()

	// Ensure shared DB/API resources are torn down once per package run.
	teardownShared()

	os.Exit(exitCode)
}
