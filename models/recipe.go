package models

// Recipe is our json format for preppy recipes
type Recipe struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Servings     int      `json:"servings"`
	URL          string   `json:"url,omitempty"`
	Src          string   `json:"src,omitempty"`
	Ingredients  []string `json:"ingredients,omitempty"`
	Instructions []string `json:"instructions,omitempty"`
	Img          string   `json:"img"`
	Video        string   `json:"video,omitempty"`
	Category     string   `json:"category,omitempty"`
}
