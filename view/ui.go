package view

import (
	"smartcalc/view/components"
	"smartcalc/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	vm     *viewmodel.ViewModel
	window fyne.Window
}

func NewUI() *UI {
	return &UI{
		vm: viewmodel.NewViewModel(),
	}
}

func (ui *UI) LoadUI(app fyne.App) {
	ui.window = app.NewWindow("SmartCalc")

	tabs := container.NewAppTabs(
		container.NewTabItem("Calculator", components.NewCalculator(ui.vm)),
		container.NewTabItem("Plot", components.NewGraphPlotter(ui.vm)),
		container.NewTabItem("Loan", widget.NewLabel("Loan calculator")),
		container.NewTabItem("Deposit", widget.NewLabel("Deposit calculator")),
	)

	ui.window.SetContent(tabs)
	ui.window.Resize(fyne.NewSize(400, 600))
	ui.window.Show()
}
