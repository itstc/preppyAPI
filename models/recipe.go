package models

// Recipe is our json format for preppy recipes
type Recipe struct {
	ID           int      `json:"id,omitempty"`
	Name         string   `json:"name,omitempty"`
	Servings     int      `json:"servings,omitempty"`
	URL          string   `json:"url,omitempty"`
	Src          string   `json:"src,omitempty"`
	Ingredients  []string `json:"ingredients,omitempty"`
	Instructions []string `json:"instructions,omitempty"`
	Img          string   `json:"img,omitempty"`
	Video        string   `json:"video,omitempty"`
	Category     string   `json:"category,omitempty"`
}
