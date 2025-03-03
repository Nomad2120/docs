package common

import "testing"

var samples = []struct {
	amount   float64
	upper    bool
	expected string
}{
	{1, true, "Один тенге 00 тиын"},
	{100.21, false, "сто тенге 21 тиын"},
	{3457890.00, true, "Три миллиона четыреста пятьдесят семь тысяч восемьсот девяносто тенге 00 тиын"},
}

func Test_Num2Str(t *testing.T) {
	for _, tt := range samples {
		res := Num2Str(tt.amount, tt.upper)
		if res != tt.expected {
			t.Errorf("Num2Str(%.2f): expected '%s', got '%s'", tt.amount, tt.expected, res)
		}
	}
}

func Test_Num2StrKaz(t *testing.T) {
	println(Num2StrKaz(100.21, false))
	println(Num2StrKaz(3457890.00, false))
}
