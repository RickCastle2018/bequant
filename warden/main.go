package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Fsym string

type Tsym struct {
	TOSYMBOL        string
	CHANGE24HOUR    string
	CHANGEPCT24HOUR string
	OPEN24HOUR      string
	VOLUME24HOUR    string
	VOLUME24HOURTO  string
	LOW24HOUR       string
	HIGH24HOUR      string
	PRICE           string
	SUPPLY          string
	MKTCAP          string
}

var schema = `
CREATE TABLE IF NOT EXISTS pairs (
    fsym varchar(10) PRIMARY KEY,
    tsym varchar(10),
    raw text,
	display text
)`

type Responce struct {
	RAW     map[Fsym]Tsym
	DISPLAY map[Fsym]Tsym
}

func updateData(db *sqlx.DB) {
	fsyms := "?fsyms=" + os.Getenv("FSYMS")
	tsyms := "&tsyms=" + os.Getenv("TSYMS")
	url := os.Getenv("API_URL") + fsyms + tsyms

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var r Responce
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Fatalln(err)
	}

	var pairMaps []map[string]interface{}
	for fsym, tsym := range r.RAW {
		raw, _ := json.Marshal(tsym)
		display, _ := json.Marshal(r.DISPLAY[fsym])
		pair := map[string]interface{}{
			"fsym":    string(fsym),
			"tsym":    tsym.TOSYMBOL,
			"raw":     string(raw),
			"display": string(display),
		}
		pairMaps = append(pairMaps, pair)
	}

	sql := `
		REPLACE INTO pairs (fsym, tsym, raw, display) 
		VALUES (:fsym, :tsym, :raw, :display)
	`
	_, err = db.NamedExec(sql, pairMaps)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("warden - pairs updated")
}

func main() {
	log.Println("warden is starting")

	url := os.Getenv("DB_USER") + ":" + os.Getenv("DB_USER_PASSWORD") + "@tcp(db)/" + os.Getenv("DB_NAME")
	db, err := sqlx.Open("mysql", url)
	if err != nil {
		log.Fatalln(err)
	}
	db.MustExec(schema)

	intervalInt, err := strconv.ParseInt(os.Getenv("UPDATE_INTERVAL"), 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	interval := time.Duration(intervalInt)
	ticker := time.NewTicker(interval * time.Second)
	for {
		<-ticker.C
		updateData(db)
	}
}
