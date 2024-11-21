package checkhealth

import (
	"encoding/json"
	"log"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	log.Println("OK!")
	host := r.Host
	log.Printf("Requested host %s", host)
	response, err := json.Marshal("I'm OK")
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
