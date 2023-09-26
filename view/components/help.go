package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const help = `
# Smart Calc
---
*by hsensor, Version 3.0*
## Description

Smart Calc is a feature-rich calculator that supports exponential notation, trigonometric functions, and graph plotting. The core of this calculator is implemented as an algorithm for expression parsing (Dijkstra's algorithm) and calculation of Polish notation. Various computational functions are implemented in C/C++ and integrated into the main codebase, while the program itself is written in Golang.

## Features

1. Arithmetic operators:

    - Parentheses (a + b),
    - Addition: a + b,
    - Subtraction: a - b,
    - Division: a / b,
    - Multiplication: a * b,
    - Exponentiation: a ^ b,
    - Modulus: a Mod b,
    - Unary minus: -a;
2. Functions:

    - Cosine: cos(x),
    - Sine: sin(x),
    - Tangent: tan(x),
    - Arccosine: acos(x),
    - Arcsine: asin(x),
    - Arctangent: atan(x),
    - Square root: sqrt(x),
    - Natural logarithm: ln(x),
    - Decimal logarithm: log(x)
`

func NewHelp() *widget.RichText {
	widget := widget.NewRichTextFromMarkdown(help)
	widget.Wrapping = fyne.TextWrapWord
	widget.Scroll = container.ScrollVerticalOnly

	return widget
}
