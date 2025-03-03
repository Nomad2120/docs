package common

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"gitlab.enterprise.qazafn.kz/oci/oci-docs/config"
	model "gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
	"go.uber.org/zap"

	"golang.org/x/exp/slices"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

// InitHTTPClient - инициализация клиента HTTP
func InitHTTPClient(conf *config.Config, log *zap.SugaredLogger) *http.Client {
	roundTripper := http.DefaultTransport
	transport, _ := roundTripper.(*http.Transport)
	transport.MaxIdleConns = 200
	transport.MaxIdleConnsPerHost = 200

	if conf.LoggingClientRequest {
		return &http.Client{
			Transport: NewLoggedRoundTripper(transport, NewDefaultLogger(log)),
			Timeout:   conf.CoreTimeout.Duration,
		}
	}

	return &http.Client{
		Transport: transport,
	}
}

// UUID -
func UUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}

// Pstr -
func Pstr(s string) *string {
	return &s
}

func PEMToCMS(pem string) (string, error) {
	n := 64
	buf := make([]byte, n+1)
	rest := 0
	result := "-----BEGIN CMS-----\n"
	for i, b := range []byte(pem) {
		rest = (i + 1) % n
		if rest == 0 {
			buf[n-1] = b
			buf[n] = '\n'
			result += string(buf)
			continue
		}
		buf[rest-1] = b
	}
	if rest > 0 {
		buf[rest] = '\n'
		result += string(buf[:rest+1])
	}
	result += "-----END CMS-----\n"
	return result, nil
}

func ExtractIIN(subject string) string {
	p := strings.Split(subject, ",")
	for _, item := range p {
		kv := strings.Split(item, "=")
		if len(kv) == 2 && len(kv[1]) > 3 && strings.ToUpper(kv[1][:3]) == "BIN" {
			return kv[1][3:]
		}
	}
	for _, item := range p {
		kv := strings.Split(item, "=")
		if len(kv) == 2 && strings.ToUpper(kv[0]) == "SERIALNUMBER" {
			return kv[1][3:]
		}
	}
	return ""
}

// LoadFixture - load fixtures
func LoadFixture(t *testing.T, name string) []byte {
	bytes, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

func FormatAmount(amount float64) string {
	p := message.NewPrinter(language.Russian)
	return p.Sprint(number.Decimal(amount, number.MaxFractionDigits(2)))
}

func FormatMonth(t time.Time) string {
	var months = [...]string{
		"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
	}
	return fmt.Sprintf("%s %d", months[t.Month()-1], t.Year())
}

func CurrentMonthYear() string {
	return strings.ToLower(FormatMonth(time.Now()))
}

func FormatDate(tm time.Time, layout string) string {
	return tm.Format(layout)
}

func GetFlatType(code string) string {
	if code == "NON_RESIDENTIAL" {
		return " Нежилое"
	}
	return ""
}

func FormatPhone(phone string) string {
	if len(phone) < 10 {
		return phone
	}

	return fmt.Sprintf("+7 %s %s %s %s", phone[0:3], phone[3:6], phone[6:8], phone[8:10])
}

func LessAlphaNum(a, b string) bool {
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		numeric, numA, numB := false, 0, 0
		for i < len(a) && a[i] >= '0' && a[i] <= '9' {
			numA = numA*10 + int(a[i]) - '0'
			numeric = true
			i++
		}
		for j < len(b) && b[j] >= '0' && b[j] <= '9' {
			numB = numB*10 + int(b[j]) - '0'
			numeric = true
			j++
		}
		if numeric {
			if numA != numB {
				return numA < numB
			}
			continue
		}
		if a[i] != b[j] {
			return a[i] < b[j]
		}
		i++
		j++
	}
	return i == len(a) && j != len(b)
}

func ExtractYearMonth(s string) (year int, month int) {
	var months = []string{
		"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь",
	}
	a := strings.Split(s, " ")
	if len(a) < 2 {
		return
	}
	idx := slices.Index(months, a[0])
	if idx < 0 {
		return
	}
	month = idx + 1
	year, _ = strconv.Atoi(a[1])
	return
}

func MapFixesToOSV(flat string, osv []model.AbonentOSVResult, fixes []model.FixesInfo) []model.AbonentOSVResult {
	for i, item := range osv {
		year, month := ExtractYearMonth(item.Period)

		for j, service := range item.Services {
			if service.SumOfFixes != 0 {
				for _, fix := range fixes {
					osv[i].Services[j].Fixes = make([]model.FixesInfo, 0, 8)
					dt, _ := time.Parse("2006-01-02T15:04:05", fix.Dt)
					if fix.Flat == flat && year == dt.Year() && month == int(dt.Month()) && fix.ServiceGroupName == service.ServiceName {
						osv[i].Services[j].Fixes = append(osv[i].Services[j].Fixes, fix)
					}
				}
			}

		}
	}
	return osv
}
