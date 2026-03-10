package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/validator"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(r *http.Request, idKey ...string) (int, error) {
	params := httprouter.ParamsFromContext(r.Context())

	key := "id" // default
	if len(idKey) > 0 && idKey[0] != "" {
		key = idKey[0]
	}

	id, err := strconv.Atoi(params.ByName(key))
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)

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
	// limit body to 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	// initialize the decoder and call DisallowUnknownFields() method on it
	// if field doesn't map to a struct it will return an error
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var MaxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		// using .Is here to match a specific error not matching a type
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshTypeError):
			if unmarshTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains an unknown key %s", fieldName)
		case errors.As(err, &MaxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", MaxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// helper that returns string value from the query string or the provided default value
func (app *application) readString(qs url.Values, key, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	return s
}

// reads a string from query string and split it into a slice on comma character if no value return default
func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

// reads query string and converts it to an int, if key not found return default
// if value couldnt be converted to int record an error msg in provided Validator instance
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	// try to convert to int, if fails add error message to validator instance
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
}
