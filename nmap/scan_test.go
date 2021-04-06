package nmap

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_scan(t *testing.T) {
	t.Log("Test runs for 15 seconds")

	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)

	target := "192.168.5.0/24"

	var lastScan ScanResult

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(time.Second * 1):
			scan, err := Scan(ctx, target)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					panic(err)
				}
			}
			if lastScan == nil {
				fmt.Println(scan.String())
			} else {
				lastScan.Diff(scan)
			}
			lastScan = scan
		}
	}
}
