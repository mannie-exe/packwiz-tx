package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func setupTestPack(t *testing.T) (Index, Pack, string) {
	t.Helper()
	dir := t.TempDir()

	packFile := filepath.Join(dir, "pack.toml")
	indexFile := filepath.Join(dir, "index.toml")

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

	err = os.WriteFile(indexFile, []byte(`hash-format = "sha256"
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	viper.Set("pack-file", packFile)
	t.Cleanup(func() { viper.Set("pack-file", nil) })

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

	index.HashFormat = "sha512"

	err := CommitChanges(&index, &pack)
	if err != nil {
		t.Fatalf("CommitChanges should succeed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "index.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "sha512") {
		t.Error("index.toml should contain 'sha512' after CommitChanges")
	}
}

func TestCommitChangesSkipsWhenNoRefreshTrue(t *testing.T) {
	index, pack, dir := setupTestPack(t)

	viper.Set("no-refresh", true)
	defer viper.Set("no-refresh", nil)

	index.HashFormat = "sha512"

	err := CommitChanges(&index, &pack)
	if err != nil {
		t.Fatalf("CommitChanges should succeed (no-op): %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "index.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(content), "sha512") {
		t.Error("index.toml should NOT contain 'sha512' when --no-refresh is set")
	}
}

func TestCommitChangesDefaultsToWrite(t *testing.T) {
	index, pack, dir := setupTestPack(t)

	viper.Set("no-refresh", nil)

	index.HashFormat = "sha512"

	err := CommitChanges(&index, &pack)
	if err != nil {
		t.Fatalf("CommitChanges should succeed with default config: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "index.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "sha512") {
		t.Error("index.toml should contain 'sha512' when no-refresh defaults to false")
	}
}
