# How I reverse-engineered the UE Boom protocol

I wanted to turn my UE Boom on from my computer without having to reach for the physical button or the buggy official app. Here’s how I got it working.

## Sniffing the traffic
The first step was seeing what the official app actually sends. I enabled "Bluetooth HCI Snoop Log" on my Android phone, toggled the power in the app a few times, and moved the logs over to my PC.

Filtering the logs in Wireshark for `btgatt` revealed a 7-byte payload being sent to UUID `c6d6dc0d-07f5-47ef-9b59-630622b01fd3`. 

The format is dead simple: `[Your Phone's MAC] + [01 or 02]`.
- `01` turns it on.
- `02` turns it off.

## Testing it manually
Before writing any code, I tried it with `gatttool`:
```bash
gatttool -b [SPEAKER_MAC] --char-write-req --handle=0x0003 --value=[MY_MAC]01
```
The speaker woke up immediately. Note the `--char-write-req` flag—that turned out to be the most important part.

## The Go struggle
I tried using standard Go BLE libraries, but they all failed. It turns out that on Linux, most of these libraries only implement "Write Commands" (Write Without Response). UE Boom speakers are picky and will only listen to a "Write Request."

Since the libraries didn't support what I needed, I had to bypass them. I wrote a manual D-Bus caller to talk to BlueZ directly. By passing the `type: request` flag in the D-Bus call, I was finally able to mimic what the official app (and `gatttool`) does.

## Final thoughts
The "pairing" is basically non-existent; the speaker just checks if the MAC address in the packet matches the device sending it. It's a nice, simple implementation once you get past the Linux Bluetooth stack headaches.
