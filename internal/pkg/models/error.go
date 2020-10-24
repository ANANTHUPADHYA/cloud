package models

//ErrorResponse model for returning errors
type ErrorResponse struct {
	ErrorCode            string   `json:"errorCode"`
	ErrorStatusCode      int      `json:"-"`
	Message              string   `json:"message"`
	Details              string   `json:"details"`
	RecommendationAction []string `json:"recommendedActions"`
	NestedErrors         string   `json:"nestedErrors"`
	ErrorSource          string   `json:"errorSource"`
	Data                 string   `json:"data"`
	CanForce             bool     `json:"canForce"`
}
