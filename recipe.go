package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lib/pq"

	"github.com/gorilla/mux"
)

// Recipe is our json format for preppy recipes
type Recipe struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	URL          string   `json:"url"`
	Src          string   `json:"src"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Img          string   `json:"img"`
	Video        string   `json:"video"`
	Category     string   `json:"category"`
}

// GetRecipes will query 20 rows of recipes and encode them to json
func GetRecipes(w http.ResponseWriter, r *http.Request) {
	// default page offset
	pageOffset := 0

	// query parameters
	page := r.URL.Query().Get("page")
	// check if page query exists
	if page != "" {
		// offset rows result by 20 * pagenumber
		pagenum, _ := strconv.Atoi(page)
		pageOffset = 20 * pagenum
	}

	// begin query for recipes
	rows, _ := Db.Query(
		`
		SELECT pid, name, url, src, ingredients, instructions, img, video, category 
		FROM recipe LIMIT 20 OFFSET $1;
		`, pageOffset)

	// add query rows to results slice
	var results []Recipe
	for rows.Next() {
		res := Recipe{}
		rows.Scan(&res.ID, &res.Name, &res.URL,
			&res.Src, pq.Array(&res.Ingredients), pq.Array(&res.Instructions),
			&res.Img, &res.Video, &res.Category)

		results = append(results, res)
	}

	// encode results to json
	json.NewEncoder(w).Encode(results)
}

// GetRecipeById returns a recipe with id given
func GetRecipeByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	row := Db.QueryRow(
		`
		SELECT pid, name, url, src, ingredients, instructions, img, video, category FROM recipe WHERE pid = $1;
		`, params["id"])

	res := Recipe{}
	row.Scan(&res.ID, &res.Name, &res.URL,
		&res.Src, pq.Array(&res.Ingredients), pq.Array(&res.Instructions),
		&res.Img, &res.Video, &res.Category)

	json.NewEncoder(w).Encode(res)
}
