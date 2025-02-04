package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type roundButtonTheme struct{}

func (m *roundButtonTheme) Color(color fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(color, variant)
}

func (m *roundButtonTheme) Icon(iconName fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(iconName)
}

var _ fyne.Theme = (*roundButtonTheme)(nil)

func (m roundButtonTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m roundButtonTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameInputRadius {
		return 100
	}
	return theme.DefaultTheme().Size(name)
}
