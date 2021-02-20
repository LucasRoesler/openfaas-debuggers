package function

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Handle a function invocation
func Handle(w http.ResponseWriter, r *http.Request) {
	var err error
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		input, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("error reading body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	log.Printf("Received: %q", string(input))

	if val, ok := os.LookupEnv("wait"); ok && len(val) > 0 {
		parsedVal, _ := time.ParseDuration(val)
		log.Printf("Waiting for %s before returning", parsedVal.String())
		time.Sleep(parsedVal)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("Received: %q", string(input))))
	if err != nil {
		log.Printf("error writing output: %s", err)
	}
}
