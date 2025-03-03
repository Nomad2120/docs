package excel

import (
	"fmt"
	"io"

	"github.com/360EntSecGroup-Skylar/excelize"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/common"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/report"
)

type abonentsXLS struct {
	f        *excelize.File
	osi      *model.OSIResult
	abonents []model.OSIAbonent
}

func NewAbonentsXLS(osi *model.OSIResult, abonents []model.OSIAbonent) report.Report {
	return &abonentsXLS{
		f:        excelize.NewFile(),
		osi:      osi,
		abonents: abonents,
	}
}

func (r *abonentsXLS) Render() error {
	const (
		top  = 1
		left = 0
	)

	sheet := "Помещения/квартиры"
	// Создать новый лист
	index := r.f.NewSheet(sheet)
	r.f.SetActiveSheet(index)

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

	// osiName := r.osi.Name
	// if !strings.HasPrefix(osiName, r.osi.UnionTypeRu) {
	// 	osiName = fmt.Sprintf(`%s "%s"`, r.osi.UnionTypeRu, osiName)
	// } else {
	// 	osiName = fmt.Sprintf(`"%s"`, osiName)
	// }
	r.f.SetCellValue(sheet, makeAxis(left, top), r.osi.OSIName())
	r.f.MergeCell(sheet, makeAxis(left, top), makeAxis(left+12, top))

	r.f.SetCellValue(sheet, makeAxis(left, top+1), "Адрес: "+r.osi.Address)
	r.f.MergeCell(sheet, makeAxis(left, top+1), makeAxis(left+12, top+1))
	r.f.SetCellStyle(sheet, makeAxis(left, top+1), makeAxis(left+12, top+1), titleStyle)

	r.f.SetCellValue(sheet, makeAxis(left, top+2), "Помещения/квартиры")
	r.f.MergeCell(sheet, makeAxis(left, top+2), makeAxis(left+12, top+2))
	r.f.SetCellStyle(sheet, makeAxis(left, top), makeAxis(left, top+2), titleStyle)

	col := left
	row := top + 3

	r.f.SetCellValue(sheet, makeAxis(col, row), "Номер")
	r.f.SetColWidth(sheet, ColIndexToLetters(col), ColIndexToLetters(col), 8)
	r.f.SetCellValue(sheet, makeAxis(col+1, row), "Собственник")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+1), ColIndexToLetters(col+1), 30)
	r.f.SetCellValue(sheet, makeAxis(col+2, row), "Владелец")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+2), ColIndexToLetters(col+2), 15)
	r.f.SetCellValue(sheet, makeAxis(col+3, row), "ИИН")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+3), ColIndexToLetters(col+3), 15)
	r.f.SetCellValue(sheet, makeAxis(col+4, row), "Телефон")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+4), ColIndexToLetters(col+4), 16)
	r.f.SetCellValue(sheet, makeAxis(col+5, row), "Тип")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+5), ColIndexToLetters(col+5), 10)
	r.f.SetCellValue(sheet, makeAxis(col+6, row), "ЛС ОСИ")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+6), ColIndexToLetters(col+6), 10)
	r.f.SetCellValue(sheet, makeAxis(col+7, row), "ЛС ЕРЦ")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+7), ColIndexToLetters(col+7), 10)
	r.f.SetCellValue(sheet, makeAxis(col+8, row), "Этаж")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+8), ColIndexToLetters(col+8), 8)
	r.f.SetCellValue(sheet, makeAxis(col+9, row), "Площадь")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+9), ColIndexToLetters(col+9), 10)
	r.f.SetCellValue(sheet, makeAxis(col+10, row), "Полезная площадь")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+10), ColIndexToLetters(col+10), 10)
	r.f.SetCellValue(sheet, makeAxis(col+11, row), "Проживает")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+11), ColIndexToLetters(col+11), 12)
	r.f.SetCellValue(sheet, makeAxis(col+12, row), "Зарегистрировано")
	r.f.SetColWidth(sheet, ColIndexToLetters(col+12), ColIndexToLetters(col+12), 10)

	r.f.SetCellStyle(sheet, makeAxis(col, row), makeAxis(col+12, row), headStyle)

	curRow := row + 1
	for _, abonent := range r.abonents {
		r.f.SetCellValue(sheet, makeAxis(col, curRow), abonent.Flat)
		r.f.SetCellValue(sheet, makeAxis(col+1, curRow), abonent.Name)
		r.f.SetCellValue(sheet, makeAxis(col+2, curRow), abonent.Owner)
		r.f.SetCellValue(sheet, makeAxis(col+3, curRow), abonent.IDN)
		r.f.SetCellValue(sheet, makeAxis(col+4, curRow), common.FormatPhone(abonent.Phone))
		r.f.SetCellValue(sheet, makeAxis(col+5, curRow), abonent.AreaTypeNameRu)
		r.f.SetCellValue(sheet, makeAxis(col+6, curRow), abonent.ID)
		r.f.SetCellValue(sheet, makeAxis(col+7, curRow), abonent.ErcAccount)
		r.f.SetCellValue(sheet, makeAxis(col+8, curRow), abonent.Floor)
		r.f.SetCellValue(sheet, makeAxis(col+9, curRow), abonent.Square)
		if abonent.EffectiveSquare != nil {
			r.f.SetCellValue(sheet, makeAxis(col+10, curRow), *abonent.EffectiveSquare)
		}
		r.f.SetCellValue(sheet, makeAxis(col+11, curRow), abonent.LivingFact)
		r.f.SetCellValue(sheet, makeAxis(col+12, curRow), abonent.LivingJur)
		r.f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col+12, curRow), rowStyle)
		curRow++
	}
	r.f.SetCellValue(sheet, makeAxis(col, curRow), "ИТОГО")
	r.f.SetCellValue(sheet, makeAxis(col+1, curRow), len(r.abonents))
	r.f.SetCellFormula(sheet, makeAxis(col+9, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+9, row+1), makeAxis(col+9, curRow-1)))
	r.f.SetCellFormula(sheet, makeAxis(col+10, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+10, row+1), makeAxis(col+10, curRow-1)))
	r.f.SetCellFormula(sheet, makeAxis(col+11, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+11, row+1), makeAxis(col+11, curRow-1)))
	r.f.SetCellFormula(sheet, makeAxis(col+12, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+12, row+1), makeAxis(col+12, curRow-1)))
	r.f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col+12, curRow), totalStyle)
	return nil
}

func (r *abonentsXLS) WriteTo(w io.Writer) (int64, error) {
	return r.f.WriteTo(w)
}
