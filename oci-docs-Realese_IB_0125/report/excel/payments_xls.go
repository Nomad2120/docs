package excel

import (
	"fmt"
	"io"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/report"
)

type paymentsXLS struct {
	f        *excelize.File
	osi      *model.OSIResult
	payments []model.PaymentInfo
	begin    time.Time
	end      time.Time
}

func NewPaymentsXLS(osi *model.OSIResult, payments []model.PaymentInfo, begin, end time.Time) report.Report {
	return &paymentsXLS{
		f:        excelize.NewFile(),
		osi:      osi,
		payments: payments,
		begin:    begin,
		end:      end,
	}
}

func (r *paymentsXLS) Render() error {
	const (
		top  = 1
		left = 0
	)

	sheet := "Sheet1"
	// Создать новый лист
	index := r.f.NewSheet(sheet)
	r.f.SetActiveSheet(index)

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
		"font":{"color": "#919eab"},	
		"fill":{"type":"pattern","color":["#212b36"],"pattern":1},
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

	rowStyle, err := r.f.NewStyle(`{"alignment":{
		"vertical": "center",
		"wrap_text":true
		}}`)
	if err != nil {
		return err
	}

	r.f.SetCellValue(sheet, makeAxis(left+1, top), r.osi.OSIName())
	r.f.MergeCell(sheet, makeAxis(left+1, top), makeAxis(left+4, top))
	r.f.SetCellValue(sheet, makeAxis(left+1, top+1), fmt.Sprintf("за период %s по %s", r.begin.Format("02/01/2006"), r.end.Format("02/01/2006")))
	r.f.MergeCell(sheet, makeAxis(left+1, top+1), makeAxis(left+4, top+1))
	r.f.SetCellStyle(sheet, makeAxis(left+1, top), makeAxis(left+4, top+1), titleStyle)

	r.f.SetCellValue(sheet, makeAxis(left+1, top+2), "адрес: "+r.osi.Address)
	r.f.MergeCell(sheet, makeAxis(left+1, top+2), makeAxis(left+4, top+2))
	r.f.SetCellStyle(sheet, makeAxis(left+1, top+2), makeAxis(left+4, top+2), titleStyle)

	col := left
	row := top + 6
	r.f.SetCellValue(sheet, makeAxis(col, row), "Дата платежа")
	r.f.SetColWidth(sheet, ColIndexToLetters(col), ColIndexToLetters(col), 20)

	r.f.SetCellValue(sheet, makeAxis(col+1, row), "Принят")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+1), ColIndexToLetters(col+1), 20)

	r.f.SetCellValue(sheet, makeAxis(col+2, row), "Собственник помещения")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+2), ColIndexToLetters(col+2), 30)

	r.f.SetCellValue(sheet, makeAxis(col+3, row), "Номер помещения")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+3), ColIndexToLetters(col+3), 20)

	r.f.SetCellValue(sheet, makeAxis(col+4, row), "Услуга")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+4), ColIndexToLetters(col+4), 30)

	r.f.SetCellValue(sheet, makeAxis(col+5, row), "Сумма (тенге)")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+5), ColIndexToLetters(col+5), 15)

	r.f.SetRowHeight(sheet, row, 30)
	r.f.SetCellStyle(sheet, makeAxis(col, row), makeAxis(col+5, row), headStyle)

	curRow := row + 1
	for _, pay := range r.payments {
		r.f.SetCellValue(sheet, makeAxis(col, curRow), pay.Dt)
		r.f.SetCellValue(sheet, makeAxis(col+1, curRow), pay.BankName)
		r.f.SetCellValue(sheet, makeAxis(col+2, curRow), pay.AbonentName)
		r.f.SetCellValue(sheet, makeAxis(col+3, curRow), pay.Flat)
		r.f.SetCellValue(sheet, makeAxis(col+4, curRow), pay.ServiceName)
		r.f.SetCellValue(sheet, makeAxis(col+5, curRow), pay.Amount)
		r.f.SetRowHeight(sheet, curRow, 40)
		r.f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col+5, curRow), rowStyle)
		curRow++
	}
	return nil
}

func (r *paymentsXLS) WriteTo(w io.Writer) (int64, error) {
	return r.f.WriteTo(w)
}
