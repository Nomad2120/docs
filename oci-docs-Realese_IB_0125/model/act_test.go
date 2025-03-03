package model

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetActResult(t *testing.T) {
	data := `{
    "id": 3,
    "createDt": "2021-08-19T18:34:11.118274",
    "signDt": null,
    "actPeriod": "2021-08-31T00:00:00",
    "actNum": "0000000090",
    "stateCode": "CREATED",
    "stateName": "Создан",
    "osiId": 68,
    "osiName": "\"12микрорайон53\" ",
    "osiIdn": "210340019400",
    "apartCount": 117,
    "planAccuralId": 90,
    "amount": 11228.03
  }`

	var res GetActResult
	err := json.Unmarshal([]byte(data), &res)
	require.NoError(t, err)
	fmt.Printf("%+v", res)
}
