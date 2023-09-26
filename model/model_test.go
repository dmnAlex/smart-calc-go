package model

import (
	"fmt"
	"math"
	"testing"
)

var errMsg = map[Err_t]string{
	ERR_MALFORMED:    "Malformed expression",
	ERR_DIVBYZERO:    "Division by zero",
	ERR_MISSBRACKET:  "Mismatched brackets",
	ERR_NEGATIVEROOT: "Negative root",
	ERR_SUCCESS:      "Success",
}

func TestCorrectEquations(t *testing.T) {
	c := NewCalcModel()
	defer DeleteCalcModel(c)

	var tests = []struct {
		eq     string
		x, res float64
		err    Err_t
	}{
		{eq: "2 + 2", x: 0, res: 4, err: ERR_SUCCESS},
		{eq: "2 - 2", x: 0, res: 0, err: ERR_SUCCESS},
		{eq: "2 * 2", x: 0, res: 4, err: ERR_SUCCESS},
		{eq: "2 / 2", x: 0, res: 1, err: ERR_SUCCESS},
		{eq: "2 mod 2", x: 0, res: 0, err: ERR_SUCCESS},
		{eq: "2 ^ 2", x: 0, res: 4, err: ERR_SUCCESS},
		{eq: "-2 + 2", x: 0, res: 0, err: ERR_SUCCESS},
		{eq: "-2 - 2", x: 0, res: -4, err: ERR_SUCCESS},
		{eq: "-2 * 2", x: 0, res: -4, err: ERR_SUCCESS},
		{eq: "-2 / 2", x: 0, res: -1, err: ERR_SUCCESS},
		{eq: "-2 mod 2", x: 0, res: 0, err: ERR_SUCCESS},
		{eq: "-2 ^ 2", x: 0, res: -4, err: ERR_SUCCESS},
		{eq: "2 + (-2)", x: 0, res: 0, err: ERR_SUCCESS},
		{eq: "2 - (-2)", x: 0, res: 4, err: ERR_SUCCESS},
		{eq: "2 * (-2)", x: 0, res: -4, err: ERR_SUCCESS},
		{eq: "2 / (-2)", x: 0, res: -1, err: ERR_SUCCESS},
		{eq: "2 mod (-2)", x: 0, res: 0, err: ERR_SUCCESS},
		{eq: "2 ^ (-2)", x: 0, res: 0.25, err: ERR_SUCCESS},
		{eq: "0", x: 0, res: 0, err: ERR_SUCCESS},
		{eq: "-1", x: 0, res: -1, err: ERR_SUCCESS},
		{eq: "9999 ^ 99999999", x: 0, res: math.Inf(1), err: ERR_SUCCESS},
		{eq: "9999 ^ (-99999999)", x: 0, res: 0, err: ERR_SUCCESS},
		{eq: "0e+0", x: 0, res: 0, err: ERR_SUCCESS},
		{eq: "0e+2", x: 0, res: 0, err: ERR_SUCCESS},
		{eq: "0e-0", x: 0, res: 0, err: ERR_SUCCESS},
		{eq: "0e-2", x: 0, res: 0, err: ERR_SUCCESS},
		{eq: "1e+2", x: 0, res: 100, err: ERR_SUCCESS},
		{eq: "1e-2", x: 0, res: 0.01, err: ERR_SUCCESS},
		{eq: "log(10)", x: 0, res: 1, err: ERR_SUCCESS},
		{eq: "acos(1)", x: 0, res: 0, err: ERR_SUCCESS},
		{eq: "2 ^ 0", x: 0, res: 1, err: ERR_SUCCESS},
		{eq: "sqrt(4)", x: 0, res: 2, err: ERR_SUCCESS},
		{eq: "(8 + 3 * (4 - 2) + (2 - 3) * (2 + 4) - 16)", x: 0, res: -8, err: ERR_SUCCESS},
		{eq: "2.0 * (3 + 8 * (2 - 7)) / ((4 - 7) + 5)", x: 0, res: -37, err: ERR_SUCCESS},
		{eq: "3e+4 * (3 + 8 * (2 - 7)) / ((4 - 7) + (-5))", x: 0, res: 138750, err: ERR_SUCCESS},
		{eq: "3e-8 * (3 + 8 * (2 - 7)) / ((4 - 5) + (-5))", x: 0, res: 1.85e-07, err: ERR_SUCCESS},
		{eq: "2 ^ 2 ^ 3", x: 0, res: 64, err: ERR_SUCCESS},
		{eq: "-(-3) - (-3)", x: 0, res: 6, err: ERR_SUCCESS},
		{eq: "-(-3) - (-(-3))", x: 0, res: 0, err: ERR_SUCCESS},
		{eq: "2.0 * (-(3 + 8))", x: 0, res: -22, err: ERR_SUCCESS},
		{eq: "2 + x", x: 3, res: 5, err: ERR_SUCCESS},
		{eq: "2 ^ x", x: 3, res: 8, err: ERR_SUCCESS},
		{eq: "2 + 2", x: 0, res: 4, err: ERR_SUCCESS},
		{eq: "x - x", x: -3, res: 0, err: ERR_SUCCESS},
		{eq: "-x - x", x: -3, res: 6, err: ERR_SUCCESS},
		{eq: "log(0)", x: 0, res: math.Inf(-1), err: ERR_SUCCESS},
		{eq: "ln(0)", x: 0, res: math.Inf(-1), err: ERR_SUCCESS},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("Test #%d (%s)", i, tt.eq)
		t.Run(testname, func(t *testing.T) {
			err := c.Calculate(tt.eq, tt.x)
			res := c.Get_result()
			if err != tt.err || res != tt.res {
				t.Errorf("got %g (%s), want %g (%s)", res, errMsg[err], tt.res, errMsg[tt.err])
			}
		})
	}
}

func TestErrorMessages(t *testing.T) {
	c := NewCalcModel()
	defer DeleteCalcModel(c)

	var tests = []struct {
		eq     string
		x, res float64
		err    Err_t
	}{
		{eq: "1 / 0", x: 0, res: 4, err: ERR_DIVBYZERO},
		{eq: "1 + (", x: 0, res: 4, err: ERR_MISSBRACKET},
		{eq: "2 + sin", x: 0, res: 4, err: ERR_MALFORMED},
		{eq: "(", x: 0, res: 4, err: ERR_MISSBRACKET},
		{eq: "sin", x: 0, res: 4, err: ERR_MALFORMED},
		{eq: "(2 + (5 - 6)", x: 0, res: 4, err: ERR_MISSBRACKET},
		{eq: "sqrt(-2)", x: 0, res: 4, err: ERR_NEGATIVEROOT},
		{eq: "2 mod 0", x: 0, res: 4, err: ERR_DIVBYZERO},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("Test #%d (%s)", i, tt.eq)
		t.Run(testname, func(t *testing.T) {
			err := c.Calculate(tt.eq, tt.x)
			if err != tt.err {
				t.Errorf("got (%s), want (%s)", errMsg[err], errMsg[tt.err])
			}
		})
	}
}

func TestIsNaN(t *testing.T) {
	c := NewCalcModel()
	defer DeleteCalcModel(c)

	var tests = []struct {
		eq     string
		x, res float64
		err    Err_t
	}{
		{eq: "log(-1)", x: 0, res: math.NaN(), err: ERR_SUCCESS},
		{eq: "ln(-1)", x: 0, res: math.NaN(), err: ERR_SUCCESS},
		{eq: "asin(2)", x: 0, res: math.NaN(), err: ERR_SUCCESS},
		{eq: "acos(2)", x: 0, res: math.NaN(), err: ERR_SUCCESS},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("Test #%d (%s)", i, tt.eq)
		t.Run(testname, func(t *testing.T) {
			err := c.Calculate(tt.eq, tt.x)
			res := c.Get_result()
			if err != tt.err || !math.IsNaN(res) {
				t.Errorf("got (%s), want (%s), isNaN = %t", errMsg[err], errMsg[tt.err], math.IsNaN(res))
			}
		})
	}
}
