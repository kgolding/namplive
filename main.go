package main

import (
	"context"
	_ "embed"
	"time"

	"github.com/kgolding/nmaplive/nmap"
	"github.com/wailsapp/wails"
)

func basic() string {
	return "Hello World!"
}

//go:embed frontend/dist/app.js
var js string

//go:embed frontend/dist/app.css
var css string

func main() {

	app := wails.CreateApp(&wails.AppConfig{
		Width:  800,
		Height: 600,
		Title:  "NmapLive",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
		Resizable: true
	})
	app.Bind(basic)
	app.Bind(scan)
	app.Run()
}

type Result struct {
	// FirstSeen time.Time
	LastSeen time.Time
	// LastNotSeen time.Time
	Changed time.Time
	Active  bool
	IP      string
	Name    string
}

var cache map[string]*Result

func init() {
	cache = make(map[string]*Result)
}

// scan return new/changed devices
func scan() ([]*Result, error) {
	ctx := context.Background()
	now := time.Now()
	newResults, err := nmap.Scan(ctx, "192.168.5.0/24")
	if err != nil {
		return nil, err
	}

	for _, item := range newResults {
		if r, ok := cache[item.IP]; ok {
			if !r.Active {
				r.Changed = now
			}
			r.LastSeen = now
			r.Active = true
		} else {
			cache[item.IP] = &Result{
				Active:   true,
				LastSeen: now,
				Changed:  now,
				IP:       item.IP,
				Name:     item.Name,
			}
		}
	}
	for ip, item := range cache {
		exists := false
		for _, r := range newResults {
			if r.IP == ip {
				exists = true
			}
		}
		if !exists {
			item.Active = false
			item.Changed = now
		}
	}

	ret := make([]*Result, 0, len(cache))
	for _, item := range cache {
		ret = append(ret, item)
	}
	return ret, nil
}
