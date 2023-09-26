package main

import (
	"smartcalc/view"

	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.NewWithID("smartcalc")

	ui := view.NewUI()
	ui.LoadUI()
	app.Run()
}
