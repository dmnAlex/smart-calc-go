package main

import (
	"smartcalc/view"

	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.New()

	ui := view.NewUI()
	ui.LoadUI(app)
	app.Run()
}
