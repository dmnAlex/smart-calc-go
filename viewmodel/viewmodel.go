package viewmodel

import (
	"encoding/json"
	"image/color"
	"log"
	"math"
	"os"
	"runtime"
	"smartcalc/model"
	"strconv"
	"strings"

	"fyne.io/fyne/v2/data/binding"
	"golang.org/x/image/font"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

const (
	MinXY    = -15.0
	MaxXY    = 15.0
	minDelta = 2
	maxDelta = 2000000
	step     = 0.001
	dpi      = 100
	filename = "historydata.json"
)

type HistoryItem struct {
	Tokens []string
	Output string
}

type PlotData struct {
	minX, maxX, minY, maxY float64
	points                 plotter.XYs
	Equation               string
}

type ViewModel struct {
	model    model.CalcModel
	tokens   binding.StringList
	Input    binding.String
	Output   binding.String
	Error    binding.String
	Variable binding.Float
	History  binding.UntypedList
	PlotData PlotData
	Canvas   *vgimg.Canvas
}

func NewViewModel() *ViewModel {
	vm := &ViewModel{
		model:    model.NewCalcModel(),
		tokens:   binding.NewStringList(),
		Input:    binding.NewString(),
		Output:   binding.NewString(),
		Error:    binding.NewString(),
		Variable: binding.NewFloat(),
		History:  binding.NewUntypedList(),
		PlotData: PlotData{
			minX:   MinXY,
			maxX:   MaxXY,
			minY:   MinXY,
			maxY:   MaxXY,
			points: make(plotter.XYs, 0)},
		Canvas: vgimg.NewWith(vgimg.UseDPI(dpi), vgimg.UseWH(1280, 720)),
	}

	vm.LoadHistory()
	runtime.SetFinalizer(vm, deleteViewModel)

	return vm
}

var errorMessage = map[model.Err_t]string{
	model.ERR_MALFORMED:    "Malformed expression",
	model.ERR_DIVBYZERO:    "Division by zero",
	model.ERR_MISSBRACKET:  "Mismatched brackets",
	model.ERR_NEGATIVEROOT: "Negative root",
}

func (vm *ViewModel) updateInput() {
	equation, _ := vm.tokens.Get()
	vm.Input.Set(strings.Join(equation, ""))
}

func deleteViewModel(vm *ViewModel) {
	model.DeleteCalcModel(vm.model)
	vm.model = nil
}

func (vm *ViewModel) Close() {
	vm.SaveHistory()
	runtime.SetFinalizer(vm, nil)
	deleteViewModel(vm)
}

func (vm *ViewModel) AddToken(token string) {
	vm.tokens.Append(token)
	vm.updateInput()
	vm.Error.Set("")
}

func (vm *ViewModel) DeleteToken() {
	size := vm.tokens.Length()

	if size > 0 {
		current, _ := vm.tokens.Get()
		vm.tokens.Set(current[:(size - 1)])
		vm.updateInput()
		vm.Error.Set("")
	}
}

func (vm *ViewModel) Clear() {
	vm.tokens.Set(make([]string, 0))
	vm.updateInput()
	vm.Error.Set("")
}

func (vm *ViewModel) DeleteHistory() {
	vm.History.Set(make([]interface{}, 0))
}

func (vm *ViewModel) SaveHistory() error {
	data, err := vm.History.Get()
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func (vm *ViewModel) LoadHistory() error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var data []HistoryItem

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	vm.DeleteHistory()
	for _, item := range data {
		vm.History.Append(item)
	}

	return nil
}

func (vm *ViewModel) SetTokens(tokens []string) {
	vm.tokens.Set(tokens)
	vm.updateInput()
}

func (vm *ViewModel) Equal() {
	xString, _ := vm.Variable.Get()
	equation, _ := vm.Input.Get()
	tokens, _ := vm.tokens.Get()

	err := vm.model.Calculate(equation, xString)
	yValue := vm.model.Get_result()

	vm.Error.Set(errorMessage[err])
	if err == model.ERR_SUCCESS {
		yString := strconv.FormatFloat(yValue, 'f', 6, 64)
		vm.Output.Set(yString)
		vm.History.Append(HistoryItem{Tokens: tokens, Output: yString})
		vm.SetTokens(strings.Split(yString, ""))
		log.Printf("%s = %s", equation, yString)
	} else {
		vm.Output.Set(errorMessage[err])
	}
}

func (vm *ViewModel) PlotFillData() {
	vm.PlotData.points = make(plotter.XYs, 0)

	for x := vm.PlotData.minX; x <= vm.PlotData.maxX; x += step {
		err := vm.model.Calculate(vm.PlotData.Equation, x)
		y := vm.model.Get_result()

		if err == model.ERR_SUCCESS && !math.IsNaN(y) && y < vm.PlotData.maxY && y > vm.PlotData.minY {
			vm.PlotData.points = append(vm.PlotData.points, plotter.XY{X: x, Y: y})
		}
	}
}

func (vm *ViewModel) PlotZoomIn() {
	if vm.PlotData.maxX-vm.PlotData.minX > minDelta && vm.PlotData.maxY-vm.PlotData.minY > minDelta {
		vm.PlotData.minX++
		vm.PlotData.maxX--
		vm.PlotData.minY++
		vm.PlotData.maxY--
		vm.PlotRedraw()
	}

}

func (vm *ViewModel) PlotZoomOut() {
	if vm.PlotData.maxX-vm.PlotData.minX < maxDelta && vm.PlotData.maxY-vm.PlotData.minY < maxDelta {
		vm.PlotData.minX--
		vm.PlotData.maxX++
		vm.PlotData.minY--
		vm.PlotData.maxY++
		vm.PlotRedraw()
	}
}

func (vm *ViewModel) PlotZoomFit() {
	vm.PlotData.minX = MinXY
	vm.PlotData.maxX = MaxXY
	vm.PlotData.minY = MinXY
	vm.PlotData.maxY = MaxXY
	vm.PlotRedraw()
}

func (vm *ViewModel) MoveRight() {
	if vm.PlotData.maxX < maxDelta/2 {
		vm.PlotData.maxX++
		vm.PlotData.minX++
		vm.PlotRedraw()
	}
}

func (vm *ViewModel) MoveLeft() {
	if vm.PlotData.maxX > -maxDelta/2 {
		vm.PlotData.maxX--
		vm.PlotData.minX--
		vm.PlotRedraw()
	}
}

func (vm *ViewModel) MoveUp() {
	if vm.PlotData.maxY < maxDelta/2 {
		vm.PlotData.maxY++
		vm.PlotData.minY++
		vm.PlotRedraw()
	}
}

func (vm *ViewModel) MoveDown() {
	if vm.PlotData.maxY > -maxDelta/2 {
		vm.PlotData.maxY--
		vm.PlotData.minY--
		vm.PlotRedraw()
	}
}

func (vm *ViewModel) PlotRedraw() {
	vm.PlotFillData()
	plot := plot.New()

	plot.X.Min = vm.PlotData.minX
	plot.X.Max = vm.PlotData.maxX
	plot.Y.Min = vm.PlotData.minY
	plot.Y.Max = vm.PlotData.maxY

	plot.X.Label.Text = "X"
	plot.X.Label.TextStyle.Font.Weight = font.WeightBold
	plot.X.Label.TextStyle.Font.Size = 32
	plot.X.Width = 3
	plot.X.Tick.Label.Font.Size = 24
	plot.X.Tick.Width = 1.5

	plot.Y.Label.Text = "Y"
	plot.Y.Label.TextStyle.Font.Weight = font.WeightBold
	plot.Y.Label.TextStyle.Font.Size = 32
	plot.Y.Width = 3
	plot.Y.Tick.Label.Font.Size = 24
	plot.Y.Tick.Width = 1.5

	grid := plotter.NewGrid()
	grid.Horizontal.Width = 2
	grid.Vertical.Width = 2
	plot.Add(grid)

	line, _ := plotter.NewScatter(vm.PlotData.points)
	line.Shape = draw.CircleGlyph{}
	line.Color = color.RGBA{B: 255, A: 255}
	line.Radius = 2
	plot.Add(line)

	plot.Draw(draw.New(vm.Canvas))
}
