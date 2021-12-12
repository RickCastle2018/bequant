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
	CHANGE24HOUR    interface{}
	CHANGEPCT24HOUR interface{}
	OPEN24HOUR      interface{}
	VOLUME24HOUR    interface{}
	VOLUME24HOURTO  interface{}
	LOW24HOUR       interface{}
	HIGH24HOUR      interface{}
	PRICE           interface{}
	SUPPLY          interface{}
	MKTCAP          interface{}
}

var schema = `
CREATE TABLE pairs (
    fsym varchar(10),
    tsym varchar(10),
    raw text,
	display text,
	PRIMARY KEY (fsym, tsym)
)`

type Responce struct {
	RAW     map[Fsym]map[string]Tsym
	DISPLAY map[Fsym]map[string]Tsym
}

func updateData(db *sqlx.DB, sql string) error {
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
	for fsym, tsyms := range r.RAW {
		for tsym, data := range tsyms {
			raw, _ := json.Marshal(data)
			display, _ := json.Marshal(r.DISPLAY[fsym][tsym])
			pair := map[string]interface{}{
				"fsym":    string(fsym),
				"tsym":    tsym,
				"raw":     string(raw),
				"display": string(display),
			}
			pairMaps = append(pairMaps, pair)
		}
	}

	_, err = db.NamedExec(sql, pairMaps)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("warden - pairs updated")
	return nil
}

func main() {
	log.Println("warden is starting")

	url := os.Getenv("DB_USER") + ":" + os.Getenv("DB_USER_PASSWORD") + "@tcp(db)/" + os.Getenv("DB_NAME")
	db, err := sqlx.Open("mysql", url)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(schema)
	if err == nil {
		err = updateData(db, `
		INSERT INTO pairs (fsym, tsym, raw, display)
		VALUES (:fsym, :tsym, :raw, :display)
		`)
		if err != nil {
			log.Fatalln(err)
		}
	}

	intervalInt, err := strconv.ParseInt(os.Getenv("UPDATE_INTERVAL"), 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	interval := time.Duration(intervalInt)
	ticker := time.NewTicker(interval * time.Second)
	for {
		<-ticker.C
		updateData(db, `
			REPLACE INTO pairs (fsym, tsym, raw, display)
			VALUES (:fsym, :tsym, :raw, :display)
		`)
	}
}
