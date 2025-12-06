package main

import (
	"io/fs"
	"os"
	"reflect"
	"sort"
	"testing"
)

type MockDirectory struct {
	name  string
	isdir bool
}

func (m MockDirectory) Name() string               { return m.name }
func (m MockDirectory) IsDir() bool                { return m.isdir }
func (m MockDirectory) Type() fs.FileMode          { return 0 }
func (m MockDirectory) Info() (fs.FileInfo, error) { return nil, nil }

func TestFileSorting(t *testing.T) {
	t.Parallel()
	TestCases := []struct {
		InputDirs      []os.DirEntry
		ExpectedOutput map[string][]string
	}{
		{
			InputDirs: []os.DirEntry{
				MockDirectory{name: "1.txt", isdir: false},
				MockDirectory{name: "2.txt", isdir: false},
				MockDirectory{name: ".hello", isdir: false},
				MockDirectory{name: "..hello", isdir: false},
				MockDirectory{name: "hello", isdir: false},
				MockDirectory{name: ".hello.vim", isdir: false},
				MockDirectory{name: "archive.tar.gz", isdir: false},
			},
			ExpectedOutput: map[string][]string{
				"txt":           {"2.txt", "1.txt"},
				"hidden":        {".hello", "..hello", ".hello.vim"},
				"uncategorized": {"hello"},
				"gz":            {"archive.tar.gz"},
			},
		},
	}

	for _, tests := range TestCases {
		sortedFiles := SortFiles(tests.InputDirs)
		t.Log(sortedFiles)
		for k, v := range tests.ExpectedOutput {
			if sortedFiles[k] == nil {
				t.Fatal("expected file extension not found in the sorted files", k)
			}
			sort.Strings(v)
			sort.Strings(sortedFiles[k])
			if !reflect.DeepEqual(v, sortedFiles[k]) {
				t.Error("the expected sorting is not same as the one done by the sorter expected", v, "output", sortedFiles[k])
			}
		}
	}
}
