package reply

import (
	"io/ioutil"
	"os"
)


func open(path string) ([]byte, error) {
    bin, err := os.Open(path)
	if os.IsNotExist(err) { return nil, &RequestError{Function: "open", StatusCode: 404, Err: err,} }
	out, err := ioutil.ReadAll(bin)
    if err == nil { return out, nil }
    return nil, &RequestError{Function: "open", StatusCode: 500, Err: err,}
}