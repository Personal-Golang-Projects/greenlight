package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(r *http.Request) (int, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	opts := []json.Options{
		jsontext.Multiline(true),
	}

	js, err := json.Marshal(data, opts...)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// UnmarshalRead reads to EOF, so trailing data after the first value is rejected too
	err := json.UnmarshalRead(r.Body, dst, json.RejectUnknownMembers(true))
	if err != nil {
		var synErr *jsontext.SyntacticError
		var semErr *json.SemanticError
		var maxBytesErr *http.MaxBytesError

		switch {
		// body larger than 1MB
		case errors.As(err, &maxBytesErr):
			return &badRequestError{
				Message: fmt.Sprintf("body must not be larger than %d bytes", maxBytesErr.Limit),
			}

		// empty or cut-off body
		case errors.Is(err, io.ErrUnexpectedEOF):
			return &badRequestError{
				Message: "body must not be empty or contain incomplete JSON",
			}

		// malformed JSON: missing quote, trailing comma, extra value after the object, etc.
		case errors.As(err, &synErr):
			return &badRequestError{
				Message: fmt.Sprintf("body contains badly-formed JSON (at character %d)", synErr.ByteOffset),
			}

		// field the struct doesn't have
		case errors.As(err, &semErr) && errors.Is(err, json.ErrUnknownName):
			return &badRequestError{
				Message: fmt.Sprintf("body contains unknown field %q", semErr.JSONPointer.LastToken()),
			}

		// valid JSON but wrong type, e.g. "year": "2016" for an int field
		case errors.As(err, &semErr):
			if semErr.JSONPointer != "" {
				return &badRequestError{
					Message: fmt.Sprintf("body contains incorrect JSON type for field %q", semErr.JSONPointer.LastToken()),
				}
			}
			return &badRequestError{Message: "body contains incorrect JSON type"}

		// anything else (e.g. dst isn't a pointer) is a bug in our code → 500
		default:
			return err
		}
	}

	return nil
}
