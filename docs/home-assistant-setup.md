# Home Assistant Setup Guide (HAOS/Supervised)

This guide explains how to integrate `ueboom-ctl` into Home Assistant to control your Ultimate Ears Boom speaker directly from your dashboard and automations.

## 1. Prerequisites
*   **Home Assistant OS (HAOS)** or **Supervised** installation.
*   **SSH & Web Terminal** add-on installed (from the official or community store).
*   **Bluetooth hardware** recognized by your host (USB dongle or built-in).

## 2. Initial Setup & Hardware Permissions

### Disable Protection Mode
For the tool to access your Bluetooth hardware from the terminal, the add-on needs high privileges:
1.  Go to **Settings > Add-ons > SSH & Web Terminal**.
2.  Switch **Protection mode** to **OFF**.
3.  **Restart** the add-on.

### Install & Prepare the Binary
Place the `ueboom-ctl` binary into your `/config` directory (which is shared between the Core and the Add-on).
```bash
# In the HA Terminal
mkdir -p /config/ueboom
# Move your binary there and make it executable
chmod +x /config/ueboom/ueboom-ctl
cd /config/ueboom
```

## 3. Pairing the Speaker

### Put Speaker in Pairing Mode
1.  Turn your UE Boom **ON**.
2.  Hold the Bluetooth button until you hear the drum sound and the light blinks rapidly.

### Run the Setup Command
```bash
./ueboom-ctl setup
```
Follow the prompts to select your speaker. This saves the configuration to `/etc/ueboom/config.json` (inside the add-on container).

### Troubleshooting: `br-connection-unknown`
If the command fails with a connection error, you may need to manually "trust" the device in the host's Bluetooth manager:
1.  Type `bluetoothctl` and press Enter.
2.  Run `trust [YOUR_SPEAKER_MAC]`.
3.  (Optional) Run `pair [YOUR_SPEAKER_MAC]`.
4.  Type `exit`.

## 4. Home Assistant Integration

### Create a Command Line Switch
Add this to your `configuration.yaml`. This uses SSH to talk from the Core container to the SSH Add-on where the binary has hardware access.

*Note: You will need to set up SSH keys between the Core container and the SSH add-on for passwordless execution.*

```yaml
# configuration.yaml
switch:
  - platform: command_line
    switches:
      ue_boom_power:
        friendly_name: "UE Boom Speaker"
        command_on: "ssh -i /config/.ssh/id_rsa -o StrictHostKeyChecking=no root@core-ssh '/config/ueboom/ueboom-ctl --on'"
        command_off: "ssh -i /config/.ssh/id_rsa -o StrictHostKeyChecking=no root@core-ssh '/config/ueboom/ueboom-ctl --off'"
        icon_template: >
          {% if states('switch.ue_boom_power') == 'on' %}
            mdi:speaker-bluetooth
          {% else %}
            mdi:speaker-off
          {% endif %}
```

### Automation Example: Auto-On with Music
```yaml
# automations.yaml
- alias: "Speaker: Auto Turn On UE Boom"
  trigger:
    - platform: state
      entity_id: media_player.your_player
      to: 'playing'
  action:
    - service: switch.turn_on
      target:
        entity_id: switch.ue_boom_power
```

## 5. Dashboard (Lovelace)
A simple stack to control your speaker:
```yaml
type: vertical-stack
cards:
  - type: button
    entity: switch.ue_boom_power
    name: UE Boom
    show_state: true
    icon: mdi:speaker-bluetooth
```
