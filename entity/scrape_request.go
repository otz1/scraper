package entity

// ScrapeSource is a source to scrape from
// e.g. wikipedia.
type ScrapeSource string

// The sources available to scrape
// that are currently supported
const (
	GOOGLE    ScrapeSource = "GOOGLE"
	DDG       ScrapeSource = "DDG"
	WIKIPEDIA ScrapeSource = "WIKIPEDIA"
)

// ScrapeRequest ...
type ScrapeRequest struct {
	Query  string        `json:"query"`
	Source *ScrapeSource `json:"source"`
}
