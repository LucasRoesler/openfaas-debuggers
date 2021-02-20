package function

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	nats "github.com/nats-io/nats.go"
)

var (
	subject        = "nats-test"
	defaultMessage = "Hello World"
)

// Handle a serverless request
func Handle(w http.ResponseWriter, r *http.Request) {

	msg := defaultMessage
	if r.Body != nil {
		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("error reading body: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		msg = string(bytes.TrimSpace(body))
	}

	natsURL := nats.DefaultURL
	val, ok := os.LookupEnv("nats_url")
	if ok {
		natsURL = val
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		errMsg := fmt.Sprintf("can not connect to nats: %s", err)
		log.Print(errMsg)

		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(errMsg))
		if err != nil {
			log.Printf("error writing output: %s", err)
		}
		return
	}
	defer nc.Close()

	s := r.Header.Get("nats-subject")
	if s != "" {
		subject = s
	}

	log.Printf("Publishing %d bytes to: %q\n", len(msg), subject)

	err = nc.Publish(subject, []byte(msg))
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(fmt.Sprintf("can not publish to nats: %s", err)))
		if err != nil {
			log.Printf("error writing output: %s", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("Published %d bytes to: %q", len(msg), subject)))
	if err != nil {
		log.Printf("error writing output: %s", err)
	}
}
