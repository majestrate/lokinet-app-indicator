package main

import (
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/gotk3/gotk3/gtk"
	"github.com/tumb1er/go-appindicator"
	"log"
	"os/exec"
	"strings"
	"time"
)

/// LokinetState represents all states lokinet can be in
type LokinetState int

const (
	Off LokinetState = iota
	Errored
	Starting
	Stopping
	NoExit
	On
)

func (st LokinetState) String() string {
	return [...]string{"Off", "Errored", "Starting", "Stopping", "No Exit", "On"}[st]
}

/// LokinetInspector gets the state of lokinet right now
type LokinetInspector interface {
	Close()
	/// determine lokinet's current state
	State() LokinetState
}

type sdLokinet struct {
	conn *dbus.Conn
}

/// NewSDLokinet creates a new systemd based lokinet inspector
func newSDLokinet() (LokinetInspector, error) {
	conn, err := dbus.New()
	if err != nil {
		return nil, err
	}
	return &sdLokinet{conn}, nil
}

func (sd sdLokinet) unitName() string {
	return "lokinet.service"
}

func (sd *sdLokinet) Close() {
	sd.conn.Close()
}

func (sd *sdLokinet) State() LokinetState {
	prop, err := sd.conn.GetUnitProperty(sd.unitName(), "ActiveState")
	if err != nil {
		return Errored
	}
	activeState := strings.Trim(strings.ToLower(prop.Value.String()), "\"")

	if activeState == "inactive" {
		return Off
	}

	if activeState == "failed" {
		return Errored
	}

	if activeState == "deactivating" {
		return Stopping
	}

	if activeState != "active" {
		return Starting
	}

	cmd := exec.Command("lokinet-vpn", "--status")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return Errored
	}

	if string(out) == "no exits\n" {
		return NoExit
	}
	return On
}

func main() {
	gtk.Init(nil)

	var iconOn = "emblem-default"
	var iconOff = "emblem-important"
	var iconErr = "emblem-error"
	var iconUnk = "emblem-question"

	menu, err := gtk.MenuNew()
	if err != nil {
		log.Fatal(err)
	}

	item, err := gtk.MenuItemNewWithLabel("exit")
	if err != nil {
		log.Fatal(err)
	}

	lokinet, err := newSDLokinet()

	if err != nil {
		log.Fatal(err)
	}

	defer lokinet.Close()

	indicator := appindicator.New("lokinet-exit-indicator", iconUnk, appindicator.CategoryApplicationStatus)
	indicator.SetStatus(appindicator.StatusActive)
	indicator.SetMenu(menu)

	item.Connect("activate", func() {
		indicator.SetLabel("activated", "")
	})

	t := time.NewTicker(time.Second)

	item.Connect("activate", func() {
		t.Stop()
		gtk.MainQuit()
	})

	menu.Add(item)
	menu.ShowAll()

	getIcon := func() string {
		switch lokinet.State() {
		case Off:
			return iconOff
		case On:
			return iconOn
		case Errored:
			return iconErr
		default:
			return iconUnk
		}
	}

	go func() {
		for {
			<-t.C
			indicator.SetIcon(getIcon())
		}
	}()

	gtk.Main()
}
