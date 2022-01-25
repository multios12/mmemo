package modules

type Memo struct {
	Id   string `json:"id"`                    // id
	Name string `json:"name" validate:"min=1"` //name
	Date string `json:"date" validate:"len=10"`
	Shop string `json:"shop"`
	Page string `json:"page"`
	Play string `json:"play"`
	Talk string `json:"talk"`
}
