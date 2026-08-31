package projectpath

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveProjectPathStaysWithinConfiguredRoots(t *testing.T) {
	allowedRoot := t.TempDir()
	projectDirectory := filepath.Join(allowedRoot, "project")
	if err := os.Mkdir(projectDirectory, 0o755); err != nil {
		t.Fatalf("create project directory: %v", err)
	}
	t.Setenv(codeStatsAllowedRootsEnv, allowedRoot)

	resolved, err := resolveProjectPath(projectDirectory)
	if err != nil {
		t.Fatalf("resolve allowed project path: %v", err)
	}
	if resolved != projectDirectory {
		t.Fatalf("resolved path = %q, want %q", resolved, projectDirectory)
	}

	outsideRoot := t.TempDir()
	if _, err := resolveProjectPath(outsideRoot); err == nil {
		t.Fatal("expected a path outside the configured root to be rejected")
	}
}

func TestResolveProjectPathRejectsSymlinkEscape(t *testing.T) {
	allowedRoot := t.TempDir()
	outsideRoot := t.TempDir()
	linkPath := filepath.Join(allowedRoot, "outside")
	if err := os.Symlink(outsideRoot, linkPath); err != nil {
		t.Fatalf("create escape symlink: %v", err)
	}
	t.Setenv(codeStatsAllowedRootsEnv, allowedRoot)

	if _, err := resolveProjectPath(linkPath); err == nil {
		t.Fatal("expected a symlink escaping the configured root to be rejected")
	}
}
