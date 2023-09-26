package components

import (
	"smartcalc/viewmodel"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Calculator struct {
	widget.BaseWidget
	vm           *viewmodel.ViewModel
	inputDisplay *widget.Label
	errorDisplay *widget.Label
	spinbox      *Spinbox
	toolbar      *widget.Toolbar
	historyList  *widget.List
	buttons      map[string]*widget.Button
	window       fyne.Window
}

func NewCalculator(vm *viewmodel.ViewModel, window fyne.Window) *Calculator {
	calculator := &Calculator{
		vm:      vm,
		buttons: make(map[string]*widget.Button),
		window:  window,
	}
	calculator.ExtendBaseWidget(calculator)
	window.Canvas().SetOnTypedKey(calculator.handleKeyPress)

	return calculator
}

func (c *Calculator) addActionBtn(name string, label string, action func(), importance widget.ButtonImportance, icon fyne.Resource) *widget.Button {
	button := widget.NewButtonWithIcon(label, icon, action)
	button.Importance = importance
	c.buttons[name] = button
	return button
}

func (c *Calculator) addInputBtn(text string) *widget.Button {
	button := widget.NewButton(text, func() {
		c.vm.AddToken(text)
	})
	c.buttons[text] = button
	return button
}

func (c *Calculator) handleEqual() {
	c.vm.Equal()
	c.historyList.ScrollToBottom()
}

func (c *Calculator) CreateRenderer() fyne.WidgetRenderer {
	c.inputDisplay = widget.NewLabelWithData(c.vm.Input)
	c.inputDisplay.Alignment = fyne.TextAlignTrailing
	c.inputDisplay.TextStyle.Bold = true

	c.errorDisplay = widget.NewLabelWithData(c.vm.Error)
	c.errorDisplay.TextStyle.Italic = true

	c.spinbox = NewSpinbox(c.vm.Variable)

	clearHistory := widget.NewToolbarAction(theme.DeleteIcon(), func() {
		dialog.ShowConfirm("Delete all history", "Are you sure you want to clear history?", func(choice bool) {
			if choice {
				c.vm.DeleteHistory()
			}
		}, c.window)
	})
	c.toolbar = widget.NewToolbar(widget.NewToolbarSpacer(), clearHistory)

	c.historyList = widget.NewListWithData(
		c.vm.History,
		func() fyne.CanvasObject {
			return NewHistoryItem()
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			component := co.(*HistoryItem)
			v, _ := di.(binding.Untyped).Get()
			data := v.(viewmodel.HistoryItem)
			component.Input.SetText(strings.Join(data.Tokens, ""))
			component.Output.SetText(data.Output)
			component.OnTappedContent = func() {
				c.vm.SetTokens(data.Tokens)
			}
			component.OnTappedAction = func() {
				d := dialog.NewCustom("", "close", NewGraphPlotter(c.vm, data.Tokens), c.window)
				d.Resize(c.window.Canvas().Size().Subtract(fyne.Size{Width: theme.Padding(), Height: theme.Padding()}))
				d.Show()
			}
		},
	)

	message := container.NewHBox(layout.NewSpacer(), c.errorDisplay)
	equation := container.NewBorder(nil, message, nil, nil, container.NewHScroll(c.inputDisplay))

	keyboard := container.NewGridWithColumns(1,
		container.NewGridWithColumns(6,
			c.addInputBtn("sin"),
			c.addInputBtn("asin"),
			c.addInputBtn("sqrt"),
			c.addInputBtn("("),
			c.addInputBtn(")"),
			c.addActionBtn("delete", "", c.vm.DeleteToken, widget.DangerImportance, theme.MailReplyIcon()),
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
			c.addActionBtn("clear", "clear", c.vm.Clear, widget.DangerImportance, nil),
			c.addInputBtn("0"),
			c.addActionBtn("equal", "=", c.handleEqual, widget.HighImportance, nil),
		),
	)

	controls := container.NewVBox(widget.NewSeparator(), c.toolbar, widget.NewSeparator(), equation, c.spinbox, keyboard)
	content := container.NewBorder(nil, controls, nil, nil, c.historyList)

	return widget.NewSimpleRenderer(content)
}

func (c *Calculator) handleKeyPress(ev *fyne.KeyEvent) {
	key := ev.Name
	if key == fyne.KeyReturn || key == fyne.KeyEnter || key == fyne.KeyEqual {
		c.buttons["equal"].OnTapped()
	}

	if key == fyne.KeyBackspace {
		c.buttons["delete"].OnTapped()
	}

	if button, ok := c.buttons[string(key)]; ok {
		button.OnTapped()
	}
}
