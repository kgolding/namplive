package nmap

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type ScanResult []ScanItem

type ScanItem struct {
	IP   string
	Name string
	Up   bool
	// Seen time.Time
}

func (si *ScanItem) String() string {
	up := "OK"
	if !si.Up {
		up = "DOWN"
	}
	return fmt.Sprintf("%s\t%s\t%s", si.IP, up, si.Name)
}

func (s ScanResult) String() string {
	if len(s) == 0 {
		return ""
	}
	a := make([]string, len(s))
	for i, b := range s {
		a[i] = b.String()
	}
	return strings.Join(a, "\n")
}

func (s ScanResult) Find(ip string) *ScanItem {
	for _, a := range s {
		if a.IP == ip {
			return &a
		}
	}
	return nil
}

func (s ScanResult) Diff(s2 ScanResult) {
	// Look for new & changed IP's
	for _, a := range s2 {
		if b := s.Find(a.IP); b != nil {
			if a.Name != b.Name {
				fmt.Println("Changed name", a.String())
			}
			if a.Up != b.Up {
				fmt.Println("Changed state", a.String())
			}
		} else {
			fmt.Println("NEW", a.String())
		}
	}
	// Look for lost IPs
	for _, a := range s {
		if b := s2.Find(a.IP); b == nil {
			fmt.Println("GONE", a.String())
		}
	}
}

var nampOgRegex = regexp.MustCompile(`(?m)Host: ([0-9\.]+) \((.*)\)\s+Status: (\S+)`)

func Scan(ctx context.Context, target string) (ScanResult, error) {
	out, err := exec.CommandContext(ctx, "nmap", "-sn", target, "-oG", "-").Output()
	if err != nil {
		return nil, err
	}

	s := make(ScanResult, 0)
	matches := nampOgRegex.FindAllSubmatch(out, -1)
	if matches != nil {
		for _, m := range matches {
			s = append(s, ScanItem{string(m[1]), string(m[2]), string(m[3]) == "Up"})
		}
	}

	return s, nil
}
