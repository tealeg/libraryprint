package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func usage() {
	log.Fatal("libraryprint csv-file print-log-file\n")
}

type LibraryRecord struct {
	Titel             string
	Autoren           string
	Schriftenreihe    string
	Kategorien        string
	Publikationsdatum string
	Verlag            string
	Seiten            string
	ISBN              string
	Gelesen           string
	Lesezeiten        string
	ArchiveNummer     string
	Zusammenfassung   string
	CoverPfad         string
}




func NewLibraryRecordFromCSV(r *csv.Reader) (*LibraryRecord, error) {
	line, err := r.Read()
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	record := &LibraryRecord{
		Titel:             line[0],
		Autoren:           line[1],
		Schriftenreihe:    line[2],
		Kategorien:        line[3],
		Publikationsdatum: line[4],
		Verlag:            line[5],
		Seiten:            line[6],
		ISBN:              line[7],
		Gelesen:           line[8],
		Lesezeiten:        line[9],
		Zusammenfassung:   line[11],
		CoverPfad:         line[12],
	}

	commentary := line[10]
	if strings.Contains(commentary, ":") {
		record.ArchiveNummer = strings.Split(commentary, ":")[1]
	}
	return record, nil
}


func (lr *LibraryRecord) MD5() ([]byte) {
	sum := md5.New()
	io.WriteString(sum, lr.Titel)
	io.WriteString(sum, lr.Autoren)
	io.WriteString(sum, lr.Schriftenreihe)
	io.WriteString(sum, lr.Kategorien)
	io.WriteString(sum, lr.Publikationsdatum)
	io.WriteString(sum, lr.Verlag)
	io.WriteString(sum, lr.Seiten)
	io.WriteString(sum, lr.ISBN)
	io.WriteString(sum, lr.Gelesen)
	io.WriteString(sum, lr.Lesezeiten)
	io.WriteString(sum, lr.Zusammenfassung)
	io.WriteString(sum, lr.CoverPfad)
	io.WriteString(sum, lr.ArchiveNummer)
	return sum.Sum(nil)
}

func (lr *LibraryRecord) WasPrinted(checkSums [][]byte) bool {
	sum := lr.MD5()
	for _, csum := range checkSums {
		if bytes.Equal(sum, csum) {
			return true
		}
	}
	return false
}

func (lr *LibraryRecord) Print() {
	return 
}


func loadCheckSums(r io.Reader) ([][]byte, error) {
	cr := bufio.NewReader(r)
	sums := make([][]byte, 0)
	for {
		line, _, err := cr.ReadLine()
		if err == io.EOF {
			return sums, nil
		}
			
		if err != nil {
			return nil, err
		}
		sum, err := base64.StdEncoding.DecodeString(string(line))
		if err != nil {
			return nil, err
		}
		sums = append(sums, sum)
	}
	return sums, nil
}


func main() {
	if len(os.Args) != 3 {
		usage()

	}
	checksumFile, err := os.Open(os.Args[2])
	if os.IsNotExist(err) {
		checksumFile, err = os.Create(os.Args[2])
	}
	if err != nil {
		log.Fatal(err)
	}

	checksums, err := loadCheckSums(checksumFile)
	if err != nil {
		log.Fatal(err)
	}
	checksumFile.Close()
	
	bookFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	r := csv.NewReader(bookFile)
	for {
		record, err := NewLibraryRecordFromCSV(r)
		if err != nil {
			log.Fatal(err)
		}
		if record == nil {
			// This signifies EOF
			break
		}
		if !record.WasPrinted(checksums) {
			fmt.Printf("%+v\n", record)
			checksums = append(checksums, record.MD5())
		}
	}
	checksumFile, err = os.OpenFile(os.Args[2], os.O_TRUNC | os.O_WRONLY , 0660)
	for _, sum := range checksums {
		checksumFile.WriteString(base64.StdEncoding.EncodeToString(sum))
		checksumFile.WriteString("\n")
		err = checksumFile.Sync()
			if err != nil {
		log.Fatal(err)
	}

	}


	checksumFile.Close()
}
