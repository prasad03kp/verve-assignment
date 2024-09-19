package verve

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var (
	Result string = "ok"
	Client *http.Client = &http.Client{}
)

func Accept(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		Result = "failed"
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		endpoint := r.URL.Query().Get("endpoint")
		if endpoint != "" {
			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				log.Printf("Error occured. Making new request to %s failed\n", endpoint)
			}
			req.URL.Query().Add("count", strconv.Itoa(id))

			res, err := Client.Do(req)
			if err != nil {
				log.Printf("Error occured. API call to %s failed\n", endpoint)
			} else {
				log.Printf("Status code for %s : %d", endpoint, res.StatusCode)
			}
		}
		Result = "ok"
	}

	w.Write([]byte(Result))
}

func Endpoint(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		Result = "Error: Count is not integer"
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(Result))
	} else {
		Result = fmt.Sprintf("Count is %d", count)
		w.Write([]byte(Result))
	}
}