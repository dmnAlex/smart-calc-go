package view

import (
	"smartcalc/view/components"
	"smartcalc/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	Padding = 20
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

func (ui *UI) LoadUI() {
	app := fyne.CurrentApp()
	ui.setTheme(app)
	ui.window = app.NewWindow("SmartCalc")
	ui.window.SetMaster()
	ui.window.SetOnClosed(ui.vm.Close)
	InitLogger()

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			dialog := dialog.NewCustom("Settings", "close", components.NewSettings(), ui.window)
			dialog.Show()
		}),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			dialog := dialog.NewCustom("", "close", components.NewHelp(), ui.window)
			dialog.Resize(ui.window.Canvas().Size().Subtract(fyne.Size{Width: theme.Padding(), Height: theme.Padding()}))
			dialog.Show()
		}),
	)

	calculator := components.NewCalculator(ui.vm, ui.window)
	content := container.NewBorder(toolbar, nil, nil, nil, calculator)

	ui.window.SetContent(content)
	ui.window.Resize(fyne.NewSize(600, 800))
	ui.window.Show()
}

func (ui *UI) setTheme(app fyne.App) {
	switch app.Preferences().String("theme") {
	case "light":
		app.Settings().SetTheme(theme.LightTheme())
	case "dark":
		app.Settings().SetTheme(theme.DarkTheme())
	default:
		app.Settings().SetTheme(theme.DarkTheme())
	}
}
