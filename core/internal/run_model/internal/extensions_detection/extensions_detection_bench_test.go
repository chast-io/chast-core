package extensionsdetection_test

import (
	"chast.io/core/internal/run_model/internal/extensions_detection"
	"testing"
)

func BenchmarkDetectExtensions(b *testing.B) {
	workingDirectory := "/var" // this path can be adjusted to any path of interest
	b.Logf("Working directory: %s", workingDirectory)

	extensions, err := extensionsdetection.DetectExtensions(workingDirectory)

	b.StopTimer()

	if err != nil {
		b.Fatalf("Expected no error, but was '%v'", err)
	}

	for _, extension := range extensions {
		b.Logf("Extension: %10s, Count: %5d, Common Parent: %s", extension.Name, extension.Count, extension.CommonParentPath)
	}

}
