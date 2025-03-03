package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config -
type Config struct {
	ServerPort           int      `json:"serverPort" env:"OCI_SERVER_PORT"`
	LogLevel             string   `json:"logLevel" env:"OCI_LOG_LEVEL"`
	LogDir               string   `json:"logDir" env:"OCI_LOG_DIR"`
	LoggingRequest       bool     `json:"loggingRequest" env:"OCI_LOGGING_REQUEST"`
	LoggingClientRequest bool     `json:"loggingClientRequest" env:"OCI_LOGGING_CLIENT_REQUEST"`
	CoreURL              string   `json:"apiURL" env:"OCI_CORE_URL"`
	CoreTimeout          Duration `json:"apiTimeout" env:"OCI_CORE_TIMEOUT"`
	CoreDBHost           string   `json:"coreDBHost" env:"OCI_CORE_DB_HOST"`
	CoreDBPort           int      `json:"coreDBPort" env:"OCI_CORE_DB_PORT"`
	CoreDBName           string   `json:"coreDBName" env:"OCI_CORE_DB_NAME"`
	CoreDBUser           string   `json:"coreDBUser" env:"OCI_CORE_DB_USER"`
	CoreDBPassword       string   `json:"coreDBPassword" env:"OCI_CORE_DB_PASSWORD"`
	InvocesPageURL       string   `json:"invocesPageURL" env:"OCI_INVOICES_PAGE_URL"`
	EarURL               string   `json:"earURL" env:"OCI_EAR_URL"`
	OsiStorePass         string   `json:"osiStorePass" env:"OCI_OSI_STORE_PASS"`
	QAZAFNStorePass      string   `json:"qazafnStorePass" env:"OCI_QAZAFN_STORE_PASS"`
	NurtauStorePass      string   `json:"nurtauStorePass" env:"OCI_NURTAU_STORE_PASS"`
}

func (c *Config) readFromFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&c)
	if err != nil {
		log.Fatal(err)
	}
}

// New ...
func New(path string) *Config {
	var config Config
	envFile := filepath.Join(filepath.Dir(path), ".env")
	if _, err := os.Stat(envFile); !os.IsNotExist(err) {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	config.mapEnv()
	return &config
}

func (c *Config) mapEnv() {
	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		tag := t.Field(i).Tag
		envTag := tag.Get("env")
		value := os.Getenv(envTag)
		if value != "" {
			switch f.Kind() {

			case reflect.String:
				f.SetString(value)

			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				i, err := strconv.ParseInt(value, 10, f.Type().Bits())
				if err == nil {
					if !f.OverflowInt(i) {
						f.SetInt(i)
					}
				}

			case reflect.Bool:
				b, err := strconv.ParseBool(value)
				if err == nil {
					f.SetBool(b)
				}

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				ui, err := strconv.ParseUint(value, 10, f.Type().Bits())
				if err == nil {
					if !f.OverflowUint(ui) {
						f.SetUint(ui)
					}
				}

			case reflect.Float32, reflect.Float64:
				fl, err := strconv.ParseFloat(value, f.Type().Bits())
				if err == nil {
					if !f.OverflowFloat(fl) {
						f.SetFloat(fl)
					}
				}
			case reflect.Struct:
				if f.Type().Name() == "Duration" {
					d, err := time.ParseDuration(value)
					fd := f.FieldByName("Duration")
					if err == nil && fd.CanSet() {
						fd.Set(reflect.ValueOf(d))
					}
				}
			case reflect.Slice:
				var ret []string
				if err := json.Unmarshal([]byte(value), &ret); err != nil {
					panic(err)
				}
				// slice := reflect.MakeSlice(reflect.TypeOf([]string{}), len(ret), len(ret))
				// slice.Set(reflect.ValueOf(ret))
				// reflect.Copy(slice, ret)
				f.Set(reflect.ValueOf(ret))

			default:
			}
		}
	}
}

// Duration -
type Duration struct {
	time.Duration
}

// UnmarshalJSON -
func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	d.Duration, err = time.ParseDuration(strings.Trim(string(b), `"`))
	return
}

// MarshalJSON -
func (d *Duration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s`, d.String())), nil
}
