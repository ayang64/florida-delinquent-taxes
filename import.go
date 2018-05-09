package main

import (
	"database/sql"
	"encoding/csv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strings"
	"unicode"
)

func main() {
	db, err := sql.Open("postgres", "database=taxes user=ayan host=/var/run/postgresql")

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(os.Stdin)
	record, err := reader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare("insert into taxes (id, business, owner, location, county, warrant, amount) values (default, NULLIF($1, ''), NULLIF($2, ''), $3, $4, $5, $6::money);")

	if err != nil {
		log.Fatal(err)
	}

	if len(record) == 0 {
		log.Fatalf("no records to import.")
	}

	ifs := make([]interface{}, len(record[0]), len(record[0]))

	for _, rec := range record {
		for i := range rec {
			rec[i] = strings.TrimSpace(rec[i])
		}
		log.Printf("%#v", rec)

		// delete all internal space from money column.  some of the money fields
		// have weird unicode \00a0 characters between the dollar sign and the
		// amount. :(
		rec[5] = func() string {
			rc := ""
			for _, i := range rec[5] {
				if unicode.IsSpace(i) {
					continue
				}
				rc += string(i)
			}

			return rc
		}()

		// ensure that the "amount" column does not end with a period!  someone
		// type-o'd an entry in the spreadsheet.
		if len(rec[5]) > 0 && rec[5][len(rec[5])-1] == '.' {
			rec[5] = "$" + rec[5]
			rec[5] = rec[5][:len(rec[5])-2]
		}

		for i := range rec {
			ifs[i] = &rec[i]
		}

		_, err := stmt.Exec(ifs...)

		if err != nil {
			log.Fatal(err)
		}

	}
}
