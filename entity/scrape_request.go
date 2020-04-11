package entity

// ScrapeSource is a source to scrape from
// e.g. wikipedia.
type ScrapeSource string

// The sources available to scrape
// that are currently supported
const (
	DDG       ScrapeSource = "DDG"
	STARTPAGE              = "STARTPAGE"
	YAHOO                  = "YAHOO"
	BING                   = "BING"
	WIKIPEDIA              = "WIKIPEDIA"
)

// ScrapeRequest ...
type ScrapeRequest struct {
	Query string `json:"query"`
}
