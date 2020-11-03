package ddb

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dgraph-io/badger/v2"
	"github.com/mrod502/logger"
	msgpack "github.com/vmihailenco/msgpack/v5"
)

// --------------- API request helper funcs --------------
func verifyAPIKey(a Action) bool {
	apiKey := os.Getenv("API_KEY")
	if a.APIKey == apiKey || apiKey == "" {
		return true
	}
	return false
}

// handle the outgoing action
func doAPIRequest(a Action, r *Result) {
	b, err := msgpack.Marshal(a)

	if err != nil {
		r.Status = StatusFailed
		logger.Error("Get", err.Error())
		return
	}
	req, err := http.NewRequest("POST", os.Getenv("SERVE_ADDR"), bytes.NewBuffer(b))

	if err != nil {
		r.Status = StatusFailed
		logger.Error("API", "create request", err.Error())
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		r.Status = StatusFailed
		logger.Error("API", "do request", err.Error())
		return
	}
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		r.Status = StatusFailed
		logger.Error("API", "read response", err.Error())
		return
	}
	defer res.Body.Close()
	err = msgpack.Unmarshal(b, r)
	if err != nil {
		r.Status = StatusFailed
		logger.Error("API", "unmarshal", err.Error())
		fmt.Println("RESPONSE:", string(b))
		return
	}
}

//handle incoming API requests
func handleRequests(w http.ResponseWriter, r *http.Request) {
	var result Result
	var resBytes []byte

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		result.Status = StatusFailed
		resBytes, _ = msgpack.Marshal(result)
		w.Write(resBytes)
		return
	}

	defer r.Body.Close()
	var action Action

	if err = msgpack.Unmarshal(b, &action); err != nil {
		result.Status = StatusFailed
		resBytes, _ = msgpack.Marshal(result)
		w.Write(resBytes)
		return
	}

	if !verifyAPIKey(action) {
		result.Status = StatusRejected
		logger.Error("API", "bad Key", action.APIKey)
		w.Write(resBytes)
		return
	}

	switch action.ActionType {
	case ActionSet:
		result = set(action)
	case ActionGet:
		result = get(action)
	case ActionDel:
		if os.Getenv("ALLOW_DELETES") != "Y" {
			result.Status = StatusRejected
			break
		}
		result = del(action)
	case ActionUpd:
		result = upd(action)
	}
	b, err = msgpack.Marshal(result)
	if err != nil {
		logger.Error("DB", "send", err.Error())
	}
	w.Write(b)

}

// ------------------- DB funcs -------------------
func set(a Action) (r Result) {
	err := db.Update(func(txn *badger.Txn) (err error) {
		return txn.Set(a.Key, a.Value)
	})
	if err != nil {
		logger.Error("Update", err.Error())
		r.Status = StatusFailed
	}
	return
}
func get(a Action) (r Result) {
	err := db.View(func(txn *badger.Txn) (err error) {
		i, err := txn.Get(a.Key)
		if err != nil {
			return err
		}
		i.Value(func(b []byte) error {
			r.Data = b
			return nil
		})
		return nil
	})
	if err != nil {
		r.Status = StatusFailed
		logger.Error("Get", err.Error())
	}

	return
}

func del(a Action) (r Result) {
	if len(a.Key) == 0 {
		r.Status = StatusFailed
		return
	}
	err := db.DropPrefix(a.Key)
	if err != nil {
		logger.Error("Del", err.Error())
		r.Status = StatusFailed
	}
	return
}

func upd(a Action) (r Result) {
	err := db.Update(func(txn *badger.Txn) (err error) {
		return txn.Set(a.Key, a.Value)
	})
	if err != nil {
		logger.Error("Update", err.Error())
		r.Status = StatusFailed
	}
	return
}
