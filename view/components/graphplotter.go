package components

import (
	"smartcalc/viewmodel"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type GraphPlotter struct {
	widget.BaseWidget
	vm     *viewmodel.ViewModel
	tokens []string
}

func NewGraphPlotter(vm *viewmodel.ViewModel, tokens []string) *GraphPlotter {
	plotter := &GraphPlotter{
		vm:     vm,
		tokens: tokens,
	}
	plotter.ExtendBaseWidget(plotter)

	return plotter
}

func (gp *GraphPlotter) CreateRenderer() fyne.WidgetRenderer {
	gp.vm.PlotData.Equation = strings.Join(gp.tokens, "")
	gp.vm.PlotRedraw()
	plot := canvas.NewImageFromImage(gp.vm.Canvas.Image())
	plot.FillMode = canvas.ImageFillContain

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
			gp.vm.MoveLeft()
			plot.Refresh()
		}),
		widget.NewToolbarAction(theme.NavigateNextIcon(), func() {
			gp.vm.MoveRight()
			plot.Refresh()
		}),
		widget.NewToolbarAction(theme.MoveUpIcon(), func() {
			gp.vm.MoveUp()
			plot.Refresh()
		}),
		widget.NewToolbarAction(theme.MoveDownIcon(), func() {
			gp.vm.MoveDown()
			plot.Refresh()
		}),
		widget.NewToolbarAction(theme.ZoomInIcon(), func() {
			gp.vm.PlotZoomIn()
			plot.Refresh()
		}),
		widget.NewToolbarAction(theme.ZoomOutIcon(), func() {
			gp.vm.PlotZoomOut()
			plot.Refresh()
		}),
		widget.NewToolbarAction(theme.ZoomFitIcon(), func() {
			gp.vm.PlotZoomFit()
			plot.Refresh()
		}),
		widget.NewToolbarSpacer(),
	)

	content := container.NewBorder(nil, toolbar, nil, nil, plot)

	return widget.NewSimpleRenderer(content)
}
