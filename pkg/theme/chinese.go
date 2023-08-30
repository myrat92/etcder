package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"github.com/myrat92/etcder/resource"
)

type MyTheme struct{}

var _ fyne.Theme = (*MyTheme)(nil)

func (m MyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}
func (m MyTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	return &fyne.StaticResource{
		StaticName:    "SourceHanSans-Normal.ttf",
		StaticContent: resource.SHStf,
	}
}

func (m MyTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
