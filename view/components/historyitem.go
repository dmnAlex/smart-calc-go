package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type HistoryItem struct {
	widget.BaseWidget
	Toolbar         *widget.Toolbar
	Input           *widget.Label
	Output          *widget.Label
	OnTappedContent func()
	OnTappedAction  func()
}

func NewHistoryItem() *HistoryItem {
	hi := &HistoryItem{}
	hi.ExtendBaseWidget(hi)

	return hi
}

func (hi *HistoryItem) CreateRenderer() fyne.WidgetRenderer {
	hi.Input = &widget.Label{Alignment: fyne.TextAlignLeading}
	hi.Output = &widget.Label{Alignment: fyne.TextAlignTrailing}
	equal := &widget.Label{Alignment: fyne.TextAlignCenter, Text: "="}

	content := container.NewGridWithColumns(3,
		container.NewHScroll(hi.Input),
		equal,
		container.NewHScroll(hi.Output),
	)

	hi.Toolbar = widget.NewToolbar(
		widget.NewToolbarAction(theme.VisibilityIcon(), func() {
			if hi.OnTappedAction != nil {
				hi.OnTappedAction()
			}
		}),
	)

	return widget.NewSimpleRenderer(container.NewBorder(nil, nil, hi.Toolbar, nil, content))
}

func (hi *HistoryItem) Tapped(event *fyne.PointEvent) {
	if hi.OnTappedContent != nil {
		hi.OnTappedContent()
	}
}
