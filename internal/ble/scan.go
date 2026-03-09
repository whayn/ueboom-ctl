package ble

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"tinygo.org/x/bluetooth"
)

var (
	Adapter = bluetooth.DefaultAdapter
)

type DiscoveredSpeaker struct {
	Name    string
	Address string
}

func Scan(timeout time.Duration) ([]DiscoveredSpeaker, error) {
	if err := Adapter.Enable(); err != nil {
		return nil, fmt.Errorf("failed to enable bluetooth adapter: %w", err)
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Writer = os.Stderr
	s.Suffix = " Searching for UE Boom speakers..."
	s.Start()
	defer s.Stop()

	var discovered []DiscoveredSpeaker
	seen := make(map[string]bool)

	// Stop scanning after timeout
	timer := time.AfterFunc(timeout, func() {
		Adapter.StopScan()
	})
	defer timer.Stop()

	err := Adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		addr := result.Address.String()
		if seen[addr] {
			return
		}

		// Check for UE OUI
		isUE := strings.HasPrefix(strings.ToLower(addr), "10:94:97")
		
		// Check for UE service UUID
		for _, uuid := range result.AdvertisementPayload.ServiceUUIDs() {
			if strings.Contains(strings.ToLower(uuid.String()), "c6d6dc0d") {
				isUE = true
				break
			}
		}

		if isUE {
			name := result.LocalName()
			if name == "" {
				name = "Unknown UE Speaker"
			}
			discovered = append(discovered, DiscoveredSpeaker{
				Name:    name,
				Address: addr,
			})
			seen[addr] = true
			s.Suffix = fmt.Sprintf(" Searching... Discovered: %s (%d total)", name, len(discovered))
		}
	})

	if err != nil && !strings.Contains(err.Error(), "scan stopped") {
		return nil, err
	}

	return discovered, nil
}

// GetLocalMAC retrieves the MAC address of the local Bluetooth adapter.
func GetLocalMAC() (string, error) {
	if err := Adapter.Enable(); err != nil {
		return "", fmt.Errorf("failed to enable adapter: %w", err)
	}
	
	addr, err := Adapter.Address()
	if err != nil {
		return "", err
	}
	
	return addr.String(), nil
}
