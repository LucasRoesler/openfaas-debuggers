package function

import (
	"io/ioutil"
	"log"
	"net/http"
)

// Handle a function invocation
func Handle(w http.ResponseWriter, r *http.Request) {
	var input string

	if r.Body != nil {
		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("error reading body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		input = string(body)
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(input))
	if err != nil {
		log.Printf("error writing output: %s", err)
	}
}
