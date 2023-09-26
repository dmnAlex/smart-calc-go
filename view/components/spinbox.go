package components

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Spinbox struct {
	widget.BaseWidget
	data binding.Float
}

func NewSpinbox(data binding.Float) *Spinbox {
	spinbox := &Spinbox{data: data}
	spinbox.ExtendBaseWidget(spinbox)

	return spinbox
}

func (s *Spinbox) CreateRenderer() fyne.WidgetRenderer {
	label := widget.NewLabelWithStyle("X :", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	entry := widget.NewEntryWithData(binding.FloatToString(s.data))
	entry.Bind(binding.FloatToString(s.data))
	entry.PlaceHolder = "Enter the value here.."

	entry.Validator = func(s string) error {
		_, err := strconv.ParseFloat(s, 64)
		return err
	}

	incBtn := widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
		value, _ := s.data.Get()
		s.data.Set(value + 1)
	})
	decBtn := widget.NewButtonWithIcon("", theme.MoveDownIcon(), func() {
		value, _ := s.data.Get()
		s.data.Set(value - 1)
	})

	buttons := container.NewHBox(incBtn, decBtn)

	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, label, buttons, entry))
}
