# UE Boom BLE Control

[![CI](https://github.com/whayn/ueboom-ctl/actions/workflows/ci.yml/badge.svg?style=flat)](https://github.com/whayn/ueboom-ctl/actions/workflows/ci.yml)
[![Release](https://github.com/whayn/ueboom-ctl/actions/workflows/release.yml/badge.svg?style=flat)](https://github.com/whayn/ueboom-ctl/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/whayn/ueboom-ctl?style=flat)](https://goreportcard.com/report/github.com/whayn/ueboom-ctl)

A small Linux tool to turn UE Boom speakers on and off over Bluetooth.

## Quick Start

### Install
```bash
go install github.com/whayn/ueboom-ctl/cmd/ueboom-ctl@latest
```

### Set it up
You only need to do this once to find your speaker and save the config.
```bash
sudo ueboom-ctl setup
```

### Use it
```bash
sudo ueboom-ctl --on
sudo ueboom-ctl --off
```

### Build from source
For other methods or to build locally:
```bash
git clone https://github.com/whayn/ueboom-ctl.git
cd ueboom-ctl
make build
```
This requires `make` and `go` to be installed.

## Home Assistant
For a full guide on setting this up with **HAOS**(including bypassing container permissions), see the [Home Assistant Setup Guide](./docs/home-assistant-setup.md).

Alternatively, add this to your `configuration.yaml` for direct desktop/server control:
```yaml
switch:
  - platform: command_line
    switches:
      ue_boom_speaker:
        command_on: "sudo /path/to/ueboom-ctl --on"
        command_off: "sudo /path/to/ueboom-ctl --off"
        friendly_name: "UE Boom Speaker"
```

## Troubleshooting & Permissions
Linux is picky about Bluetooth permissions. See the [Setup Guide](./docs/home-assistant-setup.md) for detailed HA troubleshooting, or run this with `sudo` locally:
```bash
sudo setcap 'cap_net_raw,cap_net_admin+eip' ./ueboom-ctl
```

## The Details
I wrote a [full technical write-up](./docs/write-up.md) on how I reverse-engineered the protocol by sniffing Android app traffic if you're interested in the internals.

## License
[GPL-3.0](LICENSE)
