package verve

import (
	"net/http"
	"strconv"
	"fmt"
)

var (
	Result string = "ok"
)

func Accept(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		Result = "failed"
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		endpoint := r.URL.Query().Get("endpoint")
		if endpoint != "" {
			fmt.Println("making api call", id)
		}
		Result = "ok"
	}

	w.Write([]byte(Result))
}
