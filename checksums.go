package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"io"
	"os"

)

type Checksum []byte
type ChecksumList []Checksum



func NewChecksumListFromFile(path string) (*ChecksumList, error) {
	checksumFile, err := os.Open(path)
	if os.IsNotExist(err) {
		checksumFile, err = os.Create(path)
	}
	if err != nil {
		return nil, err
	}
	defer checksumFile.Close()
	checksums := &ChecksumList{}
	if err := checksums.Load(checksumFile); err != nil {
		return nil, err
	}
	return checksums, nil
}

func DumpChecksumsToFile(checksums *ChecksumList, path string) error {
	checksumFile, err := os.OpenFile(os.Args[2], os.O_TRUNC | os.O_WRONLY , 0660)
	if err != nil {
		return err
	}
	defer checksumFile.Close()

	if err = checksums.Save(checksumFile); err != nil {
		return err
	}
	if err := checksumFile.Sync(); err != nil {
		return err
	}	
	return nil
}

func (csl *ChecksumList) Load(r io.Reader) error {
	cr := bufio.NewReader(r)
	for {
		line, _, err := cr.ReadLine()
		if err == io.EOF {
			return nil
		}
			
		if err != nil {
			return err
		}
		sum, err := base64.StdEncoding.DecodeString(string(line))
		if err != nil {
			return err
		}
		*csl = append(*csl, sum)
	}
	return nil
}

func (csl *ChecksumList) Save(w io.Writer) error {
	for _, sum := range *csl {
		_, err := io.WriteString(w, base64.StdEncoding.EncodeToString(sum) + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func (csl *ChecksumList) Contains(sum []byte) bool {
	for _, csum := range *csl {
		if bytes.Equal(sum, csum) {
			return true
		}
	}
	return false
}


func (csl *ChecksumList) Add(sum []byte) {
	*csl = append(*csl, sum)
}
