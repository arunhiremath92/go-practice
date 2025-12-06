package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func handlerHttpRequestVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handling http version request call")
	versionResponse := map[string]string{
		"version": "0.1",
		"time":    time.Now().UTC().String(),
	}
	jsonData, error := json.Marshal(versionResponse)
	if error != nil {
		errorString := fmt.Errorf("failed to marshal the object %s", error)
		fmt.Println(errorString)
		w.WriteHeader(http.StatusBadRequest)
		return

	}
	w.Write(jsonData)
	w.Write([]byte("\n"))
}

func main() {
	fmt.Println("starting the module")
	http.HandleFunc("/version", handlerHttpRequestVersion)
	port := ":6000"
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("faild to start the server %s", err)
	}
}
