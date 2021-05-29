package main

import (
	"log"
	"os"
)




func usage() {
	log.Fatal("libraryprint csv-file print-log-file output-latex-file-name \n")
}


func main() {
	if len(os.Args) != 4 {
		usage()

	}
	checksums, err := NewChecksumListFromFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
        }

	output, err := os.OpenFile(os.Args[3], os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0660)
	if err != nil {
		log.Fatal(err)
	}

	if err = writePreamble(output); err != nil {
		log.Fatal(err)
	}

	err = ForEachLibraryRecordInCSVFile(
		os.Args[1],
		WriteLibraryCard(
			WasNotPrinted(checksums),
			output,
			checksums),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = writePostscript(output); err != nil {
		log.Fatal(err)
	}

	if err := output.Sync(); err != nil {
		log.Fatal(err)
	}
	if err := output.Close(); err != nil {
		log.Fatal(err)
	}

	if err := DumpChecksumsToFile(checksums, os.Args[2]); err != nil {
		log.Fatal(err)
	}
	log.Printf(`Output stored in %q, checksums in %q

Now run:

    xelatex %q

To generate a pdf.`, os.Args[3], os.Args[2], os.Args[3])
	
		
}
