package ddb

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

func TestServer(t *testing.T) {

	err := OpenDB()
	if err != nil {
		t.Fatal(err)
	}
	defer CloseDB()
	go ServeTLS()
	fmt.Println(os.Getenv("SERVE_ADDR"))
	time.Sleep(3 * time.Second)
	ts := testStruct{UID: "test-id", Data: []byte("this is a test")}

	res := Set(ts)
	fmt.Printf("%+v\n", res)
	if res.Status != StatusOK {
		t.Fatal("write failed")
	}
	r := Get(ts.ID())

	var rstruct testStruct

	err = msgpack.Unmarshal(r.Data, &rstruct)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", rstruct)

}

type testStruct struct {
	UID  string
	Data []byte
}

func (t testStruct) ID() []byte {
	return []byte(t.UID)
}
