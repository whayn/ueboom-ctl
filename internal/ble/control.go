package ble

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/godbus/dbus/v5"
	"tinygo.org/x/bluetooth"
)

const (
	StandbyServiceUUID        = "c6d6dc0d-07f5-47ef-9b59-630622b01fd3"
	StandbyCharacteristicUUID = "c6d6dc0d-07f5-47ef-9b59-630622b01fd3"
	Attempts                  = 3
)

func SendPowerCommand(targetAddr string, hostMAC string, on bool) error {
	if err := Adapter.Enable(); err != nil {
		return fmt.Errorf("failed to enable adapter: %w", err)
	}

	mac, err := bluetooth.ParseMAC(targetAddr)
	if err != nil {
		return fmt.Errorf("invalid target address: %w", err)
	}
	addr := bluetooth.Address{MACAddress: bluetooth.MACAddress{MAC: mac}}

	cmd := "02" // OFF
	if on {
		cmd = "01" // ON
	}

	payload, err := hex.DecodeString(hostMAC + cmd)
	if err != nil {
		return fmt.Errorf("failed to encode payload: %w", err)
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Writer = os.Stderr
	s.Suffix = " Delivering power command..."
	s.Start()
	defer s.Stop()

	var lastErr error
	for i := 0; i < Attempts; i++ {
		if i > 0 {
			s.Suffix = fmt.Sprintf(" Retrying command (attempt %d/%d)...", i+1, Attempts)
			time.Sleep(500 * time.Millisecond)
		}

		err = connectAndWriteRaw(addr, payload)
		if err == nil {
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("failed after 3 attempts: %v", lastErr)
}

// connectAndWriteRaw uses direct D-Bus calls to BlueZ to perform a "Write Request",
// bypassing the limitations of the tinygo-org/bluetooth library on Linux.
func connectAndWriteRaw(addr bluetooth.Address, payload []byte) error {
	device, err := Adapter.Connect(addr, bluetooth.ConnectionParams{})
	if err != nil {
		return err
	}
	defer device.Disconnect()

	// Wait a moment for services to be resolved by BlueZ
	time.Sleep(500 * time.Millisecond)

	// Connect to System D-Bus
	conn, err := dbus.SystemBus()
	if err != nil {
		return fmt.Errorf("failed to connect to system bus: %w", err)
	}

	// We need to find the object path for the characteristic.
	// BlueZ paths follow a pattern: /org/bluez/hciX/dev_XX_XX_XX_XX_XX_XX/serviceXXXX/charXXXX
	// Instead of guessing, we'll query all managed objects.
	var objects map[dbus.ObjectPath]map[string]map[string]dbus.Variant
	err = conn.Object("org.bluez", "/").Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objects)
	if err != nil {
		return fmt.Errorf("failed to get managed objects: %w", err)
	}

	var charPath dbus.ObjectPath
	targetAddrPath := "dev_" + strings.ReplaceAll(addr.String(), ":", "_")

	for path, interfaces := range objects {
		// Filter by device address and the characteristic interface
		if !strings.Contains(string(path), targetAddrPath) {
			continue
		}

		if props, ok := interfaces["org.bluez.GattCharacteristic1"]; ok {
			if uuid, ok := props["UUID"]; ok && strings.EqualFold(uuid.Value().(string), StandbyCharacteristicUUID) {
				charPath = path
				break
			}
		}
	}

	if charPath == "" {
		return fmt.Errorf("standby characteristic not found in BlueZ via D-Bus")
	}

	// Perform the actual Write Request (WriteValue with type=request)
	// This is equivalent to gatttool --char-write-req
	obj := conn.Object("org.bluez", charPath)
	options := map[string]dbus.Variant{
		"type": dbus.MakeVariant("request"),
	}

	call := obj.Call("org.bluez.GattCharacteristic1.WriteValue", 0, payload, options)
	if call.Err != nil {
		return fmt.Errorf("D-Bus WriteValue failed: %w", call.Err)
	}

	return nil
}
