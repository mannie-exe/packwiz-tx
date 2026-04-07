package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func setupTestPack(t *testing.T) (Index, Pack, string) {
	t.Helper()
	dir := t.TempDir()

	packFile := filepath.Join(dir, "pack.toml")
	indexFile := filepath.Join(dir, "index.toml")

	// Write minimal pack.toml
	err := os.WriteFile(packFile, []byte(`name = "test-pack"
pack-format = "packwiz:1.1.0"

[index]
file = "index.toml"
hash-format = "sha256"
hash = ""

[versions]
minecraft = "1.21.1"
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Write minimal index.toml
	err = os.WriteFile(indexFile, []byte(`hash-format = "sha256"
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	viper.Set("pack-file", packFile)
	pack, err := LoadPack()
	if err != nil {
		t.Fatal(err)
	}
	index, err := pack.LoadIndex()
	if err != nil {
		t.Fatal(err)
	}

	return index, pack, dir
}

func TestCommitChangesWritesWhenNoRefreshFalse(t *testing.T) {
	index, pack, dir := setupTestPack(t)

	viper.Set("no-refresh", false)
	defer viper.Set("no-refresh", nil)

	// Modify index in memory
	index.HashFormat = "sha512"

	err := CommitChanges(&index, &pack)
	if err != nil {
		t.Fatalf("CommitChanges should succeed: %v", err)
	}

	// Verify index.toml was written with the new hash format
	content, err := os.ReadFile(filepath.Join(dir, "index.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !containsString(string(content), "sha512") {
		t.Error("index.toml should contain 'sha512' after CommitChanges")
	}
}

func TestCommitChangesSkipsWhenNoRefreshTrue(t *testing.T) {
	index, pack, dir := setupTestPack(t)

	viper.Set("no-refresh", true)
	defer viper.Set("no-refresh", nil)

	// Modify index in memory
	index.HashFormat = "sha512"

	err := CommitChanges(&index, &pack)
	if err != nil {
		t.Fatalf("CommitChanges should succeed (no-op): %v", err)
	}

	// Verify index.toml was NOT written (still has original sha256)
	content, err := os.ReadFile(filepath.Join(dir, "index.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if containsString(string(content), "sha512") {
		t.Error("index.toml should NOT contain 'sha512' when --no-refresh is set")
	}
}

func TestCommitChangesDefaultsToWrite(t *testing.T) {
	index, pack, _ := setupTestPack(t)

	// Don't set no-refresh at all; should default to false (write)
	viper.Set("no-refresh", nil)

	err := CommitChanges(&index, &pack)
	if err != nil {
		t.Fatalf("CommitChanges should succeed with default config: %v", err)
	}
}

func containsString(haystack, needle string) bool {
	return len(haystack) > 0 && len(needle) > 0 && // avoid trivial matches
		indexOf(haystack, needle) >= 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
