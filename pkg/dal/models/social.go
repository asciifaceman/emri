package models

/*
Captures the API response from the mastodon API endpoint

/api/v2/instance
*/
type AboutResponse struct {
	Domain      string           `json:"domain"`
	Title       string           `json:"title"`
	Version     string           `json:"version"`
	SourceUrl   string           `json:"source_url"`
	Description string           `json:"description"`
	Contact     *InstanceContact `json:"contact"`
	Rules       []*InstanceRule  `json:"rules"`
	Thumbnail   *InstanceImage   `json:"thumbnail"`
	Extra       interface{}      `json:"-"`
}

type InstanceRule struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Hint string `json:"hint"`
}

type InstanceContact struct {
	Email string      `json:"email"`
	Extra interface{} `json:"-"`
}

type InstanceImage struct {
	URL   string      `json:"url"`
	Extra interface{} `json:"-"`
}

type DomainBlock struct {
	Domain   string `json:"domain"`
	Digest   string `json:"digest"`
	Severity string `json:"severity"`
	Comment  string `json:"comment"`
}

// DomainBlockResponse encodes responses from /api/v1/instance/domain_blocks
type DomainBlockResponse struct {
	Blocks []*DomainBlock `json:"-"`
}

type PeeredResponse struct {
	Peers []string `json:"-"`
}

type ActivityResponse struct {
	Activity []*SingleActivityResponse `json:"-"`
}

type SingleActivityResponse struct {
	Week          string `json:"week"`
	Statuses      string `json:"statuses"`
	Logins        string `json:"logins"`
	Registrations string `json:"registrations"`
}
