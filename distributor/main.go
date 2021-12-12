package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// https://stackoverflow.com/a/44403016
func Parallelize(functions ...func()) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(functions))

	defer waitGroup.Wait()

	for _, function := range functions {
		go func(copy func()) {
			defer waitGroup.Done()
			copy()
		}(function)
	}
}

var db *sqlx.DB

type Pair struct {
	fsym    string
	tsym    string
	raw     string
	display string
}

func handler(w http.ResponseWriter, r *http.Request) {
	fsymsString := r.URL.Query()["fsyms"][0]
	fsyms := strings.Split(fsymsString, ",")

	tsymsString := r.URL.Query()["tsyms"][0]
	tsyms := strings.Split(tsymsString, ",")

	for _, fsym := range fsyms {
		sql := `
		SELECT tsym, raw, display FROM pairs
		WHERE fsym=:fsym AND tsym in (:tsyms)
		`
		params := map[string]interface{}{
			"fsym":  fsym,
			"tsyms": "'" + strings.Join(tsyms, "', '") + "'",
		}
		rows, err := db.NamedQuery(sql, params)
		if err != nil {
			log.Fatalln(err)
		}

		var pairs []Pair
		for rows.Next() {
			var pair Pair
			err := rows.StructScan(&pair)
			if err != nil {
				log.Fatalln(err)
			}
			pairs = append(pairs, pair)
			log.Println(pair)
		}
		log.Println(pairs)
	}
}

func main() {
	log.Println("distributor is starting at localhost:8080")

	url := os.Getenv("DB_USER") + ":" + os.Getenv("DB_USER_PASSWORD") + "@tcp(db)/" + os.Getenv("DB_NAME")
	var err error
	db, err = sqlx.Open("mysql", url)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/price", handler)
	http.ListenAndServe(":8080", nil)
}
