package function

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type response struct {
	Status  int         `json:"status"`
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	Message string      `json:"message"`
	Headers http.Header `json:"headers"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	var input string

	if r.Body != nil {
		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		input = string(body)
	}

	prefix, status := parsePath(r)

	resp, err := json.Marshal(response{
		Status:  status,
		Method:  r.Method,
		Path:    prefix,
		Message: input,
		Headers: r.Header,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)
	_ = w.Write(resp)
}

func parsePath(r *http.Request) (prefix string, status int) {
	path := r.URL.Path
	parts := strings.Split(path, "/")

	prefix = strings.Join(parts[:len(parts)-1], "/")

	potentialStatus := parts[len(parts)-1]
	status, err := strconv.Atoi(potentialStatus)
	if err != nil {
		prefix = prefix + "/" + potentialStatus
		status = 200
	}

	if status < 100 || status > 599 {
		prefix = prefix + "/" + potentialStatus
		status = 200
	}

	return prefix, status
}
