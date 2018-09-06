package main

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/itstc/preppyAPI/models"
)

// GetPlans retrieves all meal plans from database
func (a *App) GetPlans(w http.ResponseWriter, r *http.Request) {

	// query all plans by likes
	rows, err := a.Db.Query("SELECT * FROM plan ORDER BY likes DESC")
	CheckErr(err)

	planList := []models.Plan{}

	// append plan to slice
	for rows.Next() {
		var plan models.Plan
		err = rows.Scan(&plan.ID, &plan.User, &plan.Title, &plan.Desc, &plan.Likes)
		if err != nil {
			return
		}
		planList = append(planList, plan)
	}

	// convert slice of plans to JSON and send to client
	WriteJSON(w, planList)
}

// GetPlanByID retrieves recipes within the plan and responds with JSON of meal plan
func (a *App) GetPlanByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// meal plan data
	var planID int
	var planTitle string
	var planDesc sql.NullString

	// query to retrieve meal plan
	rows, err := a.Db.Query(`
	SELECT id, title, description, recipeId, type 
	FROM plan a INNER JOIN plan_recipe b ON a.id = b.planId
	WHERE a.id = $1;
	`, params["id"])

	// error handling
	CheckErr(err)

	// create an array of the recipes
	recipes := []models.PlanRecipe{}
	for rows.Next() {
		var recipe models.PlanRecipe
		rows.Scan(&planID, &planTitle, &planDesc, &recipe.ID, &recipe.MealType)
		recipes = append(recipes, recipe)
	}

	// write response of the meal plan in JSON
	WriteJSON(w, map[string]interface{}{
		"id":      planID,
		"desc":    planDesc,
		"title":   planTitle,
		"recipes": recipes,
	})

}

// GetPlansByUser retrieves all plans created by user
func (a *App) GetPlansByUser(w http.ResponseWriter, r *http.Request) {

}
