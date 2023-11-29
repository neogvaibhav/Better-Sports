package api

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

func CreateNewResponse(w http.ResponseWriter, statusCode int, response *Response) error {
	//Set the response code on the api
	w.WriteHeader(statusCode)
	//Convert the struct to a JSON body
	err := json.NewEncoder(w).Encode(response)
	return err
}
