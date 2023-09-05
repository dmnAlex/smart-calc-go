package components

import (
	"image"
	"smartcalc/model"
	"smartcalc/viewmodel"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type GraphPlotter struct {
	widget.BaseWidget
	vm       *viewmodel.ViewModel
	equation []string
}

func NewGraphPlotter(vm *viewmodel.ViewModel) *GraphPlotter {
	plotter := &GraphPlotter{
		vm:       vm,
		equation: make([]string, 0),
	}

	plotter.ExtendBaseWidget(plotter)

	return plotter
}

func (gp *GraphPlotter) CreateRenderer() fyne.WidgetRenderer {
	const (
		minX, maxX = -10.0, 10.0
		step       = 0.1
		dpi        = 96
	)

	size := int((maxX - minX) / step)
	var pointsX, pointsY []float64

	for i := 0; i < size; i++ {
		x := float64(minX + float64(i)*step)
		y, err := model.Evaluate("tan(x)", x)

		if err == "" {
			pointsX = append(pointsX, x)
			pointsY = append(pointsY, y)
		}
	}

	p := plot.New()
	p.Title.Text = "Graph"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	scatterPoints := make(plotter.XYs, len(pointsX))
	for i := 0; i < len(pointsX); i++ {
		scatterPoints[i].X = pointsX[i]
		scatterPoints[i].Y = pointsY[i]
	}

	err := plotutil.AddLinePoints(p, "Graph", scatterPoints)
	if err != nil {
		panic(err)
	}

	img := image.NewRGBA(image.Rect(0, 0, 10*dpi, 10*dpi))
	cnv := vgimg.NewWith(vgimg.UseImage(img))
	p.Draw(draw.New(cnv))

	plot := canvas.NewImageFromImage(cnv.Image())

	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, widget.NewLabel("Some buttons"), nil, plot))
}
