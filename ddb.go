package ddb

import (
	"encoding/json"
	"net/http"
	"time"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/gorilla/mux"
	"github.com/mrod502/logger"
)

const ()

var (
	db     *badger.DB
	config DBOptions
)

//OpenDB - open the database at path
func OpenDB(options []byte) (err error) {
	err = json.Unmarshal(options, &config)
	if err != nil {
		logger.Log("DB", err.Error())
	}
	db, err = badger.Open(badger.DefaultOptions(config.Path))

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(max(config.GCInterval, 60)))
			db.RunValueLogGC(0.7)

		}
	}()
	return
}

//DBOptions - options for configuring database
type DBOptions struct {
	Path       string
	GCInterval int64
	Servers    []string
}

func max(in ...int64) int64 {
	var max = in[0]
	for _, v := range in {
		if v > max {
			max = v
		}
	}
	return max
}

//CloseDB - closes the db
func CloseDB() {
	err := db.Close()
	if err != nil {
		logger.Log("CloseDB", err.Error())
	}
}

func Serve() {
	var router = mux.NewRouter()
	router.HandleFunc("/", handleRequests).Methods("POST")
}

func ServeDistributed() {}

func handleRequests(w http.ResponseWriter, r *http.Request) {

}

type DBAction struct {
	APIKey     string
	ActionType int
	Key        []byte
	Value      []byte
}
