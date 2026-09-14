<div align="center">

<a href="https://github.com/nguyenphanno/Discord-Username-Generator-Checker">
  <img src="https://raw.githubusercontent.com/nguyenphanno/Discord-Username-Generator-Checker/main/image/banner.png" width="100%" alt="Discord Username Generator & Checker" />
</a>

<br>

# Discord Username Generator & Checker

### High-performance · Concurrent · Proxy-powered

<p>
  <img src="https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=flat-square&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Windows-supported-18181B?style=flat-square&logo=windows&logoColor=white" />
  <img src="https://img.shields.io/badge/Linux-supported-18181B?style=flat-square&logo=linux&logoColor=white" />
  <img src="https://img.shields.io/badge/macOS-supported-18181B?style=flat-square&logo=apple&logoColor=white" />
</p>

<p>
  <img src="https://img.shields.io/github/stars/nguyenphanno/Discord-Username-Generator-Checker?style=flat-square&logo=github&color=e3b341" />
  <img src="https://img.shields.io/github/forks/nguyenphanno/Discord-Username-Generator-Checker?style=flat-square&logo=github&color=5865F2" />
  <img src="https://img.shields.io/github/issues/nguyenphanno/Discord-Username-Generator-Checker?style=flat-square&color=EF4444" />
  <img src="https://img.shields.io/github/license/nguyenphanno/Discord-Username-Generator-Checker?style=flat-square&color=22C55E" />
</p>

<br>

<a href="https://github.com/nguyenphanno/Discord-Username-Generator-Checker">
  <img src="https://img.shields.io/badge/View%20Repository-18181B?style=for-the-badge&logo=github&logoColor=white" />
</a>

<a href="https://github.com/nguyenphanno/Discord-Username-Generator-Checker/issues">
  <img src="https://img.shields.io/badge/Report%20Issue-EF4444?style=for-the-badge&logo=github&logoColor=white" />
</a>

</div>

<br>

<div align="center">

> A concurrent Go utility for generating usernames, checking availability,
> managing proxies, and organizing results.

<br>

`Fast` · `Concurrent` · `Configurable` · `Proxy-aware`

</div>

---

## Overview

**Discord Username Generator & Checker** is a terminal-based utility written in Go, designed around concurrent processing and configurable request workflows.

The project provides:

- Concurrent username processing
- Username generation
- Availability checking
- HTTP / SOCKS5 proxy support
- Proxy health checking
- Proxy scoring and prioritization
- Automatic resume
- Session logging
- Structured result files
- Interactive configuration

The goal is to keep the workflow simple while providing enough control for larger checking sessions.

---

## Features

<table>
<tr>
<td width="50%" valign="top">

### ⚡ Processing

- Concurrent workers
- Configurable thread count
- Configurable request delay
- Live generation
- File-based checking
- Generate-only mode
- Combined generation + checking

</td>

<td width="50%" valign="top">

### 🌐 Proxy Management

- HTTP proxies
- SOCKS5 proxies
- Proxy validation
- Health checking
- Performance scoring
- Automatic disabling
- Working proxy export

</td>
</tr>

<tr>
<td width="50%" valign="top">

### 💾 Session Management

- Automatic resume
- Duplicate prevention
- Session directories
- Persistent results
- Structured logs
- Runtime statistics

</td>

<td width="50%" valign="top">

### 📦 Output

- Available usernames
- Taken usernames
- Checked usernames
- Working proxies
- Session logs
- Organized result files

</td>
</tr>
</table>

---

## Preview

<div align="center">

<img
  src="https://raw.githubusercontent.com/nguyenphanno/Discord-Username-Generator-Checker/main/image/nguyenphanno.png"
  width="92%"
  alt="Application Preview"
/>
<br><br>
<img
  src="https://raw.githubusercontent.com/nguyenphanno/Discord-Username-Generator-Checker/main/image/logs.png"
  width="92%"
  alt="Logs Preview"
/>

</div>

---

## Modes

| Mode | Description |
| :--- | :--- |
| `live` | Continuously generate usernames and check availability |
| `check` | Read usernames from `usernames.txt` and process them |
| `generate` | Generate usernames without checking availability |
| `both` | Generate usernames and immediately check them |

### Example

```text
Mode: live
Threads: 30
Delay: 0.10s
Target: 5000
Length: 3-8
Quiet Mode: enabled
Proxy Check: enabled
````

---

## Architecture

```text
                         ┌──────────────────────┐
                         │   Interactive CLI    │
                         └──────────┬───────────┘
                                    │
                                    ▼
                         ┌──────────────────────┐
                         │    Configuration     │
                         └──────────┬───────────┘
                                    │
                   ┌────────────────┴────────────────┐
                   │                                 │
                   ▼                                 ▼
          ┌──────────────────┐             ┌──────────────────┐
          │ Username Engine  │             │  Proxy Manager   │
          └────────┬─────────┘             └────────┬─────────┘
                   │                                │
                   └────────────────┬───────────────┘
                                    │
                                    ▼
                         ┌──────────────────────┐
                         │  Concurrent Workers  │
                         └──────────┬───────────┘
                                    │
                                    ▼
                         ┌──────────────────────┐
                         │ Availability Check  │
                         └──────────┬───────────┘
                                    │
                         ┌──────────┴──────────┐
                         │                     │
                         ▼                     ▼
                  ┌────────────┐        ┌────────────┐
                  │  Results   │        │    Logs    │
                  └────────────┘        └────────────┘
```

---

## Proxy System

The proxy layer automatically validates and manages available endpoints.

```text
Load Proxies
     │
     ▼
Health Check
     │
     ▼
Score Proxies
     │
     ▼
Prioritize
     │
     ▼
Assign Worker
     │
     ▼
Monitor
     │
     ├───────────────┐
     │               │
     ▼               ▼
  Healthy        Unhealthy
     │               │
     ▼               ▼
Keep Active       Disable
```

### Supported

```text
HTTP
SOCKS5
```

### Proxy formats

```text
http://ip:port
socks5://ip:port
```

Authenticated proxies:

```text
http://username:password@ip:port
socks5://username:password@ip:port
```

Working proxies can be exported to:

```text
working_proxies.txt
```

---

## Browser Headers

Requests can use configurable browser-style request properties such as:

```text
User-Agent
Super-Properties
Timezone
Locale
```

These values are intended to provide a consistent request profile during processing.

---

## Automatic Resume

Processed usernames can be preserved between sessions.

Typical states:

```text
AVAILABLE
TAKEN
CHECKED
```

This allows subsequent sessions to skip entries that have already been processed.

---

## Output

### Result files

```text
available.txt
taken.txt
checked.txt
working_proxies.txt
```

### Session structure

```text
results/
│
├── session-2026-09-11/
│   ├── available.txt
│   ├── taken.txt
│   ├── checked.txt
│   └── session.log
│
└── sniper.log
```

The exact structure may vary depending on the current configuration.

---

## Logging

Runtime information can include:

```text
Timestamp
Username
Status
Proxy
Request Result
Worker
```

Example:

```text
[22:04:39] INFO  Proxies: 8023
[22:04:52] INFO  Self-testing Discord API...
[22:05:14] WARN  Self-test failed
```

### Smart Quiet Mode

When enabled, console output is reduced to the most relevant results.

```text
AVAILABLE
TAKEN
```

This keeps long-running sessions significantly cleaner.

---

## Installation

### Requirements

* Go `1.22+`
* Internet connection
* Optional HTTP / SOCKS5 proxies

### Clone

```bash
git clone https://github.com/nguyenphanno/Discord-Username-Generator-Checker.git
cd Discord-Username-Generator-Checker
```

### Dependencies

```bash
go mod tidy
```

### Run

```bash
go run main.go
```

### Build

Linux / macOS:

```bash
go build -o discord-username-checker
```

Windows:

```powershell
go build -o discord-username-checker.exe
```

---

## Configuration

The application uses an interactive terminal setup.

Available options include:

```text
Mode
Threads
Request Delay
Target Amount
Username Minimum Length
Username Maximum Length
Proxy Health Check
Proxy Validation
Smart Quiet Mode
```

Example:

```text
Mode (live/check/generate/both) [live]: live
Threads [35]: 80
Delay (seconds) [0.08]: 0.08
Target usernames [5000]: 5000
Username min length [3]: 3
Username max length [8]: 8
Health-check proxies first? [y]: y
Smart Quiet [y]: y
```

---

## Performance

Processing speed depends on:

* Worker count
* Network latency
* Proxy quality
* Request delay
* Endpoint response time
* Available system resources

A higher thread count does not always mean better performance.

### Balanced configuration

```text
Threads      25 - 40
Delay        0.05 - 0.15s
Quiet Mode   Enabled
Proxy Check  Enabled
```

For larger sessions, tune concurrency gradually rather than immediately using the highest possible value.

---

## Project Structure

```text
Discord-Username-Generator-Checker/
│
├── .github/
│   └── workflows/
│
├── image/
│   ├── banner.png
│   ├── nguyenphanno.png
│   └── logs.png
│
├── results/
│
├── main.go
├── go.mod
├── go.sum
│
├── proxies.txt
├── usernames.txt
│
├── available.txt
├── taken.txt
├── checked.txt
├── working_proxies.txt
│
└── README.md
```

---

## Platform Support

<div align="center">

| Platform |    Status   |
| :------: | :---------: |
|  Windows | ✓ Supported |
|   Linux  | ✓ Supported |
|   macOS  | ✓ Supported |

</div>

---

## Roadmap

| Feature                  |  Status |
| :----------------------- | :-----: |
| Concurrent workers       |    ✓    |
| HTTP proxy support       |    ✓    |
| SOCKS5 proxy support     |    ✓    |
| Proxy health checking    |    ✓    |
| Proxy scoring            |    ✓    |
| Automatic resume         |    ✓    |
| Session logging          |    ✓    |
| Result exporting         |    ✓    |
| Improved terminal UI     | Planned |
| Extended statistics      | Planned |
| Additional configuration | Planned |

---

## Disclaimer

This project is provided for educational and research purposes.

Users are responsible for ensuring that their use of the software complies with Discord's Terms of Service, applicable API policies, and local laws.

Do not use the software to abuse services, bypass access controls, collect credentials, or perform unauthorized activity.

Use reasonable request rates, delays, and proxy configurations.

---

# Author

<div align="center">

<a href="https://github.com/nguyenphanno">
  <img
    src="https://github.com/nguyenphanno.png"
    width="120"
    height="120"
    alt="nguyenphanno"
  />
</a>

<br><br>

## `nguyenphanno`

<br>

<a href="https://github.com/nguyenphanno">
  <img src="https://img.shields.io/badge/GitHub-nguyenphanno-18181B?style=for-the-badge&logo=github&logoColor=white" />
</a>

<a href="https://discord.com/users/1469216989158309928">
  <img src="https://img.shields.io/badge/Discord-Presence-5865F2?style=for-the-badge&logo=discord&logoColor=white" />
</a>

<br><br>

<a href="https://discord.com/users/1469216989158309928">
  <img
    src="https://lanyard.cnrad.dev/api/1469216989158309928?theme=light&bg=ffffff&borderRadius=18px&animated=true&showDisplayName=true&hideDiscrim=false&hideStatus=false&hideTimestamp=false&idleMessage=somewhere%20in%20Tokyo%2C%20thinking%20fast..."
    alt="Discord Presence"
  />
</a>

<br><br>

<img
src="https://github-profile-summary-cards.vercel.app/api/cards/profile-details?username=nguyenphanno&theme=default"
width="850"
alt="GitHub Profile"
/>

</div>

---

<div align="center">

<a href="https://github.com/nguyenphanno/Discord-Username-Generator-Checker">
  <img
    src="https://img.shields.io/github/stars/nguyenphanno/Discord-Username-Generator-Checker?style=for-the-badge&logo=github&logoColor=white&label=Star"
    alt="Star Repository"
  />
</a>

<a href="https://github.com/nguyenphanno/Discord-Username-Generator-Checker/fork">
  <img
    src="https://img.shields.io/badge/Fork-Repository-5865F2?style=for-the-badge&logo=github&logoColor=white"
    alt="Fork Repository"
  />
</a>

<br><br>

<img
src="https://img.shields.io/github/last-commit/nguyenphanno/Discord-Username-Generator-Checker?style=flat-square&color=8B5CF6"
alt="Last Commit"
/>

<br><br>

<sub>Built with Go · maintained by nguyenphanno</sub>

</div>
