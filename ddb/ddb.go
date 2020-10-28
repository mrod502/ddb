package ddb

import (
	"net/http"
	"os"
	"path"
	"time"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/gorilla/mux"
	"github.com/mrod502/logger"

	"github.com/joho/godotenv"
)

const ()

var (
	db          *badger.DB
	userHomeDir string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	userHomeDir, _ = os.UserHomeDir()

}

//OpenDB - open the database at path
func OpenDB() (err error) {

	db, err = badger.Open(badger.DefaultOptions(path.Join(userHomeDir, os.Getenv("DB_PATH"))))
	if err != nil {
		return err
	}
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			err = db.RunValueLogGC(0.7)
			if err != nil {
				logger.Log("GC", err.Error())
			}

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

//ServeTLS - serve the DB
func ServeTLS() {
	var router = mux.NewRouter()
	router.HandleFunc("/", handleRequests).Methods("POST")
	for {
		err := http.ListenAndServeTLS(os.Getenv("SERVE_ADDR"), os.Getenv("CERT_FILE_PATH"), os.Getenv("KEY_FILE_PATH"), router)
		if err != nil {
			logger.Log("Serve", err.Error())
		}
	}
}

//ServeDistributed -- future distributed server func
func ServeDistributed() {}
