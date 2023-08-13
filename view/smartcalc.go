package view

import (
	"smartcalc/model"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type smartcalc struct {
	equation      []string
	x             float64
	window        fyne.Window
	inputDisplay  *widget.Label
	outputDisplay *widget.Label
	entryField    *Spinbox
	buttons       map[string]*widget.Button
}

func NewCalculator() *smartcalc {
	return &smartcalc{
		buttons: make(map[string]*widget.Button),
	}
}

func (c *smartcalc) addBtn(text string, action func(), importance widget.ButtonImportance) *widget.Button {
	button := widget.NewButton(text, action)
	button.Importance = importance
	c.buttons[text] = button
	return button
}

func (c *smartcalc) addInputBtn(text string) *widget.Button {
	return c.addBtn(text, func() {
		c.add(text)
	}, widget.MediumImportance)
}

func (c *smartcalc) clear() {
	c.equation = nil
	c.inputDisplay.SetText("")
}

func (c *smartcalc) equal() {
	val, err := model.Evaluate(strings.Join(c.equation, ""), c.x)
	if err == "" {
		c.outputDisplay.SetText(strconv.FormatFloat(val, 'f', 6, 64))
	} else {
		c.outputDisplay.SetText(err)
	}
}

func (c *smartcalc) add(text string) {
	c.equation = append(c.equation, text)
	c.inputDisplay.SetText(strings.Join(c.equation, ""))
}

func (c *smartcalc) delete() {
	size := len(c.equation)
	if size > 0 {
		c.equation = c.equation[:(size - 1)]
		c.inputDisplay.SetText(strings.Join(c.equation, ""))
	}
}

func (c *smartcalc) LoadUI(app fyne.App) {
	c.inputDisplay = widget.NewLabel("Input")
	c.inputDisplay.Alignment = fyne.TextAlignTrailing
	c.outputDisplay = widget.NewLabel("Output")
	c.outputDisplay.Alignment = fyne.TextAlignTrailing

	c.window = app.NewWindow("SmartCalc")
	c.entryField = NewSpinbox("Enter X value here")
	c.entryField.OnChanged = func(_ string) {
		c.x = c.entryField.GetValue()
	}
	// c.entryField = widget.NewEntry()
	// c.entryField.PlaceHolder = "Enter X value here"
	// c.entryField.Validator = func(s string) error {
	// 	_, err := strconv.ParseFloat(s, 64)
	// 	return err
	// }
	// c.entryField.OnChanged = func(s string) {
	// 	if c.entryField.Validate() == nil {
	// 		c.x, _ = strconv.ParseFloat(c.entryField.Text, 64)
	// 	} else {
	// 		c.x = 0.0
	// 	}
	// }
	// c.entryField.ActionItem = container.NewVBox(
	// 	widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {}),
	// 	widget.NewButtonWithIcon("", theme.MoveDownIcon(), func() {}),
	// )

	content := container.NewGridWithColumns(1,
		c.inputDisplay,
		c.outputDisplay,
		c.entryField,
		container.NewGridWithColumns(6,
			c.addInputBtn("sin"),
			c.addInputBtn("asin"),
			c.addInputBtn("sqrt"),
			c.addInputBtn("("),
			c.addInputBtn(")"),
			c.addBtn("↩", c.delete, widget.DangerImportance),
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
			c.addBtn("C", c.clear, widget.DangerImportance),
			c.addInputBtn("0"),
			c.addBtn("=", c.equal, widget.HighImportance),
		),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Calculator", content),
		container.NewTabItem("Loan", widget.NewLabel("Loan calculator")),
		container.NewTabItem("Deposit", widget.NewLabel("Deposit calculator")),
	)

	c.window.SetContent(tabs)
	c.window.Resize(fyne.NewSize(400, 600))
	c.window.Show()
}
