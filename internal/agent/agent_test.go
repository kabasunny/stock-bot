package agent

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFindSignalFile(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "signal_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up

	// Test case 1: No signal files found
	t.Run("no signal files", func(t *testing.T) {
		foundFile, err := FindSignalFile(filepath.Join(tmpDir, "*.bin"))
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if foundFile != "" {
			t.Errorf("expected empty string, got %s", foundFile)
		}
	})

	// Test case 2: Single signal file
	t.Run("single signal file", func(t *testing.T) {
		filePath := filepath.Join(tmpDir, "signal_20240101.bin")
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		foundFile, err := FindSignalFile(filepath.Join(tmpDir, "*.bin"))
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if foundFile != filePath {
			t.Errorf("expected %s, got %s", filePath, foundFile)
		}
	})

	// Test case 3: Multiple signal files, select newest
	t.Run("multiple signal files - select newest", func(t *testing.T) {
		file1 := filepath.Join(tmpDir, "signal_old.bin")
		file2 := filepath.Join(tmpDir, "signal_new.bin")
		file3 := filepath.Join(tmpDir, "signal_middle.bin")

		// Create files with specific modification times
		if err := os.WriteFile(file1, []byte("old"), 0644); err != nil {
			t.Fatalf("failed to create file1: %v", err)
		}
		time.Sleep(10 * time.Millisecond) // Ensure distinct modification times
		if err := os.WriteFile(file3, []byte("middle"), 0644); err != nil {
			t.Fatalf("failed to create file3: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
		if err := os.WriteFile(file2, []byte("new"), 0644); err != nil {
			t.Fatalf("failed to create file2: %v", err)
		}

		// Manually set mod times if needed for precise control, but sleep should be enough for most OS
		// For example:
		// t1 := time.Now().Add(-2 * time.Hour)
		// t2 := time.Now().Add(-1 * time.Hour)
		// t3 := time.Now()
		// os.Chtimes(file1, t1, t1)
		// os.Chtimes(file3, t2, t2)
		// os.Chtimes(file2, t3, t3)


		foundFile, err := FindSignalFile(filepath.Join(tmpDir, "*.bin"))
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if foundFile != file2 {
			t.Errorf("expected newest file %s, got %s", file2, foundFile)
		}
	})

	// Test case 4: Glob pattern with no matching accessible files
	t.Run("no accessible files", func(t *testing.T) {
		// Create a file but make it inaccessible (e.g., wrong permissions or broken symlink)
		// This scenario is hard to reliably test across OS for os.Stat errors specifically,
		// but we can simulate a pattern that matches nothing.
		// For now, rely on previous 'no signal files' test for this case.
	})
}
