package extensionsdetection_test

import (
	"os"
	"path/filepath"
	"testing"

	chastlog "chast.io/core/internal/logger"
	uut "chast.io/core/internal/run_model/internal/extensions_detection"
	"github.com/joomcode/errorx"
)

func TestDetectExtensions_Valid(t *testing.T) { //nolint:gocognit // This is contains nested tests
	t.Parallel()

	baseDir, _ := os.MkdirTemp("", "TestDetectExtensions")
	paths := []string{
		"/project/PACKAGE1/file1.java",
		"/project/PACKAGE1/file2.java",
		"/project/PACKAGE2/file1.java",

		"/project/file1.cs",
		"/project/file2.cs",
		"/project/sub1/sub1/sub1/sub1/sub1/file3.cs",

		"/project/sub1/sub1/SUB1/sub1/sub1/file1.go",
		"/project/sub1/sub1/SUB2/sub1/sub1/file1.go",
	}

	if err := preparePaths(paths, baseDir); err != nil {
		t.Fatalf("Error preparing paths: %v", err)
	}

	t.Cleanup(func() { cleanupPaths(baseDir) })

	extensions, err := uut.DetectExtensions(baseDir)

	if err != nil {
		t.Fatalf("Expected no error, but was '%v'", err)
	}

	t.Run("Extensions", func(t *testing.T) {
		t.Parallel()

		if len(extensions) != 3 {
			t.Errorf("Expected 3 extensions, but was %d", len(extensions))
		}

		t.Run("Java", func(t *testing.T) {
			t.Parallel()

			if extensions["java"] == nil {
				t.Fatalf("Expected Java extension to be set, but was nil")
			}

			if extensions["java"].Count != 3 {
				t.Errorf("Expected Java extension count to be 4, but was %d", extensions["java"].Count)
			}

			expectedParent := filepath.Join(baseDir, "/project")
			if extensions["java"].CommonParentPath != expectedParent {
				t.Errorf("Expected Java extension common parent path to be '%s', but was '%s'", expectedParent, extensions["java"].CommonParentPath)
			}
		})

		t.Run("CSharp", func(t *testing.T) {
			t.Parallel()

			if extensions["cs"] == nil {
				t.Fatalf("Expected cs extension to be set, but was nil")
			}

			if extensions["cs"].Count != 3 {
				t.Errorf("Expected cs extension count to be 2, but was %d", extensions["cs"].Count)
			}

			expectedParent := filepath.Join(baseDir, "/project")
			if extensions["cs"].CommonParentPath != expectedParent {
				t.Errorf("Expected cs extension common parent path to be '%s', but was '%s'", expectedParent, extensions["cs"].CommonParentPath)
			}
		})

		t.Run("Go", func(t *testing.T) {
			t.Parallel()

			if extensions["go"] == nil {
				t.Fatalf("Expected Go extension to be set, but was nil")
			}

			if extensions["go"].Count != 2 {
				t.Errorf("Expected Go extension count to be 2, but was %d", extensions["go"].Count)
			}

			expectedParent := filepath.Join(baseDir, "/project/sub1/sub1")
			if extensions["go"].CommonParentPath != expectedParent {
				t.Errorf("Expected Go extension common parent path to be '%s', but was '%s'", expectedParent, extensions["go"].CommonParentPath)
			}
		})
	})
}

func TestDetectExtensions_Invalid(t *testing.T) {
	t.Parallel()

	_, err := uut.DetectExtensions("/invalid/path")

	if err == nil {
		t.Fatalf("Expected error, but was nil")
	}
}

func preparePaths(paths []string, baseDir string) error {
	for _, path := range paths {
		fullPath := filepath.Join(baseDir, path)

		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			return err
		}

		if _, err := os.Create(fullPath); err != nil {
			return err
		}
	}

	return nil
}

func cleanupPaths(baseDir string) {
	if err := os.RemoveAll(baseDir); err != nil {
		chastlog.Log.Fatalln(errorx.ExternalError.Wrap(err, "Error cleaning up paths"))
	}
}
