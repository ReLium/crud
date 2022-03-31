package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type IO struct {
}

func NewIO() *IO {
	return &IO{}
}

func (i *IO) writeJSON(w http.ResponseWriter, status int, v interface{}) error {
	resp, err := json.Marshal(v)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(resp)

	return err
}

func (i *IO) readJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

type errorResponse struct {
	Message string
}

func (i *IO) writeError(w http.ResponseWriter, err error) error {
	errorResponse := &errorResponse{
		Message: err.Error(),
	}
	return i.writeJSON(w, http.StatusInternalServerError, errorResponse)
}
