package model

/*
#cgo CFLAGS: -I./lib
#cgo LDFLAGS: -L./lib -lcalcmodel -Wl,-rpath=./model/lib
#include "calcmodelwrapper.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

func Evaluate(equation string, x float64) (float64, string) {
	str := C.CString(equation)
	defer C.free(unsafe.Pointer(str))
	num := C.double(x)

	res := C.calculate(str, num)

	return float64(res.value), C.GoString(res.error)
}
