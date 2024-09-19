package verve

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	memcached "github.com/prasad03kp/verve-assignment/memcached"
	"github.com/prasad03kp/verve-assignment/utilities"
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
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			memcached.WriteToMemCache(id)
		}()
		endpoint := r.URL.Query().Get("endpoint")
		if endpoint != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				utilities.MakeGetCall(endpoint, id, Client)
			}()
		}
		Result = "ok"
		wg.Wait()
	}
	w.Write([]byte(Result))
}

func Endpoint(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		log.Println("Error: Count is not integer")
	} else {
		log.Printf("Count is %d", count)
	}
}