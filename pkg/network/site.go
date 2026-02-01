package network

type ListLocalSitesRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Filter string `json:"filter"`
}

type ListLocalSitesResponse struct {
	Offset     int            `json:"offset"`
	Limit      int            `json:"limit"`
	Count      int            `json:"count"`
	TotalCount int            `json:"totalCount"`
	Data       []SiteOverview `json:"data"`
}

type SiteOverview struct {
	ID                string `json:"id"`
	InternalReference string `json:"internalReference"`
	Name              string `json:"name"`
}

// GET /v1/sites
