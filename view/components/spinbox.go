package components

import (
	"strconv"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Spinbox struct {
	widget.Entry
}

func NewSpinbox(placeholder string) *Spinbox {
	entry := &Spinbox{}
	entry.ExtendBaseWidget(entry)
	entry.PlaceHolder = placeholder
	incBtn := widget.NewButtonWithIcon("", theme.MoveUpIcon(), func() {
		value := entry.GetValue()
		entry.SetValue(value + 1)
	})
	decBtn := widget.NewButtonWithIcon("", theme.MoveDownIcon(), func() {
		value := entry.GetValue()
		entry.SetValue(value - 1)
	})
	entry.ActionItem = container.NewGridWithRows(2, incBtn, decBtn)
	entry.Validator = func(s string) error {
		_, err := strconv.ParseFloat(s, 64)
		return err
	}

	return entry
}

func (e *Spinbox) TypedRune(r rune) {
	if (r >= '0' && r <= '9') || r == '.' || r == ',' {
		e.Entry.TypedRune(r)
	}
}

func (e *Spinbox) SetValue(value float64) {
	e.SetText(strconv.FormatFloat(value, 'f', 6, 64))
}

func (e *Spinbox) GetValue() float64 {
	value, _ := strconv.ParseFloat(e.Text, 64)
	return value
}

// func (e *Spinbox) TypedShortcut(shortcut fyne.Shortcut) {
// 	paste, ok := shortcut.(*fyne.ShortcutPaste)
// 	if !ok {
// 		e.Entry.TypedShortcut(shortcut)
// 		return
// 	}

// 	content := paste.Clipboard.Content()
// 	if _, err := strconv.ParseFloat(content, 64); err == nil {
// 		e.Entry.TypedShortcut(shortcut)
// 	}
// }

// func (e *Spinbox) Keyboard() mobile.KeyboardType {
// 	return mobile.NumberKeyboard
// }
