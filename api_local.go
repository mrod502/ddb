package ddb

import (
	"os"

	"github.com/vmihailenco/msgpack/v5"
)

//GetLoc -
func GetLoc(key []byte) (r Result) {
	var a Action

	a.APIKey = os.Getenv("API_KEY")
	a.Key = key
	a.ActionType = ActionGet

	r = get(a)

	return
}

//GetPopulateLoc -
func GetPopulateLoc(key []byte, obj interface{}) (r Result) {
	var a Action

	a.APIKey = os.Getenv("API_KEY")
	a.Key = key
	a.ActionType = ActionGet

	r = get(a)

	err := msgpack.Unmarshal(r.Data, obj)

	if err != nil {
		r.Status = StatusFailedUnmarshal
	}

	return
}

//SetLoc - set value at key
func SetLoc(ix Indexer) (r Result) {
	var a Action
	b, err := msgpack.Marshal(ix)
	if err != nil {
		r.Status = StatusFailed
		return
	}
	a.Value = b
	a.APIKey = os.Getenv("API_KEY")
	a.Key = Key(ix)
	a.ActionType = ActionSet

	r = set(a)

	return
}

//DelLoc - delete items with prefix key
func DelLoc(key []byte) (r Result) {
	var a Action

	if len(key) == 0 {
		r.Status = StatusFailed
		return
	}
	a.Key = key
	a.APIKey = os.Getenv("API_KEY")

	a.ActionType = ActionDel

	r = del(a)

	return
}
