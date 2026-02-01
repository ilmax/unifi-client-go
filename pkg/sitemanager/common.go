package sitemanager

// SuccessResponse contains common fields for successful API responses.
type SuccessResponse struct {
	HttpStatusCode int    `json:"httpStatusCode"`
	TraceID        string `json:"traceId"`
}
