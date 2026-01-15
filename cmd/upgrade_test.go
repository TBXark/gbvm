package cmd

import (
	"testing"

	"github.com/TBXark/gbvm/internal/bin"
)

func TestDetermineUpgrade(t *testing.T) {
	info := &bin.Info{Name: "tool", Version: "v1.0.0"}

	shouldUpgrade, message := determineUpgrade(info, "v1.2.0")
	if !shouldUpgrade {
		t.Fatalf("expected upgrade to be required")
	}
	if message != "upgrading tool from v1.0.0 to v1.2.0" {
		t.Fatalf("unexpected message: %s", message)
	}

	shouldUpgrade, message = determineUpgrade(info, "v1.0.0")
	if shouldUpgrade {
		t.Fatalf("expected no upgrade when versions match")
	}
	if message != "tool is already the latest version" {
		t.Fatalf("unexpected message: %s", message)
	}
}
