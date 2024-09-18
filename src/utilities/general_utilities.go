package utilities

import (
	"net/http"
	"os"
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
