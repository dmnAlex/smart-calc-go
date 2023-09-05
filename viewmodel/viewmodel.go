package viewmodel

import (
	"smartcalc/model"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/data/binding"
)

type HistoryItem struct {
	Tokens []string
	Output string
}

type ViewModel struct {
	model    func(equation string, x float64) (float64, string)
	tokens   binding.StringList
	Input    binding.String
	Output   binding.String
	Variable binding.String
	History  binding.UntypedList
}

func NewViewModel() *ViewModel {
	vm := &ViewModel{
		model:    model.Evaluate,
		tokens:   binding.NewStringList(),
		Input:    binding.NewString(),
		Output:   binding.NewString(),
		Variable: binding.NewString(),
		History:  binding.NewUntypedList(),
	}

	vm.tokens.AddListener(binding.NewDataListener(func() {
		equation, _ := vm.tokens.Get()
		vm.Input.Set(strings.Join(equation, ""))
	}))

	return vm
}

func (vm *ViewModel) AddToken(token string) {
	vm.tokens.Append(token)
}

func (vm *ViewModel) DeleteToken() {
	size := vm.tokens.Length()

	if size > 0 {
		current, _ := vm.tokens.Get()
		vm.tokens.Set(current[:(size - 1)])
	}
}

func (vm *ViewModel) Clear() {
	vm.tokens.Set(make([]string, 0))
}

func (vm *ViewModel) DeleteHistory() {
	vm.History.Set(make([]interface{}, 0))
}

func (vm *ViewModel) Equal() {
	xString, _ := vm.Variable.Get()
	xValue, _ := strconv.ParseFloat(xString, 64)
	equation, _ := vm.Input.Get()
	tokens, _ := vm.tokens.Get()

	yValue, err := model.Evaluate(equation, xValue)
	if err == "" {
		yString := strconv.FormatFloat(yValue, 'f', 6, 64)
		vm.Output.Set(yString)
		vm.History.Append(HistoryItem{Tokens: tokens, Output: yString})
	} else {
		vm.Output.Set(err)
	}
}
