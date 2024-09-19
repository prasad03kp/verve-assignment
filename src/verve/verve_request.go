package verve

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"io"
	"encoding/json"
	memcached "github.com/prasad03kp/verve-assignment/memcached"
	"github.com/prasad03kp/verve-assignment/utilities"
)

var (
	Result string = "ok"
	Client *http.Client = &http.Client{}
)

type Body struct {
	FreeText string
}

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
				uniqueCount := memcached.CountUniqueIDsInCurrentMinute()
				var freeText string = "This is test string"
				utilities.MakePostCall(endpoint, freeText, uniqueCount, Client)
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

func PostEndpoint(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		log.Println("Error: Count is not integer")
	} else {
		log.Printf("Count is %d", count)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading endpoint body. %v", err)
	}

	var bodyContent Body 
	err = json.Unmarshal(body, &bodyContent)
	if err != nil {
		log.Printf("Error parsing endpoint body. %v", err)
	}

	fmt.Printf("Free text: %s \n", bodyContent.FreeText)
}