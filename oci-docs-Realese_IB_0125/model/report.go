package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/now"
	"golang.org/x/exp/slices"
)

type OSVResponse struct {
	BaseResponse
	Result OSVResult `json:"result"`
}

type OSVResult struct {
	Abonents []OSVAbonent `json:"abonents"`
}

type Abonent struct {
	AbonentID    int
	AbonentName  string
	Flat         string
	AreaTypeCode string
	Service      OSVService
}

func (o *OSVResult) GetAbonentsWithService(serviceName string) []Abonent {
	result := make([]Abonent, 0, 32)
	for _, item := range o.Abonents {
		if service := item.FindService(serviceName); service != nil {
			result = append(result, Abonent{
				AbonentID:    item.AbonentID,
				AbonentName:  item.AbonentName,
				Flat:         item.Flat,
				AreaTypeCode: item.AreaTypeCode,
				Service:      *service,
			})
		}
	}
	return result
}

type OSVAbonent struct {
	AbonentID    int    `json:"abonentId"`
	AbonentName  string `json:"abonentName"`
	Owner        string `json:"owner"`
	Flat         string `json:"flat"`
	AreaTypeCode string
	Services     []OSVService `json:"services"`
}

func (a *OSVAbonent) FindService(name string) *OSVService {
	for _, s := range a.Services {
		if strings.Contains(strings.ToLower(s.ServiceName), strings.ToLower(name)) {
			return &s
		}
	}
	return nil
}

type OSVService struct {
	ServiceName       string      `json:"serviceName"`
	ServiceNameKZ     string      `json:"serviceNameKz"`
	Begin             float64     `json:"begin"`
	Debet             float64     `json:"debet"`
	DebetWithoutFixes float64     `json:"debetWithoutFixes"`
	SumOfFixes        float64     `json:"sumOfFixes"`
	SumOfAccurals     float64     `json:"sumOfAccurals"`
	SumOfFines        float64     `json:"sumOfFines"`
	Kredit            float64     `json:"kredit"`
	End               float64     `json:"end"`
	Fixes             []FixesInfo `json:"fixes"`
}

type OSIResponse struct {
	BaseResponse
	Result OSIResult `json:"result"`
}

type AllOSIResponse struct {
	BaseResponse
	Result []OSIResult `json:"result"`
}

type OSIResult struct {
	ID                    int    `json:"id"`
	IsActive              bool   `json:"isActive"`
	IsLaunched            bool   `json:"isLaunched"`
	RegistrationID        int    `json:"registrationId"`
	WizardStep            string `json:"wizardStep"`
	Rca                   string `json:"rca"`
	HouseStateNameRu      string `json:"houseStateNameRu"`
	HouseStateNameKz      string `json:"houseStateNameKz"`
	Name                  string `json:"name"`
	Idn                   string `json:"idn"`
	Address               string `json:"address"`
	AddressKZ             string `json:"addressKz"`
	Phone                 string `json:"phone"`
	Email                 string `json:"email"`
	ConstructionYear      int    `json:"constructionYear"`
	ConstructionMaterial  string `json:"constructionMaterial"`
	Floors                int    `json:"floors"`
	ApartCount            int    `json:"apartCount"`
	HouseStateCode        string `json:"houseStateCode"`
	PersonalHeating       bool   `json:"personalHeating"`
	PersonalHotWater      bool   `json:"personalHotWater"`
	PersonalElectricPower bool   `json:"personalElectricPower"`
	Gasified              bool   `json:"gasified"`
	CoefUnlivingArea      int    `json:"coefUnlivingArea"`
	Fio                   string `json:"fio"`
	UnionTypeRu           string `json:"unionTypeRu"`
	UnionTypeKz           string `json:"unionTypeKz"`
}

func (r OSIResult) OSIName() string {
	if !strings.HasPrefix(r.Name, r.UnionTypeRu) {
		return fmt.Sprintf(`%s "%s"`, r.UnionTypeRu, r.Name)
	}
	return fmt.Sprintf(`"%s"`, r.Name)
}

// ReportRequest -
type ReportRequest struct {
	ID int `json:"id" uri:"id" binding:"required"`
}

// OSVReportRequest -
type OSVReportRequest struct {
	ID         int       `json:"id" uri:"id"`
	Begin      time.Time `json:"begin" form:"begin" binding:"required"`
	End        time.Time `json:"end" form:"end" binding:"required"`
	ForAbonent bool      `json:"forAbonent" form:"forAbonent"`
}

// ReportResponse -
type ReportResponse struct {
	BaseResponse
	Result string `json:"result"`
}

// PaymentsReportRequest -
type PaymentsReportRequest struct {
	ID    int       `json:"id"`
	Begin time.Time `json:"begin" form:"begin" binding:"required"`
	End   time.Time `json:"end" form:"end" binding:"required"`
}

type PaymentsResponse struct {
	BaseResponse
	Result []PaymentInfo `json:"result"`
}

type PaymentInfo struct {
	Dt          string  `json:"dt"`
	AbonentName string  `json:"abonentName"`
	Flat        string  `json:"flat"`
	ServiceName string  `json:"serviceName"`
	Amount      float64 `json:"amount"`
	BankName    string  `json:"bankName"`
}

type PaymentOrdersRequest struct {
	ID    int       `json:"osiId" uri:"id"`
	Begin time.Time `json:"begin" form:"begin" binding:"required"`
	End   time.Time `json:"end" form:"end" binding:"required"`
}

type PaymentOrdersResponse struct {
	BaseResponse
	Result []PaymentOrder `json:"result"`
}

type PaymentOrder struct {
	Amount           *float64 `json:"amount"`
	AmountToTransfer *float64 `json:"amountToTransfer"`
	BankName         string   `json:"bankName"`
	IBAN             string   `json:"iban"`
	ComisBank        *float64 `json:"comisBank"`
	ComisOur         *float64 `json:"comisOur"`
	Date             ISODate  `json:"date"`
}

type AbonentOSVReportRequest struct {
	ID        int    `json:"osiId" uri:"id" binding:"required"`
	AbonentID int    `json:"abonentId" uri:"abonentId" binding:"required"`
	Flat      string `json:"flat" form:"flat"`
}

type AbonentOSVResponse struct {
	BaseResponse
	Result []AbonentOSVResult `json:"result"`
}

type AbonentOSVResult struct {
	Period   string       `json:"period"`
	Services []OSVService `json:"services"`
}

// FixesRequest -
type FixesRequest struct {
	ID    int       `json:"id"`
	Begin time.Time `json:"begin" form:"begin" binding:"required"`
	End   time.Time `json:"end" form:"end" binding:"required"`
}

type FixesResponse struct {
	BaseResponse
	Result []FixesInfo `json:"result"`
}

type FixesInfo struct {
	Dt               string  `json:"dt"`
	AbonentName      string  `json:"abonentName"`
	Flat             string  `json:"flat"`
	ServiceName      string  `json:"serviceName"`
	ServiceGroupName string  `json:"serviceGroupName"`
	Reason           string  `json:"reason"`
	Amount           float64 `json:"amount"`
}

type AccountsReportRequest struct {
	Path struct {
		ID int `json:"id" uri:"id" binding:"required"`
	}
	Query struct {
		Language string `json:"language" form:"language"`
	}
}

type FillReportRequest struct {
	Query struct {
		Language string `json:"language" form:"language"`
	}
	Body AccountReportsResult
}

type AccountReportsResponse struct {
	BaseResponse
	Result AccountReportsResult `json:"result"`
}

type AccountReportsResult struct {
	Period      ISODate `json:"period"` //"2023-11-30T09:34:03.629Z",
	OsiAddress  string  `json:"osiAddress"`
	OsiName     string  `json:"osiName"`
	UnionTypeRu string  `json:"unionTypeRu"`
	UnionTypeKz string  `json:"unionTypeKz"`
	//UnionType      AccountReportsUnionType `json:"unionType"`
	Categories     []AccountReportCategory `json:"categories"`
	UnionTypeTitle string                  `json:"-"`
	UnionName      string                  `json:"-"`
	Signer         string                  `json:"signer"`
}

type AccountReportsUnionType struct {
	ID     int    `json:"id"`
	NameRu string `json:"nameRu"`
	NameKz string `json:"nameKz"`
}

type AccountReportCategory struct {
	Number string  `json:"number"`
	NameRu string  `json:"nameRu"`
	NameKz string  `json:"nameKz"`
	Amount float64 `json:"amount"`
	Name   string  `json:"-"`
}

type ReportResult struct {
	DOCBase64 string `json:"docBase64"`
}

func (ar *AccountReportsResult) Calc(lang, signer string) {
	if lang == "" {
		lang = "ru"
	}
	if strings.EqualFold(lang, "ru") {
		ar.UnionTypeTitle = ar.UnionTypeRu
		ar.UnionName = ar.OsiName

	} else {
		ar.UnionTypeTitle = "МИБ"
		ar.UnionName = ar.OsiName
	}
	for i, category := range ar.Categories {
		if strings.EqualFold(lang, "ru") {
			ar.Categories[i].Name = calcNameRu(&category)
		} else {
			ar.Categories[i].Name = calcNameKz(&category)
		}
	}
	if signer != "" {
		ar.Signer = signer
	}
}

func calcNameRu(c *AccountReportCategory) string {
	switch c.Number {
	case "1", "2":
		return "<strong>" + c.NameRu + "</strong>"
	case "3", "4":
		if idx := strings.Index(c.NameRu, " собственников"); idx >= 0 {
			return "<strong>" + c.NameRu[:idx] + "</strong>" + c.NameRu[idx:]
		}
		return c.NameRu
	case "5", "6":
		if idx := strings.Index(c.NameRu, "в том числе"); idx >= 0 {
			return "<strong>" + c.NameRu[:idx] + "</strong>" + c.NameRu[idx:]
		}
		return c.NameRu

	default:
		return c.NameRu
	}
}

func calcNameKz(c *AccountReportCategory) string {
	switch c.Number {
	case "1", "2":
		return "<strong>" + c.NameKz + "</strong>"
	case "5", "6":
		if idx := strings.Index(c.NameKz, "оның ішінде"); idx >= 0 {
			return "<strong>" + c.NameKz[:idx] + "</strong>" + c.NameKz[idx:]
		}
		return c.NameKz

	default:
		return c.NameKz
	}
}

func (ar *AccountReportsResult) GetDateBegin() string {
	location, _ := time.LoadLocation("Asia/Aqtobe")

	config := &now.Config{
		WeekStartDay: time.Monday,
		TimeLocation: location,
		TimeFormats:  []string{"2006-01-02 15:04:05"},
	}
	dt := config.With(time.Time(ar.Period)).BeginningOfMonth()
	return fmt.Sprintf(`"%02d" %s %d`, dt.Day(), formatMonth(dt), dt.Year())
}

func (ar *AccountReportsResult) GetDateBeginKZ() string {
	location, _ := time.LoadLocation("Asia/Aqtobe")

	config := &now.Config{
		WeekStartDay: time.Monday,
		TimeLocation: location,
		TimeFormats:  []string{"2006-01-02 15:04:05"},
	}
	dt := config.With(time.Time(ar.Period)).BeginningOfMonth()
	//2023 жылғы 01 желтоқсаннан
	return fmt.Sprintf("%d жылғы %02d %s", dt.Year(), dt.Day(), formatBeginMonthKZ(dt))
}

func (ar *AccountReportsResult) GetDateEnd() string {
	location, _ := time.LoadLocation("Asia/Aqtobe")

	config := &now.Config{
		WeekStartDay: time.Monday,
		TimeLocation: location,
		TimeFormats:  []string{"2006-01-02 15:04:05"},
	}
	dt := config.With(time.Time(ar.Period)).EndOfMonth()
	return fmt.Sprintf(`"%02d" %s %d`, dt.Day(), formatMonth(dt), dt.Year())
}

func (ar *AccountReportsResult) GetDateEndKZ() string {
	location, _ := time.LoadLocation("Asia/Aqtobe")

	config := &now.Config{
		WeekStartDay: time.Monday,
		TimeLocation: location,
		TimeFormats:  []string{"2006-01-02 15:04:05"},
	}
	dt := config.With(time.Time(ar.Period)).EndOfMonth()
	//2023 жылғы 31 желтоқсанға
	return fmt.Sprintf("%d жылғы %02d %s", dt.Year(), dt.Day(), formatEndMonthKZ(dt))
}

func (ar *AccountReportsResult) GetReportDate() string {
	dt := time.Now()
	return fmt.Sprintf(`"%02d" %s %d`, dt.Day(), formatMonth(dt), dt.Year())
}

func (ar *AccountReportsResult) GetReportDateKZ() string {
	dt := time.Now()
	//2024 жылғы 23 қаңтар
	return fmt.Sprintf("%d жылғы %02d %s", dt.Year(), dt.Day(), formatMonthKZ(dt))
}

func (ar *AccountReportsResult) GetClassTableRow(index int) string {
	numbers := []string{"1.", "2.", "3.", "4.", "1", "2", "3", "4"}
	if index < len(ar.Categories) {
		if slices.Contains(numbers, ar.Categories[index].Number) {
			return "text-strong"
		}
	}
	return ""
}

func (ar *AccountReportsResult) GetClassCol3(index int) string {
	numbers := []string{"1.", "2.", "3.", "4.", "5.", "6.", "1", "2", "3", "4", "5", "6", "6.1.", "6.2.", "6.2.1.", "6.1", "6.2", "6.2.1"}
	if index < len(ar.Categories) {
		if slices.Contains(numbers, ar.Categories[index].Number) {
			return "text-strong"
		}
	}
	return ""
}
