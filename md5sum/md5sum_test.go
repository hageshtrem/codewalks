package md5sum

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// The io operations with files are not mocked and are tested as they are.
// See examples in io/ioutil/ioutil_test.go.

func TestMd5all(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	if err := ioutil.WriteFile(filepath.Join(dir, "aa"), []byte("Test content1"), 0770); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(dir, "bb"), []byte("Test content2"), 0770); err != nil {
		t.Fatal(err)
	}

	result, err := Md5all(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 &&
		(result[0].Hash != "ed55f907c8accf85ed08bf0806cc130a" ||
			result[1].Hash != "d13dc04914aac361789528789c8ab2da") {
		t.Fatal("Error while walking and calculating")
	}
}
