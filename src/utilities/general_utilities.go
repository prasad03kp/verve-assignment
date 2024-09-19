package utilities

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	Version string
)

func GetVersion(w http.ResponseWriter, r *http.Request) {
	Version := os.Getenv("VERSION")
	if Version == "" {
		Version = "1.0.0"
	}
	
	w.Write([]byte(Version))
}

func MakeGetCall(endpoint string, id int, client *http.Client) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Printf("Error occured. Making new request to %s failed\n", endpoint)
	}
	query := req.URL.Query()
	query.Add("count", strconv.Itoa(id))
	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error occured. API call to %s failed\n", endpoint)
	} else {
		log.Printf("Status code for %s : %d", endpoint, res.StatusCode)
	}
}