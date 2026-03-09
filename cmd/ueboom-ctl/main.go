package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/whayn/ueboom-ctl/internal/ble"
	"github.com/whayn/ueboom-ctl/internal/config"
	"github.com/whayn/ueboom-ctl/internal/logger"

	"github.com/spf13/pflag"
)

func init() {
	// Detect if we are in a container with a shared host D-Bus socket.
	hostDbus := "/run/host/run/dbus/system_bus_socket"
	if _, err := os.Stat(hostDbus); err == nil && os.Getenv("DBUS_SYSTEM_BUS_ADDRESS") == "" {
		os.Setenv("DBUS_SYSTEM_BUS_ADDRESS", "unix:path="+hostDbus)
	}
}

func main() {
	onFlag := pflag.Bool("on", false, "Power on the speaker")
	offFlag := pflag.Bool("off", false, "Power off the speaker")
	pflag.Parse()

	checkPrivileges()

	args := pflag.Args()
	if len(args) > 0 && args[0] == "setup" {
		runSetup()
		return
	}

	if *onFlag || *offFlag {
		runHeadless(*onFlag, *offFlag)
		return
	}

	showHelp()
}

func showHelp() {
	logger.Section("UE Boom BLE Control Agent v1.1")

	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	fmt.Fprintf(os.Stderr, "  %s %s [%s]\n\n", logger.Highlight("ueboom-ctl"), logger.Highlight("[command]"), logger.Highlight("--flags"))

	fmt.Fprintf(os.Stderr, "Commands:\n")
	fmt.Fprintf(os.Stderr, "  %-15s %s\n", logger.Highlight("setup"), "Run the guided setup to pair your speaker")

	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	pflag.CommandLine.SetOutput(os.Stderr)
	pflag.VisitAll(func(f *pflag.Flag) {
		fmt.Fprintf(os.Stderr, "  --%-13s %s\n", logger.Highlight(f.Name), f.Usage)
	})

	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  ueboom-ctl %s\n", logger.Highlight("setup"))
	fmt.Fprintf(os.Stderr, "  ueboom-ctl %s\n", logger.Highlight("--on"))
	fmt.Fprintf(os.Stderr, "  ueboom-ctl %s\n", logger.Highlight("--off"))
}

func checkPrivileges() {
	if os.Geteuid() != 0 {
		logger.Error("Root privileges required for BLE operations.")
		logger.Metadata("Please run with 'sudo' or set capabilities on the binary:")
		logger.Metadata("  sudo setcap 'cap_net_raw,cap_net_admin+eip' %s", os.Args[0])
		os.Exit(1)
	}
}

func runSetup() {
	logger.Info("Starting interactive setup")

	speakers, err := ble.Scan(10 * time.Second)
	if err != nil {
		if strings.Contains(err.Error(), "system_bus_socket") || strings.Contains(err.Error(), "connect: no such file") {
			logger.Error("Could not connect to D-Bus. Is the Bluetooth service (BlueZ) running?")
			logger.Metadata("Attempted D-Bus address: %s", os.Getenv("DBUS_SYSTEM_BUS_ADDRESS"))
			os.Exit(1)
		}
		logger.Error("Scan failed: %v", err)
		os.Exit(1)
	}

	if len(speakers) == 0 {
		logger.Warn("No UE Boom speakers found. Check speaker state and range.")
		return
	}

	logger.Section("Discovered UE Boom Speakers")
	for i, s := range speakers {
		logger.List(i+1, s.Name, s.Address)
	}

	reader := bufio.NewReader(os.Stdin)
	var target *ble.DiscoveredSpeaker

	for attempts := 0; attempts < 3; attempts++ {
		fmt.Fprintf(os.Stderr, "\nSelect a speaker [1-%d] (or 'q' to quit): ", len(speakers))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "q" || input == "quit" || input == "exit" {
			logger.Info("Setup cancelled by user.")
			return
		}

		choice, err := strconv.Atoi(input)
		if err == nil && choice >= 1 && choice <= len(speakers) {
			target = &speakers[choice-1]
			break
		}

		if attempts < 2 {
			logger.Warn("Invalid selection. Please try again.")
		} else {
			logger.Error("Too many invalid attempts.")
			return
		}
	}

	if target == nil {
		return
	}

	logger.Info("Selected: %s", target.Name)
	logger.Info("Retrieving host Bluetooth MAC address")
	hostMAC, err := ble.GetLocalMAC()
	if err != nil {
		logger.Error("Failed to get host MAC: %v", err)
		os.Exit(1)
	}

	formattedHostMAC := strings.ReplaceAll(hostMAC, ":", "")
	cfg := &config.Config{
		TargetMAC: target.Address,
		HostMAC:   formattedHostMAC,
	}

	savePath := config.SystemConfigPath
	if err := cfg.Save(savePath); err != nil {
		userPath, _ := config.UserConfigPath()
		logger.Warn("Could not save to system path (%s). Using user path.", savePath)
		if err := cfg.Save(userPath); err != nil {
			logger.Error("Failed to save configuration: %v", err)
			os.Exit(1)
		}
		savePath = userPath
	}

	logger.Success("Setup complete! Configuration saved to %s", savePath)
}

func runHeadless(on, off bool) {
	cfg, _, err := config.Load()
	if err != nil {
		logger.Error("Configuration missing. Run 'ueboom-ctl setup' first.")
		os.Exit(2)
	}

	if err := ble.SendPowerCommand(cfg.TargetMAC, cfg.HostMAC, on); err != nil {
		if strings.Contains(err.Error(), "system_bus_socket") || strings.Contains(err.Error(), "connect: no such file") {
			logger.Error("Could not connect to D-Bus. Is the Bluetooth service (BlueZ) running?")
			logger.Info("Attempted D-Bus address: %s", os.Getenv("DBUS_SYSTEM_BUS_ADDRESS"))
		} else {
			logger.Error("Command failed: %v", err)
		}
		os.Exit(1)
	}

	logger.Success("Command delivered successfully")
}
