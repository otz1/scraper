package entity

type Result struct {
	Title   string `json:"title"`
	Href    string `json:"href"`
	Snippet string `json:"snippet"`
}

type ScrapeResponse struct {
	OriginalQuery string   `json:"originalQuery"`
	Results       []Result `json:"results"`
}
