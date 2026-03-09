# UE Boom BLE Control

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
Add this to your `configuration.yaml` to control your speaker from your dashboard:

```yaml
switch:
  - platform: command_line
    switches:
      ue_boom_speaker:
        command_on: "sudo /path/to/ueboom-ctl --on"
        command_off: "sudo /path/to/ueboom-ctl --off"
        friendly_name: "UE Boom Speaker"
```
*Note: Make sure your Home Assistant user can run this via sudo or set the capabilities mentioned below.*

## Troubleshooting & Permissions
Linux is picky about Bluetooth permissions. You'll need to run this with `sudo` or set capabilities on the binary:
```bash
sudo setcap 'cap_net_raw,cap_net_admin+eip' ./ueboom-ctl
```

If it still doesn't work:
- Check `systemctl status bluetooth` to make sure BlueZ is running.
- Ensure the speaker is in standby (plugged in or recently used), not completely dead.
- Use `rfkill list` to make sure your adapter isn't blocked.

## The Details
I wrote a [full write-up](./write-up.md) on how I reverse-engineered the protocol by sniffing Android app traffic if you're interested in the internals.

## License
GPL-3.0
