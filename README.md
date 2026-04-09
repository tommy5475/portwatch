# portwatch

A lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git
cd portwatch
go build -o portwatch
```

## Usage

Start monitoring with default settings (scans every 60 seconds):

```bash
portwatch start
```

Specify custom scan interval and ports to monitor:

```bash
portwatch start --interval 30s --ports 80,443,3000-3010
```

Run a one-time scan:

```bash
portwatch scan
```

View current status:

```bash
portwatch status
```

## Configuration

Create a `portwatch.yaml` in your config directory:

```yaml
interval: 60s
alert_on_new: true
alert_on_closed: true
ignore_ports: [22, 80, 443]
```

## Features

- 🔍 Real-time port monitoring
- 🔔 Alerts on new/closed ports
- ⚡ Lightweight and fast
- 🎯 Configurable port ranges
- 📝 JSON and text output formats

## License

MIT License - see [LICENSE](LICENSE) for details.