package model

import "time"

const AppName = "oci-docs"

// BaseResponse -
type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse -
type ErrorResponse struct {
	BaseResponse
	Result string `json:"result"`
}

// DocRequest -
type DocRequest struct {
	ID int `json:"id" form:"id"`
}

// DocResponse -
type DocResponse struct {
	BaseResponse
	Result string `json:"result"`
}

// FillContractRequest -
type FillContractRequest struct {
	ID          int     `json:"id"`
	IDN         string  `json:"idn"`
	Name        string  `json:"name"`
	FIO         string  `json:"fio"`
	Phone       string  `json:"phone"`
	Address     string  `json:"address"`
	ApartCount  int     `json:"apartCount"`
	Email       string  `json:"email"`
	CreateDt    string  `json:"createDt"`
	Site        string  `json:"site"`
	Tariff      float64 `json:"tariff"`
	UnionTypeRu string  `json:"unionTypeRu"`
	UnionTypeKz string  `json:"unionTypeKz"`
}

func (r *FillContractRequest) SignDate() string {
	return time.Now().Format("02/01/2006")
}

// FillContractResponse -
type FillContractResponse struct {
	BaseResponse
	Result *FillContractResult `json:"result"`
}

// FillContractResult -
type FillContractResult struct {
	HTMLBase64 string `json:"htmlBase64"`
	PDFBase64  string `json:"pdfBase64"`
}

// SignContractRequest -
type SignContractRequest struct {
	ID        int    `json:"id" uri:"id" binding:"required"`
	DocBase64 string `json:"docBase64"`
	Extension string `json:"extension"`
}

// SignContractResponse -
type SignContractResponse struct {
	BaseResponse
	Result string `json:"result"`
}

// RegistrationResponse -
type RegistrationResponse struct {
	BaseResponse
	Result RegistrationInfo `json:"result"`
}

// RegistrationInfo -
type RegistrationInfo struct {
	ID                int    `json:"id"`
	CreateDate        string `json:"createDt"`
	Name              string `json:"name"`
	IDN               string `json:"idn"`
	Address           string `json:"address"`
	UserID            int    `json:"userId"`
	Phone             string `json:"phone"`
	Email             string `json:"email"`
	ApartCount        int    `json:"apartCount"`
	StateCode         string `json:"stateCode"`
	AddressKZ         string `json:"addressKz"`
	AddressRegistryId int    `json:"addressRegistryId"`
	AtsId             int    `json:"atsId"`
	UnionTypeRu       string `json:"unionTypeRu"`
	UnionTypeKz       string `json:"unionTypeKz"`
}

// SignActRequest -
type SignActRequest struct {
	ID        int    `json:"id" uri:"id" binding:"required"`
	DocBase64 string `json:"docBase64"`
	Extension string `json:"extension"`
}

// SignActResponse -
type SignActResponse struct {
	BaseResponse
	Result string `json:"result"`
}

type GetServiceGroupsResponse struct {
	BaseResponse
	Result []ServiceGroup `json:"result"`
}

type ServiceGroup struct {
	ID                int    `json:"id"`
	NameRu            string `json:"nameRu"`
	NameKz            string `json:"nameKz"`
	AccountTypeCode   string `json:"accountTypeCode"`
	AccountTypeNameRu string `json:"accountTypeNameRu"`
	AccountTypeNameKz string `json:"accountTypeNameKz"`
}

// AddDocRequest -
type AddDocRequest struct {
	DocTypeCode string `json:"docTypeCode" binding:"required"`
	Data        string `json:"data" binding:"required"`
	Extension   string `json:"extension" binding:"required"`
}

type AddDocResponse struct {
	BaseResponse
	Result AddDocResult `json:"result"`
}

type AddDocResult struct {
	ID            int    `json:"id"`
	DocTypeCode   string `json:"docTypeCode"`
	DocTypeNameRu string `json:"docTypeNameRu"`
	DocTypeNameKz string `json:"docTypeNameKz"`
	Scan          Scan   `json:"scan"`
}

type Scan struct {
	ID       int    `json:"id"`
	FileName string `json:"fileName"`
}

// SignOsiContractRequest -
type SignOsiContractRequest struct {
	Path struct {
		ID int `json:"id" uri:"id" binding:"required"`
	}
	Body struct {
		DocBase64 string `json:"docBase64" binding:"required"`
		Extension string `json:"extension" binding:"required"`
	}
}
