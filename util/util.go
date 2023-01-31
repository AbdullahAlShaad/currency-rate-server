package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

const float64EqualityThreshold = 1e-9

func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func WriteJSONResponse(w http.ResponseWriter, code int, data interface{}) {

	setJSONHeader(w)
	w.WriteHeader(code)
	switch x := data.(type) {
	case string:
		w.Write([]byte(x))
	case []byte:
		w.Write(x)
	default:
		err := json.NewEncoder(w).Encode(x)
		if err != nil {
			data = map[string]interface{}{
				"message": err.Error(),
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
func JsonError(w http.ResponseWriter, status int,
	error string, data map[string]interface{}) {
	response := struct {
		Status int                    `json:"status"`
		Error  string                 `json:"error"`
		Data   map[string]interface{} `json:"data,omitempty"`
	}{
		Status: status,
		Error:  error,
		Data:   data,
	}
	WriteJSONResponse(w, status, response)
}

func CreateMySqlDatabase(dbAddress string, password string, databaseName string) {
	query := fmt.Sprintf("root:%s@tcp(%s)/", password, dbAddress)
	db, err := sql.Open("mysql", query)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE  IF NOT EXISTS " + databaseName)
	if err != nil {
		panic(err)
	}
}
