package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/itstc/preppyAPI/models"
	"github.com/lib/pq"

	"github.com/gorilla/mux"
)

var (
	JSONGetRecipeError = map[string]string{
		"error": "unable to retrieve recipes!",
	}
	JSONGetRecipeIDError = map[string]string{
		"error": "unable to retrieve recipe!",
	}
)

// GetRecipes will query 20 rows of recipes and encode them to json
func (a *App) GetRecipes(w http.ResponseWriter, r *http.Request) {

	// default page offset
	pageOffset := 0
	pageLimit := 20

	// query quantity parameters
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	// check if page query exists
	if page != "" {
		// offset rows result by 20 * pagenumber
		pagenum, _ := strconv.Atoi(page)
		pageOffset = 20 * pagenum
	}

	// set limit if exists
	if limit != "" {
		pageLimit, _ = strconv.Atoi(limit)

		// max return is length of 20 recipes
		if pageLimit > 20 {
			pageLimit = 20
		}
	}

	// query search parameters
	search := []string{}
	if id := r.URL.Query().Get("id"); id != "" {
		search = strings.Split(id, ",")
	}

	var rows *sql.Rows
	var err error

	// begin query
	if len(search) > 0 {
		// we are searching based on given ids
		rows, err = a.Db.Query(
			`
			SELECT id, name, img 
			FROM recipe WHERE pid = ANY($1) LIMIT $2 OFFSET $3;
			`, pq.Array(search), pageLimit, pageOffset)
	} else {
		// do regular search based on likes
		rows, err = a.Db.Query(
			`
			SELECT id, name, img 
			FROM recipe LIMIT $1 OFFSET $2;
			`, pageLimit, pageOffset)
	}

	// error has occurred during query
	if err != nil {
		WriteJSON(w, JSONGetRecipeError)
		return
	}

	// add query rows to results slice
	var results []models.Recipe
	for rows.Next() {
		res := models.Recipe{}
		rows.Scan(&res.ID, &res.Name, &res.Img)

		results = append(results, res)
	}

	// return json response of recipes
	WriteJSON(w, results)

}

// GetRecipeByID returns a recipe with id given
func (a *App) GetRecipeByID(w http.ResponseWriter, r *http.Request) {
	// retrieve uri parameters
	params := mux.Vars(r)

	rows, err := a.Db.Query(
		`
		SELECT id, name, servings, url, src, ingredients, instructions, img, video, category FROM recipe WHERE id = $1;
		`, params["id"])

	// error or no results
	if err != nil || !rows.Next() {
		w.WriteHeader(404)
		WriteJSON(w, JSONGetRecipeIDError)
		return
	}

	res := models.Recipe{}
	rows.Scan(&res.ID, &res.Name, &res.Servings, &res.URL,
		&res.Src, pq.Array(&res.Ingredients), pq.Array(&res.Instructions),
		&res.Img, &res.Video, &res.Category)

	// encode result to json
	WriteJSON(w, res)
}
