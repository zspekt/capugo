package json

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

// writes to $w a response with http status code $code and the json $payload
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	slog.Info("responding with json", "code", code)
	w.Header().Set("Content-Type", "application/json")

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		slog.Error("error marshalling json payload", "error", err)
		return
	}
	// should return status code after checking for an error
	w.WriteHeader(code)
	w.Write(jsonPayload)
}

// decodes reader with JSON-formatted content $r into the STRUCT pointed at by $st
func DecodeJson[T any](r io.Reader, st *T) error {
	decoder := json.NewDecoder(r)

	err := decoder.Decode(st)
	if err != nil {
		slog.Error("error decoding json", "error", err)
		return err
	}
	return nil
}
