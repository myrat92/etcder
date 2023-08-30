package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/myrat92/etcder/internal/engine/service/gui"
	"github.com/myrat92/etcder/pkg/theme"
)

func main() {
	a := app.New()
	gui.Start(a)
	a.Settings().SetTheme(&theme.MyTheme{})
	a.Run()
}
