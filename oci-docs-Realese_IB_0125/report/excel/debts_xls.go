package excel

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/shopspring/decimal"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/common"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/model"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/report"
)

type DebtInfo struct {
	Debt float64
	Fine float64
}

type flatDebt struct {
	Flat         string
	AreaTypeCode string
	Services     map[string]DebtInfo
}

type debtsXLS struct {
	f             *excelize.File
	osi           *model.OSIResult
	osv           *model.OSVResult
	serviceGroups []model.ServiceGroup
	begin         time.Time
	end           time.Time
}

func NewDebtsXLS(osi *model.OSIResult, osv *model.OSVResult, serviceGroups []model.ServiceGroup, begin, end time.Time) report.Report {
	return &debtsXLS{
		f:             excelize.NewFile(),
		osi:           osi,
		osv:           osv,
		serviceGroups: serviceGroups,
		begin:         begin,
		end:           end,
	}
}

func (r *debtsXLS) Render() error {
	sort.SliceStable(r.osv.Abonents, func(i, j int) bool {
		return len(r.osv.Abonents[i].Flat) < len(r.osv.Abonents[j].Flat) ||
			(len(r.osv.Abonents[i].Flat) == len(r.osv.Abonents[j].Flat) && r.osv.Abonents[i].Flat < r.osv.Abonents[j].Flat)
	})

	if err := r.fillDebts(r.f, r.osi, r.osv.Abonents, "Должники"); err != nil {
		return err
	}
	r.f.DeleteSheet("Sheet1")
	return nil
}

func (r *debtsXLS) WriteTo(w io.Writer) (int64, error) {
	return r.f.WriteTo(w)
}

func (r *debtsXLS) getDebts(abonents []model.OSVAbonent) ([]osvDebt, []string) {
	services := make([]string, 0, 3)
	debts := make([]osvDebt, 0, len(abonents))
	//today := time.Now()
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

func (r *debtsXLS) fillDebts(f *excelize.File, osi *model.OSIResult, abonents []model.OSVAbonent, sheetName string) error {
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

	debts, services := getDebts(abonents)
	colNum := len(services) + 1
	f.SetCellValue(sheet, makeAxis(left, top), osi.OSIName())
	f.MergeCell(sheet, makeAxis(left, top), makeAxis(left+colNum+1, top))
	f.SetCellValue(sheet, makeAxis(left, top+1), fmt.Sprintf("Задолженность на %s", time.Now().Format("02/01/2006")))

	f.MergeCell(sheet, makeAxis(left, top+1), makeAxis(left+colNum+1, top+1))
	f.SetCellStyle(sheet, makeAxis(left, top), makeAxis(left+colNum+1, top+1), titleStyle)

	f.SetCellValue(sheet, makeAxis(left, top+2), "№ помещения")
	f.SetColWidth(sheet, ColIndexToLetters(left), ColIndexToLetters(left), 20)
	for i, service := range services {
		f.SetCellValue(sheet, makeAxis(left+i+1, top+2), service)
	}
	f.SetColWidth(sheet, ColIndexToLetters(left+1), ColIndexToLetters(left+colNum-1), 30)

	f.SetCellValue(sheet, makeAxis(left+colNum, top+2), "Пеня")
	f.SetColWidth(sheet, ColIndexToLetters(left+colNum), ColIndexToLetters(left+colNum), 30)

	f.SetCellValue(sheet, makeAxis(left+colNum+1, top+2), "Всего")
	f.SetColWidth(sheet, ColIndexToLetters(left+colNum+1), ColIndexToLetters(left+colNum+1), 30)
	f.SetCellStyle(sheet, makeAxis(left, top+2), makeAxis(left+colNum+1, top+2), headStyle)
	f.SetRowHeight(sheet, top+2, 40)

	col := left
	row := top + 3
	curRow := row
	for _, debt := range debts {
		f.SetCellValue(sheet, makeAxis(col, curRow), debt.Flat+common.GetFlatType(debt.AreaTypeCode))
		fine := decimal.Zero
		for i, service := range services {
			dolg, ok := debt.Services[service]
			if !ok {
				dolg = DebtInfo{}
			}
			f.SetCellValue(sheet, makeAxis(col+i+1, curRow), dolg.Debt)
			fine = fine.Add(decimal.NewFromFloat(dolg.Fine))
		}
		f.SetCellValue(sheet, makeAxis(col+colNum, curRow), fine.InexactFloat64())
		f.SetCellFormula(sheet, makeAxis(col+colNum+1, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+1, curRow), makeAxis(col+len(services)+1, curRow)))
		f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col, curRow), flatStyle)
		f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col+colNum+1, curRow), rowStyle)
		curRow++
	}
	f.SetCellValue(sheet, makeAxis(col, curRow), "ИТОГО")
	if len(debts) > 0 {
		for i, _ := range services {
			f.SetCellFormula(sheet, makeAxis(col+i+1, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+i+1, row), makeAxis(col+i+1, row+len(debts)-1)))
		}
		f.SetCellFormula(sheet, makeAxis(col+colNum, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+colNum, row), makeAxis(col+colNum, row+len(debts)-1)))
		f.SetCellFormula(sheet, makeAxis(col+colNum+1, curRow), fmt.Sprintf("SUM(%s:%s)", makeAxis(col+colNum+1, row), makeAxis(col+colNum+1, row+len(debts)-1)))
		f.SetCellStyle(sheet, makeAxis(col, curRow), makeAxis(col+colNum+1, curRow), totalStyle)
	}

	return nil
}

func getDebts(abonents []model.OSVAbonent) ([]flatDebt, []string) {
	services := make([]string, 0, 3)
	debts := make([]flatDebt, 0, len(abonents))
	for _, abonent := range abonents {
		debt := flatDebt{Flat: abonent.Flat, AreaTypeCode: abonent.AreaTypeCode, Services: make(map[string]DebtInfo)}
		for _, service := range abonent.Services {
			if service.End > 0 {
				d := decimal.NewFromFloat(service.End).Sub(decimal.NewFromFloat(service.SumOfFines))
				debt.Services[service.ServiceName] = DebtInfo{
					Fine: service.SumOfFines,
					Debt: d.InexactFloat64(),
				}
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
