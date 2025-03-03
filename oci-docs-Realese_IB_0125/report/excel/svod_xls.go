package excel

import (
	"fmt"
	"io"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/report"
)

type svodXLS struct {
	f     *excelize.File
	osi   *model.OSIResult
	svod  []model.PaymentOrder
	begin time.Time
	end   time.Time
}

func NewSvodXLS(osi *model.OSIResult, svod []model.PaymentOrder, begin, end time.Time) report.Report {
	return &svodXLS{
		f:     excelize.NewFile(),
		osi:   osi,
		svod:  svod,
		begin: begin,
		end:   end,
	}
}

func (r *svodXLS) Render() error {
	const (
		top  = 1
		left = 0
	)

	sheet := "Свод"
	// Создать новый лист
	index := r.f.NewSheet(sheet)
	r.f.SetActiveSheet(index)

	if err := r.f.SetSheetViewOptions(sheet, 0, excelize.ZoomScale(100)); err != nil {
		return err
	}

	titleStyle, err := r.f.NewStyle(`{
		"font":{"bold":true},
		"alignment":{
			"horizontal": "center"
			}
		}`)
	if err != nil {
		return err
	}
	headStyle, err := r.f.NewStyle(`{
		"font":{"bold":true},
		"alignment":{
			"horizontal": "center",
			"vertical": "center",
			"wrap_text":true
			},
		"border": [
		{
			"type": "left",
			"color": "#0D0101",
			"style": 2
		},
		{
			"type": "top",
			"color": "#0D0101",
			"style": 2
		},
		{
			"type": "bottom",
			"color": "#0D0101",
			"style": 2
		},
		{
			"type": "right",
			"color": "#0D0101",
			"style": 2
		}
	]
	}`)
	if err != nil {
		return err
	}

	rowStyle, err := r.f.NewStyle(`{
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

	// osiName := r.osi.Name
	// if !strings.HasPrefix(osiName, r.osi.UnionTypeRu) {
	// 	osiName = fmt.Sprintf(`%s "%s"`, r.osi.UnionTypeRu, osiName)
	// } else {
	// 	osiName = fmt.Sprintf(`"%s"`, osiName)
	// }
	r.f.SetCellValue(sheet, makeAxis(left, top), fmt.Sprintf(`Свод по принятым платежам %s за период c %s по %s`, r.osi.OSIName(), r.begin.Format("02.01.2006"), r.end.Format("02.01.2006")))
	r.f.MergeCell(sheet, makeAxis(left, top), makeAxis(left+7, top))
	r.f.SetCellStyle(sheet, makeAxis(left, top), makeAxis(left+7, top), titleStyle)

	r.f.SetCellValue(sheet, makeAxis(left, top+1), r.osi.Address)
	r.f.MergeCell(sheet, makeAxis(left, top+1), makeAxis(left+7, top+1))
	r.f.SetCellStyle(sheet, makeAxis(left, top+1), makeAxis(left+7, top+1), titleStyle)

	col := left
	row := top + 3
	r.f.SetCellValue(sheet, makeAxis(col, row), "№ п/п")
	r.f.SetColWidth(sheet, ColIndexToLetters(col), ColIndexToLetters(col), 20)

	r.f.SetCellValue(sheet, makeAxis(col+1, row), "Источник")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+1), ColIndexToLetters(col+1), 30)

	r.f.SetCellValue(sheet, makeAxis(col+2, row), "Счет")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+2), ColIndexToLetters(col+2), 30)

	r.f.SetCellValue(sheet, makeAxis(col+3, row), "Дата")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+3), ColIndexToLetters(col+3), 20)

	r.f.SetCellValue(sheet, makeAxis(col+4, row), "Сумма платежей")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+4), ColIndexToLetters(col+4), 30)

	r.f.SetCellValue(sheet, makeAxis(col+5, row), "Комиссия Банка")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+5), ColIndexToLetters(col+5), 20)

	r.f.SetCellValue(sheet, makeAxis(col+6, row), "Комиссия eosi.kz")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+6), ColIndexToLetters(col+6), 20)

	r.f.SetCellValue(sheet, makeAxis(col+7, row), "К перечислению")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+7), ColIndexToLetters(col+7), 20)

	r.f.SetRowHeight(sheet, row, 20)
	r.f.SetCellStyle(sheet, makeAxis(col, row), makeAxis(col+7, row), headStyle)

	row++
	curRow := row
	for i, s := range r.svod {
		r.f.SetCellValue(sheet, makeAxis(col, curRow), i+1)
		r.f.SetCellValue(sheet, makeAxis(col+1, curRow), s.BankName)
		r.f.SetCellValue(sheet, makeAxis(col+2, curRow), s.IBAN)
		r.f.SetCellValue(sheet, makeAxis(col+3, curRow), s.Date.String())
		r.f.SetCellValue(sheet, makeAxis(col+4, curRow), *s.Amount)
		r.f.SetCellValue(sheet, makeAxis(col+5, curRow), *s.ComisBank)
		r.f.SetCellValue(sheet, makeAxis(col+6, curRow), *s.ComisOur)
		r.f.SetCellValue(sheet, makeAxis(col+7, curRow), *s.AmountToTransfer)
		r.f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col+7, curRow), rowStyle)
		curRow++
	}
	svodSize := len(r.svod)
	if svodSize > 0 {
		r.f.SetCellFormula(sheet, makeAxis(col+4, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+4, row), makeAxis(col+4, row+svodSize-1)))
		r.f.SetCellFormula(sheet, makeAxis(col+5, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+5, row), makeAxis(col+5, row+svodSize-1)))
		r.f.SetCellFormula(sheet, makeAxis(col+6, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+6, row), makeAxis(col+6, row+svodSize-1)))
		r.f.SetCellFormula(sheet, makeAxis(col+7, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+7, row), makeAxis(col+7, row+svodSize-1)))
		r.f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col+7, curRow), totalStyle)
	}
	return nil
}

func (r *svodXLS) WriteTo(w io.Writer) (int64, error) {
	return r.f.WriteTo(w)
}
