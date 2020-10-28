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
)

//Action - get, set, update, or delete value at key
type Action struct {
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
