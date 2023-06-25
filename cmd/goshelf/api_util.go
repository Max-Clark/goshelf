package goshelf

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func getPathFunctions(cfg *GoshelfConfig) []PathFunction {
	return []PathFunction{
		{
			Path: BookPath,
			Function: func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					ApiBookFilter(cfg, w, r)
				case http.MethodPost:
					// ApiBookCreate(cfg, w, r)
				default:
					errMsg := "unsupported HTTP method"
					returnGoshelfErrorWithMessage(&errMsg, w, r)
				}
			},
		},
		{
			Path: BookPath + "{id:[0-9]+}",
			Function: func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					ApiBookGet(cfg, w, r)
				case http.MethodDelete:
					ApiBookRemove(cfg, w, r)
				default:
					errMsg := "unsupported HTTP method"
					returnGoshelfErrorWithMessage(&errMsg, w, r)
				}
			},
		},
		{
			Path: CollectionPath,
			Function: func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPost:
					ApiCollectionCreate(cfg, w, r)
				default:
					errMsg := "unsupported HTTP method"
					returnGoshelfErrorWithMessage(&errMsg, w, r)
				}
			},
		},
		{
			Path: CollectionPath + "{title:[a-zA-Z0-9_-]+}",
			Function: func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					ApiCollectionGet(cfg, w, r)
				case http.MethodDelete:
					ApiCollectionDelete(cfg, w, r)
				default:
					errMsg := "unsupported HTTP method"
					returnGoshelfErrorWithMessage(&errMsg, w, r)
				}
			},
		},
	}
}

func checkContentType(contentTypeExpected string, r *http.Request) error {
	val, ok := r.Header["Content-Type"]
	if !ok {
		return errors.New("content-type header missing")
	}

	lowerCT := strings.ToLower(contentTypeExpected)
	lowerVal := strings.ToLower(val[0])

	if !strings.Contains(lowerVal, lowerCT) {
		return errors.New("expected content-type to contain " + contentTypeExpected)
	}

	return nil
}

func createResponseObject(status string, code int, metadata *map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type":        "sync",
		"status":      status, // "Success", "Error"
		"status_code": code,   // e.g., 400
		"metadata":    metadata,
	}
}

func returnGoshelfErrorWithMessage(msg *string, w http.ResponseWriter, r *http.Request) {
	metadata := map[string]interface{}{
		"message": *msg,
	}

	returnGoshelfResponse(StatusFailure, CodeFailure, &metadata, w, r)
}

func returnGoshelfSuccessWithNoObject(w http.ResponseWriter, r *http.Request) {
	metadata := map[string]interface{}{}
	returnGoshelfResponse(StatusSuccess, CodeSuccess, &metadata, w, r)
}

func returnGoshelfSuccessWithObject(metadata *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	returnGoshelfResponse(StatusSuccess, CodeSuccess, metadata, w, r)
}

func returnGoshelfResponse(status string, code int, metadata *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(createResponseObject(status, code, metadata))
}
