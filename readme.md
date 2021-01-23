# parse-git-index

parse git index


```go
func main() {
	indexFile, _ := filepath.Abs(filepath.Dir("") + "/example/index")
	header, entries, e := parse_git_index.ParseGitIndex(indexFile)
	if e != nil {
		log.Fatal(e)
	}
	fmt.Println(header)
	fmt.Println(entries)

}
```



```go
type Header struct {
	Signature string `json:"signature"`
	Version   uint32 `json:"version"`
	Count     uint32 `json:"count"`
}

type Entry struct {
	Ctime    uint32 `json:"ctime"`
	Mtime    uint32 `json:"mtime"`
	Dev      uint32 `json:"dev"`
	Ino      uint32 `json:"ino"`
	Mode     uint32 `json:"mode"`
	Uid      uint32 `json:"uid"`
	Gid      uint32 `json:"gid"`
	FileSize uint32 `json:"fileSize"`
	ObjectId string `json:"objectId"`
	FilePath string `json:"filePath"`
}
```