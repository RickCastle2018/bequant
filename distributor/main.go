package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

type Pair struct {
	Fsym    string
	Tsym    string
	Raw     string
	Display string
}

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

type Responce struct {
	RAW     map[string]map[string]Tsym
	DISPLAY map[string]map[string]Tsym
}

func prepareParam(r *http.Request, name string) (string, []string) {
	paramString := r.URL.Query()[name][0]
	paramArr := strings.Split(paramString, ",")
	param := "'" + strings.Join(paramArr, "', '") + "'"
	return param, paramArr
}

func handler(w http.ResponseWriter, r *http.Request) {
	fsyms, fsymsArr := prepareParam(r, "fsyms")
	tsyms, _ := prepareParam(r, "tsyms")

	sql := fmt.Sprintf(`
	SELECT * FROM pairs
	WHERE fsym IN (%s) AND tsym IN (%s)
	`, fsyms, tsyms)

	rows, err := db.Queryx(sql)
	if err != nil {
		log.Fatalln(err)
	}

	var data Responce
	data.RAW = make(map[string]map[string]Tsym)
	data.DISPLAY = make(map[string]map[string]Tsym)
	for _, fsym := range fsymsArr {
		data.RAW[fsym] = make(map[string]Tsym)
		data.DISPLAY[fsym] = make(map[string]Tsym)
	}

	var pair Pair
	for rows.Next() {
		err := rows.StructScan(&pair)
		if err != nil {
			log.Fatalln(err)
		}
		var raw Tsym
		var display Tsym
		err = json.Unmarshal([]byte(pair.Raw), &raw)
		if err != nil {
			log.Fatalln(err)
		}
		err = json.Unmarshal([]byte(pair.Display), &display)
		if err != nil {
			log.Fatalln(err)
		}
		data.RAW[pair.Fsym][pair.Tsym] = raw
		data.DISPLAY[pair.Fsym][pair.Tsym] = display
		// go cachePair(pair)
	}

	res, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(res))
}

func main() {
	log.Println("distributor is starting at localhost:8080")

	url := os.Getenv("DB_USER") + ":" + os.Getenv("DB_USER_PASSWORD") + "@tcp(db)/" + os.Getenv("DB_NAME")
	var err error
	db, err = sqlx.Open("mysql", url)
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxIdleConns(64)
	db.SetMaxOpenConns(64)

	http.HandleFunc("/price", handler)
	http.ListenAndServe(":8080", nil)
}
