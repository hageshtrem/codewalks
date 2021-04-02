package md5sum

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
)

type Entry struct {
	Hash     string
	Filename string
}

// Md5all gets the path to the folder and returns an array of entries sorted by name or error.
func Md5all(dir string) ([]Entry, error) {
	done := make(chan struct{})
	defer close(done)

	names, errc := walkFiles(done, dir)

	resultChan := make(chan result) // HL
	var wg sync.WaitGroup
	numDigesters := runtime.NumCPU()
	wg.Add(numDigesters)
	for i := 0; i < numDigesters; i++ {
		go func() {
			defer wg.Done()
			digester(done, names, resultChan) // HL
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan) // HL
	}()

	r := map[string]string{}
	keys := []string{}
	for v := range resultChan {
		if v.err != nil {
			return nil, v.err
		}
		keys = append(keys, v.path)
		r[v.path] = fmt.Sprintf("%x", v.hash)
	}

	if err := <-errc; err != nil {
		return nil, err
	}

	sort.Strings(keys)
	result := []Entry{}
	for _, v := range keys {
		result = append(result, Entry{r[v], v})
	}

	return result, nil
}

// walkFiles starts a goroutine to walk the directory tree and send the path of each
// regular files on the string channel. It sends the result of the walk on the error
// channel. if done is closed , walkFiles abandons its work.
func walkFiles(done <-chan struct{}, dir string) (<-chan string, <-chan error) {
	filenames := make(chan string)
	errc := make(chan error, 1)

	go func() {
		defer close(filenames)

		var walkFn fs.WalkDirFunc = func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.Type().IsRegular() {
				return nil
			}
			select {
			case filenames <- path:
			case <-done:
				return fmt.Errorf("walk canceled")
			}
			return nil
		}
		errc <- filepath.WalkDir(dir, walkFn)
	}()

	return filenames, errc
}

// A result is the product of reading and summing a file using MD5.
type result struct {
	path string
	hash [md5.Size]byte
	err  error
}

// digester reads path names from paths and sends of the corresponding files on resultsuntil
// either paths or done closed.
func digester(done <-chan struct{}, paths <-chan string, results chan<- result) {
	for path := range paths {
		data, err := ioutil.ReadFile(path)
		select {
		case results <- result{path, md5.Sum(data), err}:
		case <-done:
			return
		}
	}
}
