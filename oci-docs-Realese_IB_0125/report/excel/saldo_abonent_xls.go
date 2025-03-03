package excel

import (
	"fmt"
	"io"

	"github.com/360EntSecGroup-Skylar/excelize"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/report"
)

type saldoAbonentXLS struct {
	f           *excelize.File
	osi         *model.OSIResult
	osv         []model.AbonentOSVResult
	flat        string
	abonentInfo *model.OSIAbonent
}

func NewSaldoAbonentXLS(osi *model.OSIResult, osv []model.AbonentOSVResult, abonentInfo *model.OSIAbonent, flat string) report.Report {
	return &saldoAbonentXLS{
		f:           excelize.NewFile(),
		osi:         osi,
		osv:         osv,
		flat:        flat,
		abonentInfo: abonentInfo,
	}
}

func (r *saldoAbonentXLS) Render() error {
	sheet := "Сальдо по помещению"
	// Создать новый лист
	r.f.NewSheet(sheet)

	r.f.DeleteSheet("Sheet1")
	const (
		top  = 1
		left = 0
	)
	if err := r.f.SetSheetViewOptions(sheet, 0, excelize.ZoomScale(145)); err != nil {
		return err
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
	fixStyle, err := r.f.NewStyle(`{"font":{"italic":false},
	"fill":{"type":"pattern","color":["#ffc000"],"pattern":1},	
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
	fixesStyle, err := r.f.NewStyle(`{"font":{"italic":false},
	"alignment":{
		"horizontal": "left",
		"vertical": "top",
		"wrap_text":true
		},
	"fill":{"type":"pattern","color":["#ffc000"],"pattern":1},	
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
	periodStyle, err := r.f.NewStyle(`{
		"alignment":{
			"horizontal": "center",
			"vertical": "center"
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
	totalStyle, err := r.f.NewStyle(`{"font":{"italic":false, "bold": true},
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

	r.f.SetCellValue(sheet, makeAxis(left, top), r.osi.OSIName())
	r.f.MergeCell(sheet, makeAxis(left, top), makeAxis(left+8, top))
	r.f.SetCellValue(sheet, makeAxis(left, top+1), "адрес "+r.osi.Address)
	r.f.MergeCell(sheet, makeAxis(left, top+1), makeAxis(left+8, top+1))
	r.f.SetCellValue(sheet, makeAxis(left, top+2), fmt.Sprintf("Оборотно-сальдовая ведомость по квартире №%s собственник %s", r.flat, r.abonentInfo.Name))
	r.f.MergeCell(sheet, makeAxis(left, top+2), makeAxis(left+8, top+2))
	r.f.SetCellStyle(sheet, makeAxis(left, top), makeAxis(left, top+2), titleStyle)

	col := left
	row := top + 3

	r.f.SetCellValue(sheet, makeAxis(col, row), "Период")
	r.f.SetColWidth(sheet, ColIndexToLetters(col), ColIndexToLetters(col), 15)
	r.f.SetCellValue(sheet, makeAxis(col+1, row), "Услуга")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+1), ColIndexToLetters(col+1), 40)
	r.f.SetCellValue(sheet, makeAxis(col+2, row), "Сальдо на начало, тг.")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+2), ColIndexToLetters(col+2), 15)
	r.f.SetCellValue(sheet, makeAxis(col+3, row), "Начислено за период, тг.")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+3), ColIndexToLetters(col+3), 15)
	r.f.SetCellValue(sheet, makeAxis(col+4, row), "Корректировки, тг.")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+4), ColIndexToLetters(col+4), 15)
	r.f.SetCellValue(sheet, makeAxis(col+5, row), "Оплачено за период, тг.")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+5), ColIndexToLetters(col+5), 15)
	r.f.SetCellValue(sheet, makeAxis(col+6, row), "Пеня, тг.")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+6), ColIndexToLetters(col+6), 15)
	r.f.SetCellValue(sheet, makeAxis(col+7, row), "Долг, тг.")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+7), ColIndexToLetters(col+7), 15)
	r.f.SetRowHeight(sheet, row, 40)

	r.f.SetCellStyle(sheet, makeAxis(col, row), makeAxis(col+7, row), headStyle)

	curRow := row + 1
	bb, be, eb, ee := curRow, curRow, 0, 0
	for i, period := range r.osv {
		r.f.SetCellValue(sheet, makeAxis(col, curRow), period.Period)
		begin := curRow
		rows := len(period.Services)
		fixesNum := 0
		for j, service := range period.Services {
			if i == 0 && j == 0 {
				bb = curRow
			}
			if i == 0 && j == len(period.Services)+len(service.Fixes)-1 {
				be = curRow
			}
			if i == len(r.osv)-1 && j == 0 {
				eb = curRow
			}
			if i == len(r.osv)-1 && j == len(period.Services)+len(service.Fixes)-1 {
				ee = curRow
			}
			r.f.SetCellValue(sheet, makeAxis(col+1, curRow), service.ServiceName)

			r.f.SetCellValue(sheet, makeAxis(col+2, curRow), service.Begin)
			r.f.SetCellValue(sheet, makeAxis(col+3, curRow), service.SumOfAccurals)
			r.f.SetCellValue(sheet, makeAxis(col+4, curRow), service.SumOfFixes)
			r.f.SetCellValue(sheet, makeAxis(col+5, curRow), service.Kredit)
			r.f.SetCellValue(sheet, makeAxis(col+6, curRow), service.SumOfFines)
			r.f.SetCellValue(sheet, makeAxis(col+7, curRow), service.End)
			r.f.SetCellStyle(sheet, makeAxis(col+1, curRow), makeAxis(col+7, curRow), serviceStyle)
			if service.SumOfFixes > 0 {
				r.f.SetCellStyle(sheet, makeAxis(col+4, curRow), makeAxis(col+4, curRow), fixStyle)
			}
			curRow++

			if service.SumOfFixes != 0 {
				fixesNum += len(service.Fixes)

				for k, fix := range service.Fixes {
					text := fmt.Sprintf("%d. %s (%.2f)", k+1, fix.Reason, fix.Amount)
					r.f.SetCellValue(sheet, makeAxis(col+4, curRow), text)
					r.f.MergeCell(sheet, makeAxis(col+4, curRow), makeAxis(col+7, curRow))
					r.f.SetCellStyle(sheet, makeAxis(col+4, curRow), makeAxis(col+7, curRow), fixesStyle)
					if len(text) > 45 {
						r.f.SetRowHeight(sheet, curRow, 24)
					}
					curRow++
				}
			}
		}
		if rows+fixesNum > 1 {
			r.f.MergeCell(sheet, makeAxis(col, begin), makeAxis(col, begin+rows+fixesNum-1))
			r.f.SetCellStyle(sheet, makeAxis(col, begin), makeAxis(col, begin+rows+fixesNum-1), periodStyle)
		} else {
			r.f.SetCellStyle(sheet, makeAxis(col, begin), makeAxis(col, begin), periodStyle)
		}
	}

	r.f.SetCellValue(sheet, makeAxis(col, curRow), "ИТОГО")
	r.f.SetCellFormula(sheet, makeAxis(col+2, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+2, bb), makeAxis(col+2, be)))
	r.f.SetCellFormula(sheet, makeAxis(col+3, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+3, row+1), makeAxis(col+3, curRow-1)))
	r.f.SetCellFormula(sheet, makeAxis(col+4, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+4, row+1), makeAxis(col+4, curRow-1)))
	r.f.SetCellFormula(sheet, makeAxis(col+5, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+5, row+1), makeAxis(col+5, curRow-1)))
	r.f.SetCellFormula(sheet, makeAxis(col+6, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+6, row+1), makeAxis(col+6, curRow-1)))
	r.f.SetCellFormula(sheet, makeAxis(col+7, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+7, eb), makeAxis(col+7, ee)))
	r.f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col+7, curRow), totalStyle)

	curRow += 2
	r.f.SetCellValue(sheet, makeAxis(col+1, curRow), fmt.Sprintf("Председатель %s", r.osi.OSIName()))
	r.f.SetCellValue(sheet, makeAxis(col+3, curRow), r.osi.Fio)

	return nil
}

func (r *saldoAbonentXLS) WriteTo(w io.Writer) (int64, error) {
	return r.f.WriteTo(w)
}
