package client

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var cmd = map[string]string{
	"All groups":  "--groups",
	"New group":   "--create",
	"Join group":  "--join",
	"Leave group": "--leave",
	"Quit":        "--quit",
}

var (
	viewArr = []string{"v2", "v4", "v1"}
	active  = 0
)

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var command string
	var err error
	_, cy := v.Cursor()
	if command, err = v.Line(cy); err != nil {
		command = ""
	}
	msg, ok := cmd[command]
	if !ok {
		return nil
	}
	send(msg)
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	out, err := g.View("v2")
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "Going from view "+v.Name()+" to "+name)

	g.SetCurrentView(name)

	if nextIndex == 1 || nextIndex == 0 {
		g.Cursor = true
	} else {
		// v.Clear()
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

func sendMsg(g *gocui.Gui, v *gocui.View) error {
	var msg string
	var err error
	_, cy := v.Cursor()
	if msg, err = v.Line(cy); err != nil {
		msg = ""
	}
	v.Clear()
	v.MoveCursor(-len(msg), 0, true)
	send(msg)
	return nil
}

func send(command string) {
	if command == "" {
		return
	}
	// oi.input <- command
	// oi.output <- command
	oi.input <- command
}

func viewUpdate(g *gocui.Gui) {
	defer wg.Done()
	for {
		select {
		case <-stop:
			return
		case msg := <-oi.output:
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("v2")
				if err != nil {
					return err
				}
				fmt.Fprintln(v, msg)
				return nil
			})
		}
	}
}
