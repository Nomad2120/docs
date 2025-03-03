package model

type BuildingInfoResponse struct {
	ID           int    `json:"id"`
	Rca          string `json:"rca"`
	Number       string `json:"number"`
	FullPathRus  string `json:"fullPathRus"`
	FullPathKaz  string `json:"fullPathKaz"`
	ShortPathRus string `json:"shortPathRus"`
	ShortPathKaz string `json:"shortPathKaz"`
	Geonim       struct {
		ID           int    `json:"id"`
		NameRus      string `json:"nameRus"`
		NameKaz      string `json:"nameKaz"`
		FullPathRus  string `json:"fullPathRus"`
		FullPathKaz  string `json:"fullPathKaz"`
		ShortPathRus string `json:"shortPathRus"`
		ShortPathKaz string `json:"shortPathKaz"`
	} `json:"geonim"`
	Ats struct {
		ID           int    `json:"id"`
		NameRus      string `json:"nameRus"`
		NameKaz      string `json:"nameKaz"`
		FullPathRus  string `json:"fullPathRus"`
		FullPathKaz  string `json:"fullPathKaz"`
		ShortPathRus string `json:"shortPathRus"`
		ShortPathKaz string `json:"shortPathKaz"`
	} `json:"ats"`
}
