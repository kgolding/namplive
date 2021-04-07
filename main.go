package main

import (
	"context"
	_ "embed"
	"net"
	"time"

	"github.com/kgolding/nmaplive/network"
	"github.com/kgolding/nmaplive/nmap"
	"github.com/wailsapp/wails"
)

//go:embed frontend/dist/app.js
var js string

//go:embed frontend/dist/app.css
var css string

func main() {

	app := wails.CreateApp(&wails.AppConfig{
		Width:     800,
		Height:    600,
		Title:     "NmapLive",
		JS:        js,
		CSS:       css,
		Colour:    "#fff",
		Resizable: true,
	})
	app.Bind(scan)
	app.Bind(GetIPv4NonLocalInterfaces)

	app.Run()
}

func GetIPv4NonLocalInterfaces() []string {
	return network.GetIPv4NonLocalInterfaces()
}

type Result struct {
	LastSeen time.Time
	Changed  time.Time
	Active   bool
	IP       string
	Name     string
}

var cache map[string]*Result

func init() {
	cache = make(map[string]*Result)
}

// scan return new/changed devices
func scan(target string) ([]*Result, error) {
	if _, _, err := net.ParseCIDR(target); err != nil {
		return nil, err
	}
	ctx := context.Background()
	now := time.Now()
	newResults, err := nmap.Scan(ctx, target)
	if err != nil {
		return nil, err
	}

	for _, item := range newResults { // Check new items
		if r, ok := cache[item.IP]; ok { // Already in cache
			if !r.Active {
				r.Changed = now
			}
			r.LastSeen = now
			r.Active = true
		} else { // Add to cache
			cache[item.IP] = &Result{
				Active:   true,
				LastSeen: now,
				Changed:  now,
				IP:       item.IP,
				Name:     item.Name,
			}
		}
	}
	for ip, item := range cache { // Check cache
		exists := false
		for _, r := range newResults {
			if r.IP == ip {
				exists = true
			}
		}
		if !exists {
			if item.Active {
				item.Active = false
				item.Changed = now.Add(-time.Second) // sorts inactive below active
			}
		}
	}

	ret := make([]*Result, 0, len(cache))
	for _, item := range cache {
		ret = append(ret, item)
	}
	return ret, nil
}
