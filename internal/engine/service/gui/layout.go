package gui

import (
	"strconv"

	"fyne.io/fyne/v2/widget"

	"github.com/myrat92/etcder/internal/engine/domain/etcd"
)

// BelowValueBoard More information about values
type BelowValueBoard struct {
	Version   *widget.Label
	CreateRev *widget.Label
	ModRev    *widget.Label
	Lease     *widget.Label
}

func NewBelowValueBoard() *BelowValueBoard {
	return &BelowValueBoard{
		Version:   widget.NewLabel(""),
		CreateRev: widget.NewLabel(""),
		ModRev:    widget.NewLabel(""),
		Lease:     widget.NewLabel(""),
	}
}

func (b *BelowValueBoard) Refresh(resp *etcd.GetResp) {
	b.Version.SetText("Version: " + strconv.FormatInt(resp.Version, 10))
	b.CreateRev.SetText("CreateRev: " + strconv.FormatInt(resp.CreateRev, 10))
	b.ModRev.SetText("ModRev: " + strconv.FormatInt(resp.ModRev, 10))
	b.Lease.SetText("Lease: " + strconv.FormatInt(resp.Lease, 10))
}
