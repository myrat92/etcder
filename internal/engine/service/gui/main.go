package gui

import (
	"encoding/json"
	"strings"

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

var (
	Unselected = -1
)

func Start(app fyne.App) {
	tabs := container.NewDocTabs()

	connectButtonOnTapped := func(title string) {
		t := container.NewTabItem(title, NewDataPage())
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

func NewLoginPage(connectButtonOnTapped func(title string)) fyne.CanvasObject {
	nameBind := binding.NewString()
	hostBind := binding.NewString()
	portBind := binding.NewString()
	idbind := binding.NewInt()
	idbind.Set(Unselected)

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
		if id == Unselected {
			session.AddSession(name.Text, host.Text, port.Text)
		} else {
			session.UpdateSession(name.Text, host.Text, port.Text, id)
		}

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

		connectButtonOnTapped(name.Text)
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
	listAll, err := etcd.Instance().ListAll()

	if err != nil {
		slog.Warn("list key in etcd", err)
	}

	value := binding.NewString()
	valueEntry := widget.NewMultiLineEntry()
	valueEntry.Wrapping = fyne.TextWrapWord
	valueEntry.Bind(value)

	list := widget.NewList(
		func() int {
			return len(listAll)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template Object")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(listAll[id])
		},
	)

	var key string
	list.OnSelected = func(id widget.ListItemID) {
		err = value.Set(etcd.Instance().Get(listAll[id]))
		if err != nil {
			slog.Warn("get value", err)
		}
		key = listAll[id]
	}

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search")

	searchEntry.OnChanged = func(query string) {
		listAll, err = etcd.Instance().ListAll()

		if err != nil {
			slog.Warn("list key in etcd", err)
		}

		var newList []string
		for _, d := range listAll {
			if strings.Contains(d, query) {
				newList = append(newList, d)
			}
		}

		listAll = newList
		list.Refresh()
	}

	// toolbar
	// refresh and save
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.ViewRefreshIcon(), func() {
			listAll, err = etcd.Instance().ListAll()
			if err != nil {
				slog.Warn("refresh list", err)
			}
			list.Refresh()

			err = value.Set(etcd.Instance().Get(key))
			if err != nil {
				slog.Warn("get value", err)
			}
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			err = etcd.Instance().Update(key, valueEntry.Text)
			if err != nil {
				slog.Warn("update value", err)
			}
		}),
	)

	dropdown := makeFormatSelect(valueEntry)

	topBox := container.New(layout.NewHBoxLayout(), dropdown, toolbar)

	listBorder := container.NewBorder(searchEntry, nil, nil, nil, list)

	split := container.NewHSplit(listBorder, container.NewBorder(topBox, nil, nil, nil, valueEntry))
	content := container.NewBorder(nil, nil, nil, nil, split)

	return content
}

func makeFormatSelect(valueEntry *widget.Entry) *widget.Select {
	options := []string{"json", "txt"}
	dropdown := widget.NewSelect(options, func(selected string) {
		switch selected {
		case "json":
			var jsonData map[string]interface{}
			err := json.Unmarshal([]byte(valueEntry.Text), &jsonData)
			if err != nil {
				slog.Warn("json unmarshal", err)
				return
			}

			formattedJson, err := json.MarshalIndent(jsonData, "", "    ")
			if err != nil {
				slog.Warn("json marshal", err)
				return
			}
			valueEntry.SetText(string(formattedJson))
		case "txt":
			var jsonData map[string]interface{}
			err := json.Unmarshal([]byte(valueEntry.Text), &jsonData)
			if err != nil {
				slog.Warn("json unmarshal", err)
				return
			}

			formattedStr, err := json.Marshal(jsonData)
			if err != nil {
				slog.Warn("json marshal", err)
				return
			}
			valueEntry.SetText(string(formattedStr))
		}
	})

	return dropdown
}
