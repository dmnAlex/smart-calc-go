package main

import (
	"smartcalc/view"

	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.New()

	c := view.NewCalculator()
	c.LoadUI(app)
	app.Run()
}
