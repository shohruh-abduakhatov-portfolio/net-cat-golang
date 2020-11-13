package client

import (
	"fmt"
	"sync"

	"github.com/jroimartin/gocui"
)

var oi *IO
var wg sync.WaitGroup
var stop = make(chan string)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("v1", 0, 0, int(0.25*float32(maxX))-1, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v1 (editable)"
		v.Wrap = true
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for k, _ := range cmd {
			fmt.Fprintln(v, k)
		}
		g.SetCurrentView("v1")
	}

	if v, err := g.SetView("v2", int(0.25*float32(maxX)), 0, maxX-1, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v2"
		v.Wrap = true
		v.Autoscroll = true
		// v.Editable = true
	}
	if v, err := g.SetView("v3", 0, maxY-4, int(0.25*float32(maxX))-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v3"
		v.Wrap = true
		// v.Autoscroll = true
		fmt.Fprintln(v, "Tab: Change View")
		fmt.Fprintln(v, "^C: Exit")
	}
	if v, err := g.SetView("v4", int(0.25*float32(maxX)), maxY-4, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v4 (editable)"
		v.Editable = true
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(stop)
	return gocui.ErrQuit
}

func Main(o *IO) error {
	oi = o
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()
	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen
	g.SetManagerFunc(layout)
	if err := keybindings(g); err != nil {
		return err
	}
	wg.Add(1)
	go viewUpdate(g)
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	wg.Wait()
	return nil
}
