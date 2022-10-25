package mangodex

// AuthorAttributes : Attributes for an Author.
type AuthorAttributes struct {
	Name      string           `json:"name"`
	ImageURL  string           `json:"imageUrl"`
	Biography LocalisedStrings `json:"biography"`
	Version   int              `json:"version"`
	CreatedAt string           `json:"createdAt"`
	UpdatedAt string           `json:"updatedAt"`
}
