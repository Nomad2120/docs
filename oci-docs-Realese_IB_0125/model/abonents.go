package model

type OSIAbonentsResponse struct {
	BaseResponse
	Result []OSIAbonent `json:"result"`
}

type OSIAbonent struct {
	AreaTypeCode    string      `json:"areaTypeCode"`
	AreaTypeNameKz  string      `json:"areaTypeNameKz"`
	AreaTypeNameRu  string      `json:"areaTypeNameRu"`
	ErcAccount      interface{} `json:"ercAccount"`
	Flat            string      `json:"flat"`
	Floor           int64       `json:"floor"`
	ID              int64       `json:"id"`
	IDN             string      `json:"idn"`
	LivingFact      int64       `json:"livingFact"`
	LivingJur       int64       `json:"livingJur"`
	Name            string      `json:"name"`
	OsiID           int64       `json:"osiId"`
	Owner           string      `json:"owner"`
	Phone           string      `json:"phone"`
	Square          float64     `json:"square"`
	EffectiveSquare *float64    `json:"effectiveSquare"`
}

// OSIAbonentsRequest -
type OSIAbonentsRequest struct {
	ID int `json:"id"`
}

type AbonentsInfoResponse struct {
	BaseResponse
	Result OSIAbonent `json:"result"`
}
