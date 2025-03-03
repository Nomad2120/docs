package model

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebtorNotificationResponse(t *testing.T) {
	data := `{
		"result": {
		  "flat": "434",
		  "address": "г.Актобе, ул.Бокенбай батыра, дом 155, корпус 7, кв. 434",
		  "osiName": "ОСИ \"ЖК Юнис Сити\" корпус 7",
		  "osiChairman": "Ахметов Арман Ахметұлы",
		  "debtDate": "2022-02-28T00:00:00",
		  "servicesDebts": [
			{
			  "serviceName": "Взнос на капитальный ремонт здания",
			  "saldo": 2118,
			  "saldoString": "две тысячи сто восемнадцать тенге 00 тиын"
			},
			{
			  "serviceName": "Взнос на содержание общего имущества",
			  "saldo": 4356,
			  "saldoString": "четыре тысячи триста пятьдесят шесть тенге 00 тиын"
			}
		  ]
		},
		"code": 0,
		"message": "Success"
	  }`

	var resp DebtorNotificationResponse

	err := json.Unmarshal([]byte(data), &resp)
	require.NoError(t, err)
	fmt.Printf("%+v\n", resp.Result.DebtDate)
}

func TestGetDebtDate(t *testing.T) {
	var na NotaryApplicationResult
	dt := na.GetDebtDate()
	assert.Equal(t, "25", dt[:2])

}
