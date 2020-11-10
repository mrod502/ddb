package ddb

import (
	"bytes"
	"crypto/tls"
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
	hub         *Hub
)

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	userHomeDir, _ = os.UserHomeDir()
	err := godotenv.Load(path.Join(userHomeDir, ".env"))
	if err != nil {
		panic(err)
	}
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
				logger.Error("GC", err.Error())
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
		logger.Error("CloseDB", err.Error())
	}
}

//ServeTLS - serve the DB
func ServeTLS() {

	var router = mux.NewRouter()
	router.HandleFunc("/", handleRequests) //.Methods("POST")

	// do we enable wss streaming?
	if os.Getenv("DB_ENABLE_WSS") == "Y" {
		hub = newHub(func(b []byte) error {
			if bytes.Contains(b, []byte("unsubscribe")) {
				return ErrUnsubscribe
			}
			return nil
		})
		router.HandleFunc("/wss", func(w http.ResponseWriter, r *http.Request) {
			wssServe(hub, w, r)
		})
	}

	for {
		err := http.ListenAndServeTLS(os.Getenv("PORT_ADDR"), os.Getenv("CERT_FILE_PATH"), os.Getenv("KEY_FILE_PATH"), router)
		if err != nil {
			logger.Error("Serve", err.Error())
			time.Sleep(time.Second)
		}
	}
}

//ServeDistributed -- future distributed server func
func ServeDistributed() {}

func pushUpdates() {}
