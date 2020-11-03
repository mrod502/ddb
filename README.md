# Overview

# client
Use the functions in `api.go`.

# server example
```go
func main(){
    OpenDB() 
    defer CloseDB()

    go ServeTLS()
    //Wait for crtl+c
}
```

# Client example
```go
type Value struct{
    UID string
    Val int32
    Words []string
    Cfg struct{
        A string
        B string
        }
}

}
func (v Value) ID() []byte{
    return []byte(v.UID)
}

var val Value
GetPopulate([]byte("some_key"),&val)

val.Val = 42
Set(val)
```

