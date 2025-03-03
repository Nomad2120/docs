package model

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

type GetActResponse struct {
	BaseResponse
	Result GetActResult `json:"result"`
}

type GetActResult struct {
	ID                  int       `json:"id"`
	CreateDt            *ISODate  `json:"createDt"`
	SignDt              *ISODate  `json:"signDt"`
	ActPeriod           ISODate   `json:"actPeriod"`
	ActNum              string    `json:"actNum"`
	StateCode           string    `json:"stateCode"`
	StateName           string    `json:"stateName"`
	OsiID               int       `json:"osiId"`
	OsiName             string    `json:"osiName"`
	OsiIDN              string    `json:"osiIdn"`
	OsiAddress          string    `json:"osiAddress"`
	OsiPhone            string    `json:"osiPhone"`
	OsiRegistrationDate *ISODate  `json:"osiRegistrationDate"`
	ApartCount          int       `json:"apartCount"`
	PlanAccuralID       int       `json:"planAccuralId"`
	Amount              float64   `json:"amount"`
	Comission           float64   `json:"comission"`
	Debt                float64   `json:"debt"`
	Tariff              float64   `json:"tariff"`
	ActItems            []ActItem `json:"actItems"`
}

type ActItem struct {
	Description string   `json:"description"`
	DateWork    string   `json:"dateWork"`
	Quantity    *float64 `json:"quantity"`
	Price       *float64 `json:"price"`
	Amount      *float64 `json:"amount"`
	Note        *string  `json:"note"`
}

func (a *GetActResult) GetDateContract() string {
	if a.OsiRegistrationDate != nil {
		return a.OsiRegistrationDate.String()
	}
	return a.CreateDt.String()
}

func (a *GetActResult) GetOldDiscount() float64 {
	return -a.Amount / 2
}

func (a *GetActResult) GetOldTotal() float64 {
	return a.Amount // a.Amount / 2
}

func (a *GetActResult) GetDiscount() float64 {
	var result float64
	for _, item := range a.ActItems {
		if item.Price != nil {
			result += *item.Price
		}
	}
	return result
}

func (a *GetActResult) GetTotal() float64 {
	var result float64
	for _, item := range a.ActItems {
		if item.Amount != nil {
			result += *item.Amount
		}
	}
	return result
}

func (a *GetActResult) GetTotalQuantity() float64 {
	var result float64
	for _, item := range a.ActItems {
		if item.Quantity != nil {
			result += *item.Quantity
		}
	}
	return result
}

func (a *GetActResult) GetFullName() string {
	res := a.OsiName
	if a.OsiAddress != "" {
		res += ", " + a.OsiAddress
	}
	if a.OsiPhone != "" {
		res += ", тел. " + a.OsiPhone
	}
	return res
}

func (a *GetActResult) GetPeriod() string {
	return strings.ToLower(a.FormatMonth(time.Time(a.ActPeriod)))
}

func (a *GetActResult) FormatMonth(t time.Time) string {
	var months = [...]string{
		"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
	}
	return fmt.Sprintf("%s %d", months[t.Month()-1], t.Year())
}

// FillActRequest -
type FillActRequest struct {
	ID int `json:"id" uri:"id" binding:"required"`
}

// FillActResult -
type FillActResult struct {
	//HTMLBase64 string `json:"htmlBase64"`
	PDFBase64 string `json:"pdfBase64"`
}

// FillActResponse -
type FillActResponse struct {
	BaseResponse
	Result *FillActResult `json:"result"`
}

type ISODate time.Time

func (t *ISODate) UnmarshalText(text []byte) (err error) {
	if strings.Contains(string(text), `/`) {
		err = unmarshalTime(text, (*time.Time)(t), "02/01/2006")
		return
	}
	if len(text) == 19 {
		err = unmarshalTime(text, (*time.Time)(t), "2006-01-02T15:04:05")
	} else {
		err = unmarshalTime(text, (*time.Time)(t), "2006-01-02T15:04:05.000")
		if err != nil {
			err = unmarshalTime(text, (*time.Time)(t), "2006-01-02T15:04:05.999999")
		}
		if err != nil {
			err = unmarshalTime(text, (*time.Time)(t), "2006-01-02T15:04:05Z07:00") //2021-10-24T12:00:00+05:00
		}
	}
	return
}

func (t ISODate) MarshalText() ([]byte, error) {
	return []byte((time.Time)(t).Format("02/01/2006")), nil
}

func unmarshalTime(text []byte, t *time.Time, format string) (err error) {
	s := string(bytes.TrimSpace(text))
	*t, err = time.Parse(format, s)
	return
}

func (t ISODate) String() string {
	return (time.Time)(t).Format("02.01.2006")
}

type DateTime time.Time

func (t *DateTime) UnmarshalText(text []byte) (err error) {
	if len(text) == 10 {
		err = unmarshalTime(text, (*time.Time)(t), "2006-01-02")
		return
	}
	err = unmarshalTime(text, (*time.Time)(t), "2006-01-02T15:04:05")
	if err != nil {
		err = unmarshalTime(text, (*time.Time)(t), "2006-01-02T15:04:05.000")
	}
	return
}

func (t DateTime) MarshalText() ([]byte, error) {
	return []byte((time.Time)(t).Format(time.RFC3339)), nil
}

func (t DateTime) String() string {
	return (time.Time)(t).Format(time.RFC3339)
}

type SignedAct struct {
	ID                  int      `json:"id"`
	CreateDt            string   `json:"createDt"`
	SignDt              string   `json:"signDt"`
	ActPeriod           string   `json:"actPeriod"`
	ActNum              string   `json:"actNum"`
	StateCode           string   `json:"stateCode"`
	StateName           string   `json:"stateName"`
	OsiID               int      `json:"osiId"`
	OsiName             string   `json:"osiName"`
	OsiIdn              string   `json:"osiIdn"`
	OsiAddress          string   `json:"osiAddress"`
	OsiPhone            string   `json:"osiPhone"`
	OsiRegistrationDate string   `json:"osiRegistrationDate"`
	ApartCount          int      `json:"apartCount"`
	PlanAccuralID       int      `json:"planAccuralId"`
	Amount              *float64 `json:"amount"`
	Comission           *float64 `json:"comission"`
	Debt                *float64 `json:"debt"`
}

type ActDoc struct {
	DocTypeCode   string `json:"docTypeCode"`
	DocTypeNameRu string `json:"docTypeNameRu"`
	DocTypeNameKz string `json:"docTypeNameKz"`
	Scan          struct {
		FileName string `json:"fileName"`
		ID       int    `json:"id"`
	} `json:"scan"`
	ID int `json:"id"`
}
