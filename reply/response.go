package reply

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	c "github.com/SpeedSlime/Asahi/config"
)

func RespondWithError(w http.ResponseWriter, err error) error {
	req, ok := err.(*RequestError)
	if ok { RespondWithJSON(w, req.Code(), map[string]string{"error": http.StatusText(req.Code())}) } else { RespondWithJSON(w, 500, map[string]string{"error": http.StatusText(500)}) }
	if req.Code() - 500 >= 0 && req.Code() - 500 <= 99 { return fmt.Errorf("RespondWithError: %w", err) }
	return nil
}

func RespondWithFile(w http.ResponseWriter, httpcode int, path string) error {
	bin, err := open(c.Assets() + filepath.Clean(path));
	if err != nil { return RespondWithError(w, err) }
	RespondWithDetect(w, httpcode, bin)
	return nil
}

func RespondWithCode(w http.ResponseWriter, httpcode int, code string) error {
	err := RespondWithJSON(w, httpcode, map[string]string{"result": code})
	if err != nil { return err }
	return nil
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	res, err := json.Marshal(payload)
	if err != nil { return RespondWithError(w, &RequestError{Function: "RespondWithJSON", StatusCode: 500, Err: err,}) }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
	return nil
}

func RespondWithDetect(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Set("Content-Type", http.DetectContentType(payload))
	w.WriteHeader(code)
	w.Write(payload)
}

func RespondWithRaw(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(code)
	w.Write(payload)
}

func open(path string) ([]byte, error) {
    bin, err := os.Open(path)
	if os.IsNotExist(err) { return nil, &RequestError{Function: "open", StatusCode: 404, Err: err,} }
	out, err := ioutil.ReadAll(bin)
    if err == nil { return out, nil }
    return nil, &RequestError{Function: "open", StatusCode: 500, Err: err,}
}
