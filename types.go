package ddb

//DB action types
const (
	ActionSet int = iota
	ActionGet
	ActionDel
	ActionUpd
)

//Action statuses
const (
	StatusOK int = iota
	StatusFailed
	StatusRejected
	StatusFailedUnmarshal
)

//Action - get, set, update, or delete value at key
type Action struct {
	_msgpack   struct{} `msgpack:",omitempty"`
	APIKey     string
	ActionType int
	Key        []byte
	Value      []byte
}

//Result - result of an action
type Result struct {
	Status int
	Data   []byte
}

//Indexer - an indexable object that returns a unique ID
type Indexer interface {
	ID() []byte
	Type() []byte
}

//Key in database
func Key(i Indexer) (b []byte) {
	b = append(b, i.Type()...)
	b = append(b, byte('^'))
	b = append(b, i.ID()...)
	return b
}
