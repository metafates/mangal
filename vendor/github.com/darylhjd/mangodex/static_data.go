package mangodex

// Publication demographic
const (
	ShonenDemographic string = "shounen"
	ShoujoDemographic string = "shoujo"
	JoseiDemographic  string = "josei"
	SeinenDemograpic  string = "seinen"
)

// Manga publication status
const (
	OngoingStatus   string = "ongoing"
	CompletedStatus string = "completed"
	HiatusStatus    string = "hiatus"
	CancelledStatus string = "cancelled"
)

// Manga reading status
const (
	Reading    string = "reading"
	OnHold     string = "on_hold"
	PlanToRead string = "plan_to_read"
	Dropped    string = "dropped"
	ReReading  string = "re_reading"
	Completed  string = "completed"
)

// Manga content rating
const (
	Safe       string = "safe"
	Suggestive string = "suggestive"
	Erotica    string = "erotica"
	Porn       string = "pornographic"
)

// Relationship types. Useful for reference expansions
const (
	MangaRel           string = "manga"
	ChapterRel         string = "chapter"
	CoverArtRel        string = "cover_art"
	AuthorRel          string = "author"
	ArtistRel          string = "artist"
	ScanlationGroupRel string = "scanlation_group"
	TagRel             string = "tag"
	UserRel            string = "user"
	CustomListRel      string = "custom_list"
)
