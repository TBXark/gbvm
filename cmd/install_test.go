package cmd

import (
	"errors"
	"os"
	"testing"

	"github.com/TBXark/gbvm/internal/bin"
)

func TestParseBackupFile(t *testing.T) {
	data := []byte(`[{"name":"tool","version":"v1.2.3","mod":"example.com/tool","path":"example.com/tool/cmd"}]`)
	versions, err := parseBackupFile(data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(versions) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(versions))
	}
	entry := versions[0]
	if entry.Name != "tool" || entry.Version != "v1.2.3" || entry.Mod != "example.com/tool" || entry.Path != "example.com/tool/cmd" {
		t.Fatalf("unexpected entry: %#v", entry)
	}
}

func TestParseBackupFileInvalidJSON(t *testing.T) {
	_, err := parseBackupFile([]byte(`{not-json}`))
	if err == nil {
		t.Fatalf("expected error for invalid JSON")
	}
}

func TestDecideInstall(t *testing.T) {
	desired := &bin.Info{Name: "tool", Version: "v1.0.0"}

	decision := decideInstall(desired, &bin.Info{Version: "v1.0.0"}, nil)
	if decision.shouldInstall {
		t.Fatalf("expected skip when versions match")
	}
	if decision.message != "skip tool" {
		t.Fatalf("unexpected message: %s", decision.message)
	}

	decision = decideInstall(desired, &bin.Info{Version: "v0.9.0"}, nil)
	if !decision.shouldInstall {
		t.Fatalf("expected install when versions differ")
	}
	if decision.message != "installing tool@v1.0.0" {
		t.Fatalf("unexpected message: %s", decision.message)
	}

	decision = decideInstall(desired, nil, os.ErrNotExist)
	if !decision.shouldInstall {
		t.Fatalf("expected install when binary missing")
	}
	if decision.message != "installing tool@v1.0.0" {
		t.Fatalf("unexpected message: %s", decision.message)
	}

	loadErr := errors.New("boom")
	decision = decideInstall(desired, nil, loadErr)
	if decision.shouldInstall {
		t.Fatalf("expected skip on load error")
	}
	if decision.message != "failed to load tool: boom" {
		t.Fatalf("unexpected message: %s", decision.message)
	}
}
