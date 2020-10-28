package ddb

import (
	"net/http"
	"os"
	"time"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/gorilla/mux"
	"github.com/mrod502/logger"
	msgpack "github.com/vmihailenco/msgpack/v5"

	"github.com/joho/godotenv"
)

const ()

var (
	db     *badger.DB
	config DBOptions
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

//OpenDB - open the database at path
func OpenDB(options []byte) (err error) {
	err = msgpack.Unmarshal(options, &config)
	if err != nil {
		panic(err)
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

//ServeTLS - serve the DB somewhere
func ServeTLS() {
	var router = mux.NewRouter()
	router.HandleFunc("/", handleRequests).Methods("POST")
	http.ListenAndServeTLS(os.Getenv("SERVE_ADDR"), os.Getenv("CERT_FILE_PATH"), os.Getenv("KEY_FILE_PATH"), router)
}

//ServeDistributed -- future distributed server func
func ServeDistributed() {}
