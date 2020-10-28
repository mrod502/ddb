package ddb

import (
	"os"

	"github.com/vmihailenco/msgpack/v5"
)

//Get -
func Get(key []byte) (r Result) {
	var a Action

	a.APIKey = os.Getenv("API_KEY")
	a.Key = key
	a.ActionType = ActionGet

	doAPIRequest(a, &r)

	return
}

//Set - set value at key
func Set(ix Indexer) (r Result) {
	var a Action
	b, err := msgpack.Marshal(ix)
	if err != nil {
		r.Status = StatusFailed
		return
	}
	a.Value = b
	a.APIKey = os.Getenv("API_KEY")
	a.Key = ix.ID()
	a.ActionType = ActionSet

	doAPIRequest(a, &r)

	return
}

//Del - delete items with prefix key
func Del(key []byte) (r Result) {
	var a Action

	if len(key) == 0 {
		r.Status = StatusFailed
		return
	}
	a.Key = key
	a.APIKey = os.Getenv("API_KEY")

	a.ActionType = ActionDel

	doAPIRequest(a, &r)

	return
}
