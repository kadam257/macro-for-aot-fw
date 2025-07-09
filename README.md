# Roblox Shift Lock Macro

A macro tool specifically designed with aot freedom war on roblox in mind, designed for shift lock players and built for linux (sorry if anyone on windows wants to use)

## Features

- **Smart Shift Lock Cycling**: Toggle between shift lock modes with context-aware behavior
- **Target Enemy Support**: Enable/disable target enemy mode while maintaining shift lock
- **Manual Toggle**: Enable/disable macros with a hotkey for precise control

## Controls

- `=` - Toggle macros ON/OFF
- `[` - Cycle through shift lock modes:
  1. Enter shift lock (no targeting)
  2. Enable target enemies (stay in shift lock)
  3. Disable target enemies (stay in shift lock)
- `l` - Exit shift lock mode entirely
- `Ctrl+C` - Exit program

## Installation

### Prerequisites

Install required system packages:

```bash
# Fedora/RHEL
sudo dnf install libX11-devel libXext-devel libXtst-devel libXinerama-devel libXrandr-devel libXScrnSaver-devel libxkbcommon-devel libxkbcommon-x11-devel

# Ubuntu/Debian
sudo apt install libx11-dev libxext-dev libxtst-dev libxinerama-dev libxrandr-dev libxss-dev libxkbcommon-dev libxkbcommon-x11-dev
```

### Build

```bash
go mod tidy
go build -o macro
```

## Usage

1. Run the program with sudo (required for system-wide key capture):
   ```bash
   sudo ./macro
   ```

2. The program starts with macros **disabled**

3. When ready to use in your game:
   - Press `=` to enable macros
   - Use `[` and `l` keys as needed
   - Press `=` again to disable when done

## How It Works

The macro uses the double-shift technique to maintain shift lock while enabling target enemy mode:
- Single shift tap: Toggles shift lock on/off
- Double shift (tap + hold): Exits and re-enters shift lock while holding targeting

## Requirements

- Linux (tested on Fedora)
- Go 1.19+
- Root privileges (for system-wide key capture)
- X11 display server

## Notes

- Designed for Roblox games running through Sober on Linux
- Macros only active when manually enabled
- All key combinations are configurable in the source code

## License

MIT License
