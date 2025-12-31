package main

import (
	"testing"
)

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
}

func TestBuildTime(t *testing.T) {
	if BuildTime == "" {
		t.Error("BuildTime should not be empty")
	}
}
