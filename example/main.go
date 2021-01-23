package main

import (
	parse_git_index "github.com/zhaopengme/parse-git-index"
	"log"
	"path/filepath"
)

func main() {
	indexFile, _ := filepath.Abs(filepath.Dir("") + "/example/index")
	header, entries, e := parse_git_index.ParseGitIndex(indexFile)
	if e != nil {
		log.Fatal(e)
	}
	println(header)
	println(entries)

}
