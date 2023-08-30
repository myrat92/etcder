package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/exp/slog"

	"github.com/myrat92/etcder/internal/engine/domain/etcd"
	"github.com/myrat92/etcder/internal/engine/infrastructure/session"
)

func Start(app fyne.App) {
	tabs := container.NewDocTabs()

	connectButtonOnTapped := func() {
		t := container.NewTabItem("New Tab", NewDataPage())
		tabs.Items = append(tabs.Items, t)
		tabs.Select(t)
		tabs.Refresh()
	}

	tabs.CreateTab = func() *container.TabItem {
		slog.Warn("create tab newLoginPage")
		return container.NewTabItem("New Tab", NewLoginPage(connectButtonOnTapped))
	}

	login := NewLoginPage(connectButtonOnTapped)

	tabs.Append(container.NewTabItem("New Tab", login))

	w := app.NewWindow("Login")
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(1080, 768))
	w.Show()
}

func NewLoginPage(connectButtonOnTapped func()) fyne.CanvasObject {
	nameBind := binding.NewString()
	hostBind := binding.NewString()
	portBind := binding.NewString()
	idbind := binding.NewInt()
	idbind.Set(-1)

	sessions := session.ListSession()

	list := widget.NewList(
		func() int {
			return len(sessions)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template Object")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(sessions[id].Name + "   " + sessions[id].Host + ":" + sessions[id].Port)
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		nameBind.Set(sessions[id].Name)
		hostBind.Set(sessions[id].Host)
		portBind.Set(sessions[id].Port)
		idbind.Set(id)
	}

	list.OnUnselected = func(id widget.ListItemID) {
		idbind.Set(-1)
	}

	name := widget.NewEntryWithData(nameBind)
	host := widget.NewEntryWithData(hostBind)
	port := widget.NewEntryWithData(portBind)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: name},
			{Text: "Host", Widget: host},
			{Text: "Port", Widget: port},
		},
	}

	tipsLabel := widget.NewLabel("")

	saveButton := widget.NewButton("Save", func() {
		id, _ := idbind.Get()
		session.UpdateSession(name.Text, host.Text, port.Text, id)

		sessions = session.ListSession()
		list.Refresh()
	})

	connectButton := widget.NewButton("Connect", func() {})

	connectButton.OnTapped = func() {
		etcd.NewInstance(host.Text + ":" + port.Text)

		err := etcd.Instance().Health()
		if err != nil {
			tipsLabel.SetText(err.Error())
			return
		}
		tipsLabel.SetText("")

		connectButtonOnTapped()
	}

	listToolBar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			session.AddSession("", "", "")
			sessions = session.ListSession()
			list.Refresh()
			list.Select(list.Length() - 1)

		}),

		widget.NewToolbarAction(theme.DeleteIcon(), func() {
			id, _ := idbind.Get()
			session.DeleteSession(id)
			sessions = session.ListSession()
			list.Refresh()
			list.Select(list.Length() - 1)
		}),
	)

	buttonBox := container.NewVBox(form, container.NewCenter(container.NewHBox(saveButton, connectButton)))
	box := container.NewVBox(buttonBox, tipsLabel)

	connection := container.NewBorder(widget.NewLabel("New Connection"), nil, nil, nil, box)

	right := container.New(layout.NewMaxLayout(), connection)
	left := container.NewBorder(nil, listToolBar, nil, nil, list)

	split := container.NewHSplit(left, right)

	return split
}

func NewDataPage() fyne.CanvasObject {
	d, err := etcd.Instance().ListAll()
	fmt.Println(d)
	if err != nil {
		slog.Warn("list key in etcd", err)
	}

	value := binding.NewString()
	valueEntry := widget.NewEntryWithData(value)
	valueEntry.MultiLine = true
	valueEntry.Disabled()

	list := widget.NewList(
		func() int {
			return len(d)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template Object")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(d[id])
		},
	)

	var key string
	list.OnSelected = func(id widget.ListItemID) {
		err = value.Set(etcd.Instance().Get(d[id]))
		if err != nil {
			slog.Warn("get value", err)
		}
		key = d[id]
	}

	// toolbar
	// refresh and save
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			d, err = etcd.Instance().ListAll()
			if err != nil {
				slog.Warn("refresh list", err)
			}
			list.Refresh()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			err = etcd.Instance().Update(key, valueEntry.Text)
			if err != nil {
				slog.Warn("update value", err)
			}
		}),
	)

	split := container.NewHSplit(list, valueEntry)
	content := container.NewBorder(toolbar, nil, nil, nil, split)

	return content
}
