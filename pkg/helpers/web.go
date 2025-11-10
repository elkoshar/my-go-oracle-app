package helpers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-chi/chi/v5"
	"oracle.com/oracle/my-go-oracle-app/pkg/validator"
)

func GetUrlPathString(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

func GetUrlPathInt(r *http.Request, key string) (int, error) {
	if !IsEmpty(chi.URLParam(r, key)) {
		i, err := strconv.Atoi(chi.URLParam(r, key))
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, nil
}

func GetUrlPathInt64(r *http.Request, key string) (int64, error) {
	if !IsEmpty(chi.URLParam(r, key)) {
		i, err := strconv.ParseInt(chi.URLParam(r, key), 10, 64)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, nil
}

func IsEmpty(object interface{}) bool {
	if object == nil || object == "" {
		return true
	}

	if reflect.ValueOf(object).Kind() == reflect.Slice ||
		reflect.ValueOf(object).Kind() == reflect.Array {
		arr := reflect.ValueOf(object)
		if arr.Len() == 0 || arr.IsZero() {
			return true
		}
	}

	if reflect.ValueOf(object).Kind() == reflect.Struct {
		empty := reflect.New(reflect.TypeOf(object)).Elem().Interface()
		if reflect.DeepEqual(object, empty) {
			return true
		}
	}

	return false
}

func ParseBodyAndValidate(r *http.Request, req interface{}) error {
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	_, err = validator.ValidateStruct(req)
	if err != nil {
		return err
	}

	return nil
}
