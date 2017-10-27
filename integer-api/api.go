package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/guilherme-santos/dbmapping"
	"github.com/guilherme-santos/dbmapping/mysql"
	"github.com/julienschmidt/httprouter"
)

type (
	IntegerAPI struct {
		dbClient *sql.DB
		router   *httprouter.Router
		table    dbmapping.Table
	}

	Entity struct {
		ID        int64     `json:"id"`
		Number    int       `json:"integer"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func NewIntegerAPI(dbClient *sql.DB, router *httprouter.Router) *IntegerAPI {
	api := &IntegerAPI{
		dbClient: dbClient,
		router:   router,
	}

	router.GET("/", api.index)
	router.GET("/v1", api.index)
	router.GET("/v1/integers", api.getIntegers)
	router.GET("/v1/integers/last", api.getLastInteger)
	router.POST("/v1/integers", api.newInteger)

	mapping := mysql.NewMapping(dbClient)
	api.table, _ = mapping.NewTable("integers", []dbmapping.Field{
		{Name: "id", Type: "INT AUTO_INCREMENT", Null: false, PrimaryKey: true},
		{Name: "number", Type: "INT UNSIGNED", Null: false},
		{Name: "created_at", Type: "TIMESTAMP", Default: "CURRENT_TIMESTAMP"},
	})

	err := api.table.CreateTable()
	if err != nil {
		log.Fatalf("Cannot create mapping to table: %s", err)
	}

	return api
}

func (api *IntegerAPI) Run(port string) error {
	return http.ListenAndServe(":"+port, api.router)
}

func (api *IntegerAPI) index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	responseJSON(w, map[string]interface{}{
		"endpoints": map[string]interface{}{
			"root": map[string]string{
				"method": "GET",
				"path":   "/v1",
			},
			"get_integers": map[string]string{
				"method": "GET",
				"path":   "/v1/integers",
			},
			"get_last_integer": map[string]string{
				"method": "GET",
				"path":   "/v1/integers/last",
			},
			"new_integer": map[string]string{
				"method": "POST",
				"path":   "/v1/integers",
			},
		},
	})
}

func (api *IntegerAPI) getIntegers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sort := []dbmapping.OrderByClause{
		{Field: "id", Type: dbmapping.OrderByDESC},
	}

	results, err := api.table.Query(nil, nil, sort, []int{0, 10})
	if err != nil {
		responseJSONWithStatus(w, http.StatusBadRequest, Error{
			Code:    "database_failed",
			Message: fmt.Sprintf("Cannot read data from database: %s", err),
		})
		return
	}

	entities := make([]Entity, 0, len(results))
	for _, result := range results {
		entity := Entity{
			ID:        result["id"].(int64),
			Number:    int(result["number"].(int64)),
			CreatedAt: result["created_at"].(time.Time),
		}

		entities = append(entities, entity)
	}

	responseJSONWithStatus(w, http.StatusCreated, entities)
}

func (api *IntegerAPI) getLastInteger(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sort := []dbmapping.OrderByClause{
		{Field: "id", Type: dbmapping.OrderByDESC},
	}

	result, err := api.table.QueryOne(nil, nil, sort)
	if err != nil {
		responseJSONWithStatus(w, http.StatusBadRequest, Error{
			Code:    "database_failed",
			Message: fmt.Sprintf("Cannot read data from database: %s", err),
		})
		return
	}

	if result == nil {
		responseJSONWithStatus(w, http.StatusBadRequest, Error{
			Code:    "empty_database",
			Message: "No data found into database",
		})
		return
	}

	entity := Entity{
		ID:        result["id"].(int64),
		Number:    int(result["number"].(int64)),
		CreatedAt: result["created_at"].(time.Time),
	}

	responseJSONWithStatus(w, http.StatusCreated, entity)
}

func (api *IntegerAPI) newInteger(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var body Entity

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		responseJSONWithStatus(w, http.StatusBadRequest, Error{
			Code:    "invalid_body",
			Message: fmt.Sprintf("Cannot read body: %s", err),
		})
		return
	}

	defer r.Body.Close()

	if body.Number == 0 {
		responseJSONWithStatus(w, http.StatusBadRequest, Error{
			Code:    "invalid_integer",
			Message: fmt.Sprintf("Integer \"%v\" is not valid", body.Number),
		})
		return
	}

	err = api.table.Insert(map[string]interface{}{
		"id":     nil,
		"number": body.Number,
	})
	if err != nil {
		log.Println("Cannot insert integer:", err)
		responseJSONWithStatus(w, http.StatusInternalServerError, Error{
			Code:    "database_error",
			Message: err.Error(),
		})
		return
	}

	responseJSONWithStatus(w, http.StatusCreated, nil)
}
