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

type Responce struct {
	RAW     map[string]map[string]string
	DISPLAY map[string]map[string]string
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
	data.RAW = make(map[string]map[string]string)
	data.DISPLAY = make(map[string]map[string]string)
	for _, fsym := range fsymsArr {
		data.RAW[fsym] = make(map[string]string)
		data.DISPLAY[fsym] = make(map[string]string)
	}

	var pair Pair
	for rows.Next() {
		err := rows.StructScan(&pair)
		if err != nil {
			log.Fatalln(err)
		}
		data.RAW[pair.Fsym][pair.Tsym] = pair.Raw
		data.DISPLAY[pair.Fsym][pair.Tsym] = pair.Display
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

	http.HandleFunc("/price", handler)
	http.ListenAndServe(":8080", nil)
}
