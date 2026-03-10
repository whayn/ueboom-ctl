# Controlling UE Boom from Home Assistant (HAOS/Supervised)

If you're running Home Assistant OS or Supervised, you're dealing with a bunch of Docker containers. Since `ueboom-ctl` needs direct access to your Bluetooth hardware, getting it to work from the "Core" container is a bit of a hack, but it works great. We're essentially going to run the command inside the SSH add-on because that's where the hardware access is.

## What you'll need
- **HAOS** or **Supervised** installed.
- The **SSH & Web Terminal** add-on.
- A working **Bluetooth adapter** (USB or built-in).

## 1. Sorting out permissions

The SSH add-on needs to be able to talk to your hardware. 
1.  Go to **Settings > Add-ons > SSH & Web Terminal**.
2.  Turn **Protection mode** to **OFF**. If this is on, the script won't be able to see your Bluetooth adapter.
3.  **Restart** the add-on so the change takes effect.

## 2. Placing the binary

You should put the `ueboom-ctl` binary in your `/config` folder. This is the easiest way to make sure both the Core container and the SSH add-on can see it.

```bash
# Open your HA Terminal
mkdir -p /config/ueboom
# Move the binary there and make it executable
chmod +x /config/ueboom/ueboom-ctl
cd /config/ueboom
```

## 3. Pairing the speaker

We're going to pair the speaker through the CLI. 

1.  Turn the UE Boom **ON**.
2.  Hold the Bluetooth button until it makes that drum sound and the light starts blinking.
3.  Run the setup:
    ```bash
    ./ueboom-ctl setup
    ```
    Pick your speaker from the list. The tool will save your config to `/etc/ueboom/config.json` (this lives inside the add-on container).

### Getting a "br-connection-unknown" error?
Bluetooth on Linux is finicky. If it fails to connect, try manually "trusting" the speaker in the host's manager:
1.  Type `bluetoothctl` and hit Enter.
2.  Run `trust [YOUR_SPEAKER_MAC]`.
3.  Maybe run `pair [YOUR_SPEAKER_MAC]` if it's still being difficult.
4.  Type `exit` to get out.

## 4. The Home Assistant integration

### Setting up the switch
We'll use a `command_line` switch. Because the binary runs in the SSH add-on, the Core container has to SSH into the add-on to actually fire the command.

*Note: You'll need to set up SSH keys between Core and the add-on for this to work without asking for a password.*

```yaml
# configuration.yaml
switch:
  - platform: command_line
    switches:
      ue_boom_power:
        friendly_name: "UE Boom"
        # We SSH from Core into the add-on container to run the command
        command_on: "ssh -i /config/.ssh/id_rsa -o StrictHostKeyChecking=no root@core-ssh '/config/ueboom/ueboom-ctl --on'"
        command_off: "ssh -i /config/.ssh/id_rsa -o StrictHostKeyChecking=no root@core-ssh '/config/ueboom/ueboom-ctl --off'"
        icon_template: >
          {% if states('switch.ue_boom_power') == 'on' %}
            mdi:speaker-bluetooth
          {% else %}
            mdi:speaker-off
          {% endif %}
```

### Automation: Wake on music
This turns the speaker on automatically when your media player starts playing.

```yaml
# automations.yaml
- alias: "Auto-wake UE Boom"
  trigger:
    - platform: state
      entity_id: media_player.spotify # or whatever you use
      to: 'playing'
  action:
    - service: switch.turn_on
      target:
        entity_id: switch.ue_boom_power
```

## 5. Dashboard button
Add a simple button to your UI:
```yaml
type: button
entity: switch.ue_boom_power
name: Wake Boom
icon: mdi:speaker-bluetooth
```
