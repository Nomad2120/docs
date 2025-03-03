package model

// SignWSSERequest -
type SignWSSERequest struct {
	SignNodeID string `json:"signNodeId" form:"signNodeId" binding:"required"`
	Alias      string `json:"alias" form:"alias"`
	Data       string
}
