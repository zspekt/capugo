package json

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// writes to $w a response with http status code $code and the json $payload
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	log.Printf("Responding with code %v and provided payload...\n", code)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		return
	}

	w.Write(jsonPayload)
}

// decodes reader with JSON-formatted content $r into the STRUCT pointed at by $st
func DecodeJson[T any](r io.Reader, st *T) error {
	decoder := json.NewDecoder(r)
	err := decoder.Decode(st)
	if err != nil {
		log.Printf("\nError decoding json %v\n", err)
		return err
	}
	return nil
}
