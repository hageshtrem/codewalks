package main

import (
	"fmt"
	"os"

	"github.com/hageshtrem/codewalks/md5sum"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage:\n\tmd5sum dir\n")
		return
	}
	entries, err := md5sum.Md5all(os.Args[1])
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	for _, e := range entries {
		fmt.Printf("%s    %s\n", e.Hash, e.Filename)
	}
}
