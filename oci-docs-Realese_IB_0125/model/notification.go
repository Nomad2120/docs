package model

import (
	"fmt"
	"strings"
	"time"
)

type DebtorNotificationRequest struct {
	AbonentID int `json:"abonentId" uri:"abonentId" binding:"required"`
}

type DebtorNotificationResponse struct {
	BaseResponse
	Result DebtorNotificationResult `json:"result"`
}

type DebtorNotificationResult struct {
	Address       string          `json:"address"`
	DebtDate      ISODate         `json:"debtDate"`
	Flat          string          `json:"flat"`
	OsiChairman   string          `json:"osiChairman"`
	OsiName       string          `json:"osiName"`
	OsiID         int             `json:"osiId"`
	ServicesDebts []ServicesDebts `json:"servicesDebts"`
	ContactPhone  string          `json:"contactPhone"`
	Number        string          `json:"number"`
	AddressKZ     string          `json:"addressKz"`
	UnionTypeRu   string          `json:"unionTypeRu"`
	UnionTypeKz   string          `json:"unionTypeKz"`
}

type ServicesDebts struct {
	Saldo         float64 `json:"saldo"`
	SaldoString   string  `json:"saldoString"`
	ServiceName   string  `json:"serviceName"`
	ServiceNameKZ string  `json:"serviceNameKz"`
}

// FillDebetorResult -
type FillDebetorResult struct {
	DOCBase64 string `json:"docBase64"`
}

// FillDebetorResponse -
type FillDebetorResponse struct {
	BaseResponse
	Result *FillDebetorResult `json:"result"`
}

func (r *DebtorNotificationResult) GetDay() string {
	return fmt.Sprintf("%02d", time.Time(r.DebtDate).Day())
}

func (r *DebtorNotificationResult) GetYear() string {
	return fmt.Sprintf("%d", time.Time(r.DebtDate).Year())
}

func (r *DebtorNotificationResult) GetMonth() string {
	return formatMonth(time.Time(r.DebtDate))
}

func (r *DebtorNotificationResult) GetMonthKZ() string {
	return formatMonthKZ(time.Time(r.DebtDate))
}

func (r *DebtorNotificationResult) GetFullDate() string {
	return fmt.Sprintf("« %d » %s %d", time.Time(r.DebtDate).Day(), formatMonth(time.Time(r.DebtDate)), time.Time(r.DebtDate).Year())
}

func (r *DebtorNotificationResult) GetFullDateKZ() string {
	return fmt.Sprintf("%d жылғы « %d » %s ", time.Time(r.DebtDate).Year(), time.Time(r.DebtDate).Day(), formatMonthKZ(time.Time(r.DebtDate)))
}

func formatMonth(t time.Time) string {
	var months = [...]string{
		"января", "февраля", "марта", "апреля", "мая", "июня",
		"июля", "августа", "сентября", "октября", "ноября", "декабря",
	}
	return months[t.Month()-1]
}

func formatMonthNom(t time.Time) string {
	var months = [...]string{
		"январь", "февраль", "март", "апрель", "май", "июнь",
		"июль", "август", "сентябрь", "октябрь", "ноябрь", "декабрь",
	}
	return months[t.Month()-1]
}

func formatMonthKZ(t time.Time) string {
	var months = [...]string{
		"қаңтар", "ақпан", "наурыз", "сәуір", "мамыр", "маусым",
		"шілде", "тамыз", "қыркүйек", "қазан", "қараша", "желтоқсан",
	}
	return months[t.Month()-1]
}

func formatMonthNomKZ(t time.Time) string {
	var months = [...]string{
		"қаңтар", "ақпан", "наурыз", "сәуір", "мамыр", "маусым",
		"шілде", "тамыз", "қыркүйек", "қазан", "қараша", "желтоқсан",
	}
	return months[t.Month()-1]
}

func formatBeginMonthKZ(t time.Time) string {
	var months = [...]string{
		"қаңтардан", "ақпаннан", "наурыздан", "сәуірден", "мамырдан", "маусымнан",
		"шілдеден", "тамыздан", "қыркүйектен", "қазаннан", "қарашадан", "желтоқсаннан",
	}
	return months[t.Month()-1]
}

func formatEndMonthKZ(t time.Time) string {
	var months = [...]string{
		"қаңтарға", "ақпанға", "наурызға", "сәуірге", "мамырға", "маусымға",
		"шілдеге", "тамызға", "қыркүйекке", "қазанға", "қарашаға", "желтоқсанға",
	}
	return months[t.Month()-1]
}

type NotaryApplicationRequest struct {
	AbonentID int `json:"abonentId" uri:"abonentId" binding:"required"`
}

type NotaryApplicationResponse struct {
	BaseResponse
	Result NotaryApplicationResult `json:"result"`
}

type NotaryApplicationResult struct {
	OsiID            int                         `json:"osiId"`
	OsiName          string                      `json:"osiName"`
	OsiIdn           string                      `json:"osiIdn"`
	OsiAddress       string                      `json:"osiAddress"`
	OsiAddressKZ     string                      `json:"osiAddressKz"`
	OsiChairman      string                      `json:"osiChairman"`
	OsiPhone         string                      `json:"osiPhone"`
	AbonentName      string                      `json:"abonentName"`
	AbonentIdn       string                      `json:"abonentIdn"`
	AbonentFlat      string                      `json:"abonentFlat"`
	AbonentAddress   string                      `json:"abonentAddress"`
	AbonentAddressKZ string                      `json:"abonentAddressKz"`
	AbonentPhone     string                      `json:"abonentPhone"`
	DebtDate         ISODate                     `json:"debtDate"`
	TotalDebt        float64                     `json:"totalDebt"`
	TotalDebtString  string                      `json:"totalDebtString"`
	Registry         []NotaryApplicationRegistry `json:"registry"`
	UnionTypeRu      string                      `json:"unionTypeRu"`
	UnionTypeKz      string                      `json:"unionTypeKz"`
}

type NotaryApplicationRegistry struct {
	ServiceName     string                   `json:"serviceName"`
	ServiceNameKZ   string                   `json:"serviceNameKz"`
	TotalDebt       float64                  `json:"totalDebt"`
	TotalDebtString string                   `json:"totalDebtString"`
	Debts           []NotaryApplicationDebts `json:"debts"`
}

type NotaryApplicationDebts struct {
	Number         int      `json:"number"`
	Period         DebtDate `json:"period"`
	Debt           float64  `json:"debt"`
	CumulativeDebt float64  `json:"cumulativeDebt"`
}

// FillNotaryApplicationResult -
type FillNotaryApplicationResult struct {
	DOCBase64 string `json:"docBase64"`
}

// FillNotaryApplicationResponse -
type FillNotaryApplicationResponse struct {
	BaseResponse
	Result *FillNotaryApplicationResult `json:"result"`
}

func parseDate(date string) (time.Time, error) {
	r := strings.NewReplacer(
		"Январь", "January",
		"Февраль", "February",
		"Март", "March",
		"Апрель", "April",
		"Май", "May",
		"Июнь", "June",
		"Июль", "July",
		"Август", "August",
		"Сентябрь", "September",
		"Октябрь", "October",
		"Ноябрь", "November",
		"Декабрь", "December")
	return time.Parse("January 2006", r.Replace(date))
}

type DebtDate time.Time

func (t *DebtDate) UnmarshalText(text []byte) error {
	tm, err := parseDate(string(text))
	if err != nil {
		return err
	}
	*t = DebtDate(tm)
	return nil
}

func (t DebtDate) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%s %d", formatMonth(time.Time(t)), time.Time(t).Year())), nil
}

func (t *DebtDate) String() string {
	return fmt.Sprintf("%s %d", formatMonth(time.Time(*t)), time.Time(*t).Year())
}

func (na *NotaryApplicationResult) GetDateBegin() string {
	if len(na.Registry) == 0 || len(na.Registry[0].Debts) == 0 {
		return ""
	}
	return formatTime(time.Time(na.Registry[0].Debts[0].Period))
}

func (na *NotaryApplicationResult) GetDateEnd() string {
	if len(na.Registry) == 0 || len(na.Registry[0].Debts) == 0 {
		return ""
	}
	return formatTime(time.Time(na.Registry[len(na.Registry)-1].Debts[len(na.Registry[len(na.Registry)-1].Debts)-1].Period))
}

func (na *NotaryApplicationResult) GetDebtDate() string {
	dt := time.Now()
	if dt.Day() > 25 {
		dt = time.Date(dt.Year(), dt.Month(), 25, 0, 0, 0, 0, time.Local)
	} else {
		dt = dt.AddDate(0, -1, 0)
		dt = time.Date(dt.Year(), dt.Month(), 25, 0, 0, 0, 0, time.Local)
	}
	return formatTime(dt)
}

func (na *NotaryApplicationResult) GetChairman() string {
	if na.OsiChairman == "" {
		return "____________________________"
	}
	return na.OsiChairman
}

func (na *NotaryApplicationResult) FormatDateNom(d DebtDate) string {
	t := time.Time(d)
	return fmt.Sprintf("%s %04d г.", formatMonthNom(t), t.Year())
}

func (na *NotaryApplicationResult) FormatDate(d ISODate) string {
	t := time.Time(d)
	return fmt.Sprintf("%02d %s %04d г.", t.Day(), formatMonth(t), t.Year())
}

func formatTime(t time.Time) string {
	return fmt.Sprintf("%02d %s %04d г.", t.Day(), formatMonth(t), t.Year())
}
