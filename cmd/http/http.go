package http

import (
	"log"
	"net/http"
)

type PathFunction struct {
	path     string
	function func(http.ResponseWriter, *http.Request)
}

func startServer(address string, pathFunctions []PathFunction) {
	for _, v := range pathFunctions {
		http.HandleFunc(v.path, v.function)
	}

	err := http.ListenAndServe(address, nil)

	if err != nil {
		log.Fatal(err.Error())
	}
}
