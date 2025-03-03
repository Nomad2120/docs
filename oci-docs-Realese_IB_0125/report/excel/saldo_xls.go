package excel

import (
	"fmt"
	"io"
	"sort"
	"time"
	"unicode/utf8"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/shopspring/decimal"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/common"
	model "gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
	report "gitlab.enterprise.qazafn.kz/oci/oci-docs/report"
)

type total struct {
	Begin        []string
	Debet        []string
	Fixes        []string
	Kredit       []string
	End          []string
	Fines        []string
	BeginAmount  decimal.Decimal
	DebetAmount  decimal.Decimal
	FixesAmount  decimal.Decimal
	KreditAmount decimal.Decimal
	EndAmount    decimal.Decimal
	FinesAmount  decimal.Decimal
}

type osvDebt struct {
	Flat         string
	AreaTypeCode string
	Services     map[string]float64
}

type saldoXLS struct {
	f             *excelize.File
	osi           *model.OSIResult
	osv           *model.OSVResult
	serviceGroups []model.ServiceGroup
	begin         time.Time
	end           time.Time
	isAbonent     bool
}

func NewSaldoXLS(osi *model.OSIResult, osv *model.OSVResult, serviceGroups []model.ServiceGroup, begin, end time.Time, isAbonent bool) report.Report {
	return &saldoXLS{
		f:             excelize.NewFile(),
		osi:           osi,
		osv:           osv,
		serviceGroups: serviceGroups,
		begin:         begin,
		end:           end,
		isAbonent:     isAbonent,
	}
}

func (r *saldoXLS) Render() error {
	const (
		top  = 1
		left = 0
	)

	sort.SliceStable(r.osv.Abonents, func(i, j int) bool {
		return len(r.osv.Abonents[i].Flat) < len(r.osv.Abonents[j].Flat) ||
			(len(r.osv.Abonents[i].Flat) == len(r.osv.Abonents[j].Flat) && r.osv.Abonents[i].Flat < r.osv.Abonents[j].Flat)
	})

	sheet := "Общая ведомость"
	// Создать новый лист
	r.f.NewSheet(sheet)

	r.f.DeleteSheet("Sheet1")

	if err := r.f.SetSheetViewOptions(sheet, 0, excelize.ZoomScale(145)); err != nil {
		panic(err)
	}

	titleStyle, err := r.f.NewStyle(`{
		"font":{"bold":false},
		"alignment":{
			"horizontal": "center"
			}
		}`)

	if err != nil {
		return err
	}

	headStyle, err := r.f.NewStyle(`{
		"font":{"bold":false},
		"alignment":{
			"horizontal": "center",
			"vertical": "center",
			"wrap_text":true
			},
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}
	]
	}`)
	if err != nil {
		return err
	}
	rowStyle, err := r.f.NewStyle(`{
		"number_format": 4,
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}]
		}`)
	if err != nil {
		return err
	}

	totalStyle, err := r.f.NewStyle(`{
		"font":{"bold":true},
		"number_format": 4,
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}]
		}`)
	if err != nil {
		return err
	}

	serviceStyle, err := r.f.NewStyle(`{"font":{"italic":false},
	"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}]
	}`)
	if err != nil {
		return err
	}
	flatStyle, err := r.f.NewStyle(`{
		"alignment":{
			"horizontal": "center",
			"vertical": "top"
		},
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}
		]}`)
	if err != nil {
		return err
	}

	flatTotalStyle, err := r.f.NewStyle(`{
		"font":{"bold":true},
		"alignment":{
			"horizontal": "center",
			"vertical": "top"
		},
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}
		]}`)
	if err != nil {
		return err
	}
	abonentStyle, err := r.f.NewStyle(`{
		"alignment":{
			"horizontal": "left",
			"vertical": "top",
			"wrap_text":true
		},
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}
		]}`)
	if err != nil {
		return err
	}

	r.f.SetCellValue(sheet, makeAxis(left, top), r.osi.OSIName())
	r.f.MergeCell(sheet, makeAxis(left, top), makeAxis(left+15, top))
	r.f.SetCellValue(sheet, makeAxis(left, top+1), fmt.Sprintf("Оборотно-сальдовая ведомость за период %s по %s", r.begin.Format("02/01/2006"), r.end.Format("02/01/2006")))
	r.f.MergeCell(sheet, makeAxis(left, top+1), makeAxis(left+15, top+1))
	r.f.SetCellStyle(sheet, makeAxis(left, top), makeAxis(left+15, top+1), titleStyle)

	r.f.SetCellValue(sheet, makeAxis(left, top+2), "Адрес: "+r.osi.Address)
	r.f.MergeCell(sheet, makeAxis(left, top+2), makeAxis(left+15, top+2))
	r.f.SetCellStyle(sheet, makeAxis(left, top+2), makeAxis(left+15, top+2), titleStyle)

	col := left
	row := top + 3
	corr := 1
	if r.isAbonent {
		corr = -1
	}
	r.f.SetCellValue(sheet, makeAxis(col, row), "№ помещения")
	r.f.SetColWidth(sheet, ColIndexToLetters(col), ColIndexToLetters(col), 15)
	if !r.isAbonent {
		r.f.SetCellValue(sheet, makeAxis(col+1, row), "ФИО собственника")
		r.f.SetColWidth(sheet, ColIndexToLetters(col+1), ColIndexToLetters(col+1), 25)
		r.f.SetCellValue(sheet, makeAxis(col+2, row), "Владелец")
		r.f.SetColWidth(sheet, ColIndexToLetters(col+2), ColIndexToLetters(col+2), 15)
	}

	r.f.SetCellValue(sheet, makeAxis(col+2+corr, row), "Услуга")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+2+corr), ColIndexToLetters(col+2+corr), 30)
	r.f.SetCellValue(sheet, makeAxis(col+3+corr, row), "Сальдо на начало, тг.")
	r.f.MergeCell(sheet, makeAxis(col+3+corr, row), makeAxis(col+4+corr, row))
	r.f.SetColWidth(sheet, ColIndexToLetters(col+3+corr), ColIndexToLetters(col+4+corr), 15)
	r.f.SetCellValue(sheet, makeAxis(col+5+corr, row), "Начислено за период, тг.")
	r.f.MergeCell(sheet, makeAxis(col+5+corr, row), makeAxis(col+6+corr, row))
	r.f.SetColWidth(sheet, ColIndexToLetters(col+5+corr), ColIndexToLetters(col+6+corr), 15)
	r.f.SetCellValue(sheet, makeAxis(col+7+corr, row), "Корректировки, тг.")
	r.f.MergeCell(sheet, makeAxis(col+7+corr, row), makeAxis(col+8+corr, row))
	r.f.SetColWidth(sheet, ColIndexToLetters(col+7+corr), ColIndexToLetters(col+8+corr), 15)

	r.f.SetCellValue(sheet, makeAxis(col+9+corr, row), "Оплачено за период, тг.")
	r.f.MergeCell(sheet, makeAxis(col+9+corr, row), makeAxis(col+10+corr, row))
	r.f.SetColWidth(sheet, ColIndexToLetters(col+9+corr), ColIndexToLetters(col+10+corr), 15)

	r.f.SetCellValue(sheet, makeAxis(col+11+corr, row), "Пеня, тг.")
	r.f.MergeCell(sheet, makeAxis(col+11+corr, row), makeAxis(col+12+corr, row))
	r.f.SetColWidth(sheet, ColIndexToLetters(col+11+corr), ColIndexToLetters(col+12+corr), 15)

	r.f.SetCellValue(sheet, makeAxis(col+13+corr, row), "Долг, тг.")
	r.f.MergeCell(sheet, makeAxis(col+13+corr, row), makeAxis(col+14+corr, row))
	r.f.SetColWidth(sheet, ColIndexToLetters(col+13+corr), ColIndexToLetters(col+14+corr), 15)

	r.f.SetCellStyle(sheet, makeAxis(col, row), makeAxis(col+14+corr, row), headStyle)

	totals := make(map[string]total)
	curRow := row + 1

	for _, abonent := range r.osv.Abonents {
		r.f.SetCellValue(sheet, makeAxis(col, curRow), abonent.Flat+common.GetFlatType(abonent.AreaTypeCode))
		if !r.isAbonent {
			r.f.SetCellValue(sheet, makeAxis(col+1, curRow), abonent.AbonentName)
			r.f.SetCellValue(sheet, makeAxis(col+2, curRow), abonent.Owner)
		}
		r.f.SetCellValue(sheet, makeAxis(col+2+corr, curRow), "")

		r.f.SetCellFormula(sheet, makeAxis(col+3+corr, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+4+corr, curRow+1), makeAxis(col+4+corr, curRow+len(abonent.Services))))
		r.f.SetCellFormula(sheet, makeAxis(col+5+corr, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+6+corr, curRow+1), makeAxis(col+6+corr, curRow+len(abonent.Services))))
		r.f.SetCellFormula(sheet, makeAxis(col+7+corr, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+8+corr, curRow+1), makeAxis(col+8+corr, curRow+len(abonent.Services))))
		r.f.SetCellFormula(sheet, makeAxis(col+9+corr, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+10+corr, curRow+1), makeAxis(col+10+corr, curRow+len(abonent.Services))))
		r.f.SetCellFormula(sheet, makeAxis(col+11+corr, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+12+corr, curRow+1), makeAxis(col+12+corr, curRow+len(abonent.Services))))
		r.f.SetCellFormula(sheet, makeAxis(col+13+corr, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+14+corr, curRow+1), makeAxis(col+14+corr, curRow+len(abonent.Services))))
		r.f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col+14+corr, curRow), rowStyle)

		curRow++
		for _, service := range abonent.Services {
			item, ok := totals[service.ServiceName]

			if !ok {
				item = total{
					Begin:        []string{makeAxis(col+4+corr, curRow)},
					Debet:        []string{makeAxis(col+6+corr, curRow)},
					Fixes:        []string{makeAxis(col+8+corr, curRow)},
					Kredit:       []string{makeAxis(col+10+corr, curRow)},
					End:          []string{makeAxis(col+12+corr, curRow)},
					Fines:        []string{makeAxis(col+14+corr, curRow)},
					BeginAmount:  decimal.NewFromFloat(service.Begin),
					DebetAmount:  decimal.NewFromFloat(service.SumOfAccurals),
					FixesAmount:  decimal.NewFromFloat(service.SumOfFixes),
					KreditAmount: decimal.NewFromFloat(service.Kredit),
					FinesAmount:  decimal.NewFromFloat(service.SumOfFines),
					EndAmount:    decimal.NewFromFloat(service.End),
				}
			} else {
				item.Begin = append(item.Begin, makeAxis(col+4+corr, curRow))
				item.Debet = append(item.Debet, makeAxis(col+6+corr, curRow))
				item.Fixes = append(item.Fixes, makeAxis(col+8+corr, curRow))
				item.Kredit = append(item.Kredit, makeAxis(col+10+corr, curRow))
				item.Fines = append(item.Fines, makeAxis(col+12+corr, curRow))
				item.End = append(item.End, makeAxis(col+14+corr, curRow))
				item.BeginAmount = item.BeginAmount.Add(decimal.NewFromFloat(service.Begin))
				item.DebetAmount = item.DebetAmount.Add(decimal.NewFromFloat(service.SumOfAccurals))
				item.FixesAmount = item.FixesAmount.Add(decimal.NewFromFloat(service.SumOfFixes))
				item.KreditAmount = item.KreditAmount.Add(decimal.NewFromFloat(service.Kredit))
				item.FinesAmount = item.FinesAmount.Add(decimal.NewFromFloat(service.SumOfFines))
				item.EndAmount = item.EndAmount.Add(decimal.NewFromFloat(service.End))
			}
			totals[service.ServiceName] = item

			r.f.SetCellValue(sheet, makeAxis(col+2+corr, curRow), service.ServiceName)
			r.f.MergeCell(sheet, makeAxis(col+2+corr, curRow), makeAxis(col+3+corr, curRow))

			r.f.SetCellValue(sheet, makeAxis(col+4+corr, curRow), service.Begin)
			r.f.SetCellValue(sheet, makeAxis(col+6+corr, curRow), service.SumOfAccurals)
			r.f.SetCellValue(sheet, makeAxis(col+8+corr, curRow), service.SumOfFixes)
			r.f.SetCellValue(sheet, makeAxis(col+10+corr, curRow), service.Kredit)
			r.f.SetCellValue(sheet, makeAxis(col+12+corr, curRow), service.SumOfFines)
			r.f.SetCellValue(sheet, makeAxis(col+14+corr, curRow), service.End)
			r.f.SetCellStyle(sheet, makeAxis(col+1+corr, curRow), makeAxis(col+14+corr, curRow), serviceStyle)
			curRow++
		}
		r.f.MergeCell(sheet, makeAxis(left, curRow-len(abonent.Services)-1), makeAxis(left, curRow-1))
		r.f.SetCellStyle(sheet, makeAxis(left, curRow-len(abonent.Services)-1), makeAxis(left, curRow-1), flatStyle)
		if !r.isAbonent {
			r.f.MergeCell(sheet, makeAxis(col+corr, curRow-len(abonent.Services)-1), makeAxis(col+corr, curRow-1))
			r.f.SetCellStyle(sheet, makeAxis(col+corr, curRow-len(abonent.Services)-1), makeAxis(col+corr, curRow-1), abonentStyle)
		}
		r.f.MergeCell(sheet, makeAxis(col+1+corr, curRow-len(abonent.Services)-1), makeAxis(col+1+corr, curRow-1))
		r.f.SetCellStyle(sheet, makeAxis(col+1+corr, curRow-len(abonent.Services)-1), makeAxis(col+1+corr, curRow-1), abonentStyle)

	}
	r.f.SetCellValue(sheet, makeAxis(col, curRow), "ИТОГО")

	r.f.SetCellFormula(sheet, makeAxis(col+3+corr, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+3+corr, row+1), makeAxis(col+3+corr, curRow-1)))
	r.f.SetCellFormula(sheet, makeAxis(col+5+corr, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+5+corr, row+1), makeAxis(col+5+corr, curRow-1)))
	r.f.SetCellFormula(sheet, makeAxis(col+7+corr, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+7+corr, row+1), makeAxis(col+7+corr, curRow-1)))
	r.f.SetCellFormula(sheet, makeAxis(col+9+corr, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+9+corr, row+1), makeAxis(col+9+corr, curRow-1)))
	r.f.SetCellFormula(sheet, makeAxis(col+11+corr, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+11+corr, row+1), makeAxis(col+11+corr, curRow-1)))
	r.f.SetCellFormula(sheet, makeAxis(col+13+corr, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+13+corr, row+1), makeAxis(col+13+corr, curRow-1)))
	r.f.SetCellStyle(sheet, makeAxis(col+1, curRow), makeAxis(col+14+corr, curRow), totalStyle)

	for key, value := range totals {
		curRow++
		r.f.SetCellValue(sheet, makeAxis(col+2+corr, curRow), key)
		r.f.MergeCell(sheet, makeAxis(col+2+corr, curRow), makeAxis(col+3+corr, curRow))

		// r.f.SetCellFormula(sheet, makeAxis(col+4+corr, curRow), fmt.Sprintf("=SUM(%s)", strings.Join(value.Begin, ",")))
		// r.f.SetCellFormula(sheet, makeAxis(col+6+corr, curRow), fmt.Sprintf("=SUM(%s)", strings.Join(value.Debet, ",")))
		// r.f.SetCellFormula(sheet, makeAxis(col+8+corr, curRow), fmt.Sprintf("=SUM(%s)", strings.Join(value.Kredit, ",")))
		// r.f.SetCellFormula(sheet, makeAxis(col+10+corr, curRow), fmt.Sprintf("=SUM(%s)", strings.Join(value.End, ",")))
		r.f.SetCellValue(sheet, makeAxis(col+4+corr, curRow), value.BeginAmount.InexactFloat64())
		r.f.SetCellValue(sheet, makeAxis(col+6+corr, curRow), value.DebetAmount.InexactFloat64())
		r.f.SetCellValue(sheet, makeAxis(col+8+corr, curRow), value.FixesAmount.InexactFloat64())
		r.f.SetCellValue(sheet, makeAxis(col+10+corr, curRow), value.KreditAmount.InexactFloat64())
		r.f.SetCellValue(sheet, makeAxis(col+12+corr, curRow), value.FinesAmount.InexactFloat64())
		r.f.SetCellValue(sheet, makeAxis(col+14+corr, curRow), value.EndAmount.InexactFloat64())
		r.f.SetCellStyle(sheet, makeAxis(col+1, curRow), makeAxis(col+14+corr, curRow), totalStyle)
	}
	r.f.MergeCell(sheet, makeAxis(col, curRow-len(totals)), makeAxis(col+1+corr, curRow))
	r.f.SetCellStyle(sheet, makeAxis(col, curRow-len(totals)), makeAxis(col+1+corr, curRow), flatTotalStyle)

	for _, service := range r.serviceGroups {
		serviceName := service.NameRu
		abonents := r.osv.GetAbonentsWithService(serviceName)
		if len(abonents) > 0 {
			shortName := serviceName
			if utf8.RuneCountInString(serviceName) > 30 {
				runes := []rune(shortName)
				shortName = string(runes[:30])
			}
			if err := r.fillService(r.f, r.osi, abonents, shortName, serviceName); err != nil {
				return err
			}
		}
	}
	r.fillDebts(r.f, r.osi, r.osv.Abonents, "Должники")
	r.f.SetActiveSheet(0)

	return nil
}

func (r *saldoXLS) fillService(f *excelize.File, osi *model.OSIResult, abonents []model.Abonent, sheetName, title string) error {
	sheet := sheetName
	// Создать новый лист
	f.NewSheet(sheet)
	const (
		top  = 1
		left = 0
	)

	if err := f.SetSheetViewOptions(sheet, 0, excelize.ZoomScale(145)); err != nil {
		return err
	}
	titleStyle, err := f.NewStyle(`{
		"font":{"bold":false},
		"alignment":{
			"horizontal": "center"
			}
		}`)
	if err != nil {
		return err
	}

	headStyle, err := f.NewStyle(`{
		"font":{"bold":false},
		"alignment":{
			"horizontal": "center",
			"vertical": "center",
			"wrap_text":true
			},
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}
	]
	}`)
	if err != nil {
		return err
	}

	flatStyle, err := f.NewStyle(`{
		"alignment":{
			"horizontal": "right"
		},
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}
		]}`)
	if err != nil {
		return err
	}

	rowStyle, err := f.NewStyle(`{
		"number_format": 4,
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}]
		}`)
	if err != nil {
		return err
	}

	totalStyle, err := f.NewStyle(`{
		"number_format": 4,
		"font":{"bold":true},
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}]
		}`)
	if err != nil {
		return err
	}

	f.SetCellValue(sheet, makeAxis(left, top), osi.OSIName())
	f.MergeCell(sheet, makeAxis(left, top), makeAxis(left+5, top))
	//f.SetCellStyle(sheet, makeAxis(left, top), makeAxis(left+5, top), headStyle)
	f.SetCellValue(sheet, makeAxis(left, top+1), fmt.Sprintf("Оборотно-сальдовая ведомость (%s) за период %s по %s", title, r.begin.Format("02/01/2006"), r.end.Format("02/01/2006")))

	f.MergeCell(sheet, makeAxis(left, top+1), makeAxis(left+4, top+1))
	f.SetCellStyle(sheet, makeAxis(left, top), makeAxis(left+4, top+1), titleStyle)

	f.SetCellValue(sheet, makeAxis(left, top+2), "№ квартиры/помещения")
	f.SetColWidth(sheet, ColIndexToLetters(left), ColIndexToLetters(left), 28)
	f.SetCellValue(sheet, makeAxis(left+1, top+2), "Долг на начало, тг.")
	f.SetCellValue(sheet, makeAxis(left+2, top+2), "Начислено, тг.")
	f.SetCellValue(sheet, makeAxis(left+3, top+2), "Корректировки, тг.")
	f.SetCellValue(sheet, makeAxis(left+4, top+2), "Оплачено, тг.")
	f.SetCellValue(sheet, makeAxis(left+5, top+2), "Пеня, тг.")
	f.SetCellValue(sheet, makeAxis(left+6, top+2), "Долг на конец, тг.")
	f.SetCellStyle(sheet, makeAxis(left, top+2), makeAxis(left+6, top+2), headStyle)
	f.SetColWidth(sheet, ColIndexToLetters(left+1), ColIndexToLetters(left+6), 18)

	col := left
	row := top + 3
	curRow := row
	for _, abonent := range abonents {
		f.SetCellValue(sheet, makeAxis(col, curRow), abonent.Flat+common.GetFlatType(abonent.AreaTypeCode))
		f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col, curRow), flatStyle)
		f.SetCellValue(sheet, makeAxis(col+1, curRow), abonent.Service.Begin)
		f.SetCellValue(sheet, makeAxis(col+2, curRow), abonent.Service.SumOfAccurals)
		f.SetCellValue(sheet, makeAxis(col+3, curRow), abonent.Service.SumOfFixes)
		f.SetCellValue(sheet, makeAxis(col+4, curRow), abonent.Service.Kredit)
		f.SetCellValue(sheet, makeAxis(col+5, curRow), abonent.Service.SumOfFines)
		f.SetCellValue(sheet, makeAxis(col+6, curRow), abonent.Service.End)
		f.SetCellStyle(sheet, makeAxis(col+1, curRow), makeAxis(col+6, curRow), rowStyle)
		curRow++
	}
	f.SetCellValue(sheet, makeAxis(col, curRow), "ИТОГО")
	if len(abonents) > 0 {
		f.SetCellFormula(sheet, makeAxis(col+1, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+1, row), makeAxis(col+1, row+len(abonents)-1)))
		f.SetCellFormula(sheet, makeAxis(col+2, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+2, row), makeAxis(col+2, row+len(abonents)-1)))
		f.SetCellFormula(sheet, makeAxis(col+3, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+3, row), makeAxis(col+3, row+len(abonents)-1)))
		f.SetCellFormula(sheet, makeAxis(col+4, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+4, row), makeAxis(col+4, row+len(abonents)-1)))
		f.SetCellFormula(sheet, makeAxis(col+5, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+5, row), makeAxis(col+5, row+len(abonents)-1)))
		f.SetCellFormula(sheet, makeAxis(col+6, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+6, row), makeAxis(col+6, row+len(abonents)-1)))
	}

	f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col+6, curRow), totalStyle)

	return nil
}

func (r *saldoXLS) fillDebts(f *excelize.File, osi *model.OSIResult, abonents []model.OSVAbonent, sheetName string) error {
	sheet := sheetName
	// Создать новый лист
	f.NewSheet(sheet)
	const (
		top  = 1
		left = 0
	)

	if err := f.SetSheetViewOptions(sheet, 0, excelize.ZoomScale(145)); err != nil {
		return err
	}
	titleStyle, err := f.NewStyle(`{
		"font":{"bold":false},
		"alignment":{
			"horizontal": "center"
			}
		}`)
	if err != nil {
		return err
	}

	headStyle, err := f.NewStyle(`{
		"font":{"bold":false},
		"alignment":{
			"horizontal": "center",
			"vertical": "center",
			"wrap_text":true
			},
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}
	]
	}`)
	if err != nil {
		return err
	}
	flatStyle, err := f.NewStyle(`{
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}]
		}`)
	if err != nil {
		return err
	}
	rowStyle, err := f.NewStyle(`{
		"number_format": 4,
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}]
		}`)
	if err != nil {
		return err
	}
	totalStyle, err := f.NewStyle(`{
		"font":{"bold":true},
		"number_format": 4,
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 1
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 1
		}]
		}`)
	if err != nil {
		return err
	}

	debts, services := getOsvDebts(abonents)
	colNum := len(services) + 1
	f.SetCellValue(sheet, makeAxis(left, top), osi.OSIName())
	f.MergeCell(sheet, makeAxis(left, top), makeAxis(left+colNum, top))
	f.SetCellValue(sheet, makeAxis(left, top+1), fmt.Sprintf("Задолженность на %s", r.end.Format("02/01/2006")))

	f.MergeCell(sheet, makeAxis(left, top+1), makeAxis(left+colNum, top+1))
	f.SetCellStyle(sheet, makeAxis(left, top), makeAxis(left+colNum, top+1), titleStyle)

	f.SetCellValue(sheet, makeAxis(left, top+2), "№ помещения")
	f.SetColWidth(sheet, ColIndexToLetters(left), ColIndexToLetters(left), 20)
	f.SetRowHeight(sheet, top+2, 30)
	for i, service := range services {
		f.SetCellValue(sheet, makeAxis(left+i+1, top+2), service)
	}
	f.SetColWidth(sheet, ColIndexToLetters(left+1), ColIndexToLetters(left+colNum-1), 30)
	f.SetCellValue(sheet, makeAxis(left+colNum, top+2), "Всего")
	f.SetColWidth(sheet, ColIndexToLetters(left+colNum), ColIndexToLetters(left+colNum), 30)
	f.SetCellStyle(sheet, makeAxis(left, top+2), makeAxis(left+colNum, top+2), headStyle)

	col := left
	row := top + 3
	curRow := row
	for _, debt := range debts {
		f.SetCellValue(sheet, makeAxis(col, curRow), debt.Flat+common.GetFlatType(debt.AreaTypeCode))
		for i, service := range services {
			dolg, ok := debt.Services[service]
			if !ok {
				dolg = 0
			}
			f.SetCellValue(sheet, makeAxis(col+i+1, curRow), dolg)
		}
		f.SetCellFormula(sheet, makeAxis(col+colNum, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+1, curRow), makeAxis(col+len(services), curRow)))
		f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col, curRow), flatStyle)
		f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col+colNum, curRow), rowStyle)
		curRow++
	}
	f.SetCellValue(sheet, makeAxis(col, curRow), "ИТОГО")
	if len(debts) > 0 {
		for i, _ := range services {
			f.SetCellFormula(sheet, makeAxis(col+i+1, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+i+1, row), makeAxis(col+i+1, row+len(debts)-1)))
		}
		f.SetCellFormula(sheet, makeAxis(col+colNum, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+colNum, row), makeAxis(col+colNum, row+len(debts)-1)))
		f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col+colNum, curRow), totalStyle)
	}

	return nil
}

func containsString(a []string, v string) bool {
	for _, item := range a {
		if item == v {
			return true
		}
	}
	return false
}

func getOsvDebts(abonents []model.OSVAbonent) ([]osvDebt, []string) {
	services := make([]string, 0, 3)
	debts := make([]osvDebt, 0, len(abonents))
	for _, abonent := range abonents {
		debt := osvDebt{Flat: abonent.Flat, AreaTypeCode: abonent.AreaTypeCode, Services: map[string]float64{}}
		for _, service := range abonent.Services {
			if service.End > 0 {
				debt.Services[service.ServiceName] = service.End
				if !containsString(services, service.ServiceName) {
					services = append(services, service.ServiceName)
				}
			}
		}
		if len(debt.Services) > 0 {
			debts = append(debts, debt)
		}
	}
	return debts, services
}

func (r *saldoXLS) WriteTo(w io.Writer) (int64, error) {
	return r.f.WriteTo(w)
}
