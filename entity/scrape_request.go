package entity

// ScrapeSource is a source to scrape from
// e.g. wikipedia.
type ScrapeSource string

// The sources available to scrape
// that are currently supported
const (
	GOOGLE    ScrapeSource = "GOOGLE"
	WIKIPEDIA              = "WIKIPEDIA"
)

// ScrapeRequest ...
type ScrapeRequest struct {
	Query  string        `json:"query"`
	Source *ScrapeSource `json:"source"`
}
