package model

import (
	"testing"
	"time"
)

func TestGetDateBegin(t *testing.T) {
	ar := &AccountReportsResult{
		Period: ISODate(time.Now()),
	}
	s := ar.GetReportDateKZ()
	println(s)
}

