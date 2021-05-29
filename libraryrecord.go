package main

import (
	"crypto/md5"
	"encoding/csv"
	"io"
	"os"
	"strings"
)

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



type LibraryRecordVistorF func(lr *LibraryRecord) error 


func ForEachLibraryRecordInCSVFile(path string, vf LibraryRecordVistorF) error {
	bookFile, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	defer bookFile.Close()
	r := csv.NewReader(bookFile)
	for {
		record, err := NewLibraryRecordFromCSV(r)
		if err != nil {
			return err
		}
		if record == nil {
			// This signifies EOF
			break
		}
		if record.Titel == "Titel" {
			// This is the heading row in the CSV file
			continue
		}
		if err = vf(record); err != nil {
			return err
		}
	}
	return nil
}

type Predicate func(lr *LibraryRecord) bool

func WasPrinted(checksums *ChecksumList) Predicate {
	return func(lr *LibraryRecord) bool {
		sum := lr.MD5()
		return checksums.Contains(sum)
	}
}

func WasNotPrinted(checksums *ChecksumList) Predicate {
	sub := WasPrinted(checksums)
	return func(lr *LibraryRecord) bool {
		return !sub(lr)
	}
}


func WriteLibraryCard(predicate Predicate, w io.Writer, csl *ChecksumList) LibraryRecordVistorF {
	
	return func(lr *LibraryRecord) error {
		if predicate(lr) {
			if err := writeCard(w, lr); err != nil {
				return err
			}
			csl.Add(lr.MD5())
		}
		return nil
	}
}
