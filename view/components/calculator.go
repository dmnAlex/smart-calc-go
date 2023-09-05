package components

import (
	"smartcalc/viewmodel"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Calculator struct {
	widget.BaseWidget
	vm            *viewmodel.ViewModel
	inputDisplay  *widget.Label
	outputDisplay *widget.Label
	entryField    *widget.Entry
	historyList   *widget.List
	buttons       map[string]*widget.Button
}

func NewCalculator(vm *viewmodel.ViewModel) *Calculator {
	calculator := &Calculator{
		vm:      vm,
		buttons: make(map[string]*widget.Button),
	}

	calculator.ExtendBaseWidget(calculator)

	return calculator
}

func (c *Calculator) addBtn(text string, action func(), importance widget.ButtonImportance) *widget.Button {
	button := widget.NewButton(text, action)
	button.Importance = importance
	c.buttons[text] = button
	return button
}

func (c *Calculator) addInputBtn(text string) *widget.Button {
	return c.addBtn(text, func() {
		c.vm.AddToken(text)
	}, widget.MediumImportance)
}

func (c *Calculator) handleEqual() {
	c.vm.Equal()
	c.historyList.ScrollToBottom()
}

func (c *Calculator) CreateRenderer() fyne.WidgetRenderer {
	c.inputDisplay = widget.NewLabelWithData(c.vm.Input)
	c.inputDisplay.Alignment = fyne.TextAlignTrailing
	c.outputDisplay = widget.NewLabelWithData(c.vm.Output)
	c.outputDisplay.Alignment = fyne.TextAlignTrailing

	c.entryField = widget.NewEntryWithData(c.vm.Variable)

	c.historyList = widget.NewListWithData(
		c.vm.History,
		func() fyne.CanvasObject {
			label := &widget.Label{Alignment: fyne.TextAlignTrailing}
			return label
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			label := co.(*widget.Label)
			item, _ := di.(binding.Untyped).Get()
			historyItem := item.(viewmodel.HistoryItem)

			label.SetText(strings.Join(historyItem.Tokens, ""))
		},
	)

	content := container.NewGridWithColumns(1,
		c.historyList,
		c.inputDisplay,
		c.outputDisplay,
		c.entryField,
		container.NewGridWithColumns(6,
			c.addInputBtn("sin"),
			c.addInputBtn("asin"),
			c.addInputBtn("sqrt"),
			c.addInputBtn("("),
			c.addInputBtn(")"),
			c.addBtn("↩", c.vm.DeleteToken, widget.DangerImportance),
		),
		container.NewGridWithColumns(6,
			c.addInputBtn("cos"),
			c.addInputBtn("acos"),
			c.addInputBtn("^"),
			c.addInputBtn("."),
			c.addInputBtn("e"),
			c.addInputBtn("+"),
		),
		container.NewGridWithColumns(6,
			c.addInputBtn("tan"),
			c.addInputBtn("atan"),
			c.addInputBtn("7"),
			c.addInputBtn("8"),
			c.addInputBtn("9"),
			c.addInputBtn("-"),
		),
		container.NewGridWithColumns(6,
			c.addInputBtn("log"),
			c.addInputBtn("ln"),
			c.addInputBtn("4"),
			c.addInputBtn("5"),
			c.addInputBtn("6"),
			c.addInputBtn("*"),
		),
		container.NewGridWithColumns(6,
			c.addInputBtn("mod"),
			c.addInputBtn("x"),
			c.addInputBtn("1"),
			c.addInputBtn("2"),
			c.addInputBtn("3"),
			c.addInputBtn("/"),
		),
		container.NewGridWithColumns(3,
			c.addBtn("C", c.vm.Clear, widget.DangerImportance),
			c.addInputBtn("0"),
			c.addBtn("=", c.handleEqual, widget.HighImportance),
		),
	)

	return widget.NewSimpleRenderer(content)
}
