package parse_git_index

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
)

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


type header struct {
	Signature [4]byte
	Version   uint32
	Count     uint32
}

type entry struct {
	CtimeSecond     uint32
	CtimeNanosecond uint32
	MtimeSecond     uint32
	MtimeNanosecond uint32
	Dev             uint32
	Ino             uint32
	Mode            uint32
	Uid             uint32
	Gid             uint32
	FileSize        uint32
	ObjectId        [20]byte
	FilePathSize    uint16
}

func convertToHeader(header header) Header {
	h := Header{
		Signature: string(header.Signature[:]),
		Version:   header.Version,
		Count:     header.Count,
	}
	return h
}

func convertToEntry(entry entry, filePath string) Entry {
	en := Entry{
		Ctime:    entry.CtimeSecond,
		Mtime:    entry.MtimeSecond,
		Dev:      entry.Dev,
		Ino:      entry.Ino,
		Mode:     entry.Mode,
		Uid:      entry.Uid,
		Gid:      entry.Gid,
		FileSize: entry.FileSize,
		ObjectId: hex.EncodeToString(entry.ObjectId[:]),
		FilePath: filePath,
	}
	return en
}


func padEntry(index []byte, offset *int, size int) {
	var padLen int
	if (8 - (size % 8)) != 0 {
		padLen = 8 - (size % 8)
	} else {
		padLen = 8
	}

	/*断言padding字符都为0*/
	for i := 0; i < padLen; i++ {
		if index[*offset+i] != 0 {
			panic("error")
		}
	}

	*offset += padLen
}


func getBytes(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

func ParseGitIndex(indexFile string) (*Header, *[]Entry, error) {
	var headerSize = 12
	var entrySize = 62

	index := getBytes(indexFile)

	var header header
	var currentEntryBody entry
	var currentEntrySize int

	e := binary.Read(bytes.NewBuffer(index), binary.BigEndian, &header)
	if e != nil {
		return nil, nil, e
	}
	h := convertToHeader(header)

	offset := int(headerSize)
	num := int(header.Count)
	var entries []Entry
	for i := 0; i < num; i++ {
		currentEntrySize = entrySize
		binary.Read(bytes.NewBuffer(index[offset:]), binary.BigEndian, &currentEntryBody)
		offset = offset + currentEntrySize
		currentFilePathLength := int(currentEntryBody.FilePathSize & (0xffff >> 4))
		if currentFilePathLength < 0xfff {
			currentFilePath := string(index[offset : offset+currentFilePathLength])
			offset = offset + currentFilePathLength
			currentEntrySize = currentEntrySize + currentFilePathLength
			padEntry(index, &offset, currentEntrySize) /*每个entry都为8的整数倍*/
			en := convertToEntry(currentEntryBody, currentFilePath)
			entries = append(entries, en)
		}
	}
	return &h, &entries, nil
}