package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewSettings() *widget.Form {
	app := fyne.CurrentApp()

	themeWidget := widget.NewSelect([]string{"dark", "light"}, func(choice string) {
		switch choice {
		case "dark":
			app.Preferences().SetString("theme", "dark")
			app.Settings().SetTheme(theme.DarkTheme())
		case "light":
			app.Preferences().SetString("theme", "light")
			app.Settings().SetTheme(theme.LightTheme())
		}
	})
	if app.Preferences().String("theme") == "" {
		themeWidget.Selected = "dark"
	} else {
		themeWidget.Selected = app.Preferences().String("theme")
	}

	periodWidget := widget.NewSelect([]string{"hourly", "dayly", "monthly"}, func(choice string) {
		app.Preferences().SetString("rotation", choice)
	})

	if app.Preferences().String("rotation") == "" {
		periodWidget.Selected = "hourly"
	} else {
		periodWidget.Selected = app.Preferences().String("rotation")
	}

	return widget.NewForm(
		widget.NewFormItem("Log rotation (restart may be required)", periodWidget),
		widget.NewFormItem("Color theme", themeWidget),
	)
}
