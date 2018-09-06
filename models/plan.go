package models

import (
	"database/sql"
)

// Plan is our json format for preppy plans
type Plan struct {
	ID    int            `json:"id"`
	User  int            `json:"user"`
	Title string         `json:"name"`
	Desc  sql.NullString `json:"desc,omitempty"`
	Likes int            `json:"likes"`
}

// PlanRecipe are the recipes within a plan
type PlanRecipe struct {
	ID       int    `json:"id"`
	MealType string `json:"type"`
}
