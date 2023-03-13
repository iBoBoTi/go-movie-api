package domain

type MovieCharacter struct {
	Name      string   `json:"name"`
	Mass      string   `json:"mass"`
	Height    string   `json:"height"`
	HairColor string   `json:"hair_color"`
	SkinColor string   `json:"skin_color"`
	EyeColor  string   `json:"eye_color"`
	BirthYear string   `json:"birth_year"`
	Gender    string   `json:"gender"`
	HomeWorld string   `json:"homeworld"`
	Films     []string `json:"films"`
	Vehicles  []string `json:"vehicles"`
	Starships []string `json:"starships"`
	Created   string   `json:"created"`
	Edited    string   `json:"edited"`
	Url       string   `json:"url"`
}

type CharacterListResponse struct {
	Characters                    []MovieCharacter `json:"characters"`
	CharactersCount               int              `json:"characters_count"`
	TotalHeightOfCharactersInCM   string           `json:"total_height_of_characters_in_cm"`
	TotalHeightOfCharactersInFeet string           `json:"total_height_of_characters_in_feet"`
}
