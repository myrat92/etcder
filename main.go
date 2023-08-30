package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/myrat92/etcder/internal/engine/service/gui"
)

func main() {
	a := app.New()
	gui.Start(a)
	a.Run()
}
