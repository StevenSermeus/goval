package file_test

import (
	"os"
	"path/filepath"
	"testing"
)

// Dummy FileExists function to test
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func TestFileExists(t *testing.T) {
	testDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	testFile, err := os.CreateTemp(testDir, "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFile.Name())

	testFilePath := testFile.Name()

	nonexistentFilePath := filepath.Join(testDir, "nonexistentfile")

	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{
			name:     "File exists",
			filename: testFilePath,
			expected: true,
		},
		{
			name:     "File does not exist",
			filename: nonexistentFilePath,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := FileExists(test.filename)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}
