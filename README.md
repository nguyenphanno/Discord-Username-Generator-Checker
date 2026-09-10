<div align="center">

# `Discord Username Generator & Checker`

### High-performance · Concurrent · Proxy-powered

<p>
  <img src="https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.22+" />
  <img src="https://img.shields.io/badge/Windows-Linux-macOS-18181B?style=for-the-badge" alt="Platforms" />
  <img src="https://img.shields.io/badge/HTTP-SOCKS5-F97316?style=for-the-badge" alt="Proxy Support" />
</p>

<p>
  <img src="https://img.shields.io/github/stars/nguyenphanno/Discord-Username-Generator-Checker?style=for-the-badge&logo=github&logoColor=white&color=e3b341" alt="Stars" />
  <img src="https://img.shields.io/github/forks/nguyenphanno/Discord-Username-Generator-Checker?style=for-the-badge&logo=github&logoColor=white&color=5865F2" alt="Forks" />
  <img src="https://img.shields.io/github/issues/nguyenphanno/Discord-Username-Generator-Checker?style=for-the-badge&color=EF4444" alt="Issues" />
  <img src="https://img.shields.io/github/license/nguyenphanno/Discord-Username-Generator-Checker?style=for-the-badge&color=22C55E" alt="License" />
</p>

<br>

<img
src="https://count.getloli.com/@nguyenphanno?name=nguyenphanno&theme=booru-qualityhentais&padding=1&offset=0&align=center&scale=1&pixelated=0&darkmode=auto"
alt="Profile Views"
/>

<br><br>

<a href="https://github.com/nguyenphanno/Discord-Username-Generator-Checker">
  <img
    src="https://img.shields.io/badge/View%20Repository-18181B?style=for-the-badge&logo=github&logoColor=white"
    alt="View Repository"
  />
</a>

</div>

---

<div align="center">

**Generate usernames · Check availability · Manage proxies · Export results**

</div>

## Overview

Discord Username Generator & Checker is a high-performance Go utility designed for concurrent username generation and availability checking.

The application provides an interactive terminal interface, configurable concurrency, HTTP/SOCKS5 proxy support, automatic proxy health management, session persistence, and structured result logging.

---

## Preview

<div align="center">

<img
src="https://raw.githubusercontent.com/nguyenphanno/Discord-Username-Generator-Checker/main/image/nguyenphanno.png"
width="900"
alt="Application Preview"
/>
<img
src="https://raw.githubusercontent.com/nguyenphanno/Discord-Username-Generator-Checker/main/image/logs.png"
width="900"
alt="Logs Preview"
/>

</div>

---

## Features

| Feature                       | Description                                                               |
| ----------------------------- | ------------------------------------------------------------------------- |
| **Concurrent Processing**     | Process multiple username checks concurrently using configurable workers. |
| **Smart Quiet Mode**          | Reduce console output and display only relevant availability results.     |
| **HTTP / SOCKS5 Proxies**     | Supports both HTTP and SOCKS5 proxy endpoints.                            |
| **Proxy Health Checking**     | Validate proxies before and during a session.                             |
| **Proxy Scoring**             | Prioritize reliable proxies based on their performance.                   |
| **Automatic Proxy Disabling** | Temporarily remove unhealthy proxies from the active pool.                |
| **Automatic Resume**          | Previously processed usernames can be skipped automatically.              |
| **Session Logging**           | Keep results and activity organized by session.                           |
| **Request Logging**           | Record timestamps and proxy information for processed requests.           |
| **Browser Headers**           | Uses configurable browser-style request properties.                       |
| **Interactive Configuration** | Configure the session directly from the terminal.                         |
| **Result Export**             | Export available, taken, checked usernames and working proxies.           |

---

## Modes

### `live`

Continuously generate usernames and check their availability.

### `check`

Read usernames from:

```text
usernames.txt
```

and check them against the configured endpoint.

### `generate`

Generate usernames without performing availability checks.

### `both`

Generate usernames and immediately process them through the checker.

<div align="center">

<img src="https://img.shields.io/badge/LIVE-Continuous-8B5CF6?style=for-the-badge" alt="Live Mode" />
<img src="https://img.shields.io/badge/CHECK-File%20Based-5865F2?style=for-the-badge" alt="Check Mode" />
<img src="https://img.shields.io/badge/GENERATE-Generation-00ADD8?style=for-the-badge" alt="Generate Mode" />
<img src="https://img.shields.io/badge/BOTH-Combined-e3b341?style=for-the-badge" alt="Both Mode" />

</div>

---

## Architecture

The application follows a concurrent processing model:

```text
                   ┌─────────────────────┐
                   │   Interactive CLI   │
                   └──────────┬──────────┘
                              │
                              ▼
                   ┌─────────────────────┐
                   │    Configuration    │
                   └──────────┬──────────┘
                              │
               ┌──────────────┴──────────────┐
               │                             │
               ▼                             ▼
      ┌─────────────────┐          ┌─────────────────┐
      │ Username Engine │          │  Proxy Manager  │
      └────────┬────────┘          └────────┬────────┘
               │                            │
               └──────────────┬─────────────┘
                              ▼
                   ┌─────────────────────┐
                   │ Concurrent Workers  │
                   └──────────┬──────────┘
                              │
                              ▼
                   ┌─────────────────────┐
                   │ Availability Check  │
                   └──────────┬──────────┘
                              │
                    ┌─────────┴─────────┐
                    │                   │
                    ▼                   ▼
             ┌────────────┐      ┌────────────┐
             │  Results   │      │    Logs    │
             └────────────┘      └────────────┘
```

---

## Installation

### Requirements

* Go `1.22+`
* Internet connection
* Optional HTTP or SOCKS5 proxies

### Clone Repository

```bash
git clone https://github.com/nguyenphanno/Discord-Username-Generator-Checker.git
cd Discord-Username-Generator-Checker
```

### Install Dependencies

```bash
go mod tidy
```

### Configure Proxies

Create:

```text
proxies.txt
```

Supported formats:

```text
http://ip:port
```

```text
socks5://ip:port
```

Authentication proxies can use:

```text
http://username:password@ip:port
```

```text
socks5://username:password@ip:port
```

### Start

```bash
go run main.go
```

For a compiled binary:

```bash
go build -o discord-username-checker
```

Windows:

```powershell
go build -o discord-username-checker.exe
```

---

## Configuration

The application provides an interactive configuration flow before starting a session.

Available options include:

```text
Mode
Threads
Request Delay
Target Amount
Username Length
Smart Quiet Mode
Proxy Health Check
Proxy Validation
```

Example:

```text
Mode: live
Threads: 30
Delay: 0.10s
Target: 10000
Length: 4-5
Quiet Mode: enabled
Proxy Check: enabled
```

Configuration is applied per session and does not require manual editing of source code.

---

## Proxy System

The proxy manager is designed around automatic validation and prioritization.

```text
Load Proxies
     │
     ▼
Health Check
     │
     ▼
Calculate Score
     │
     ▼
Prioritize
     │
     ▼
Assign Worker
     │
     ▼
Monitor Performance
     │
     ├───────────────┐
     │               │
     ▼               ▼
 Healthy          Unhealthy
     │               │
     ▼               ▼
 Keep Active     Disable
```

### Supported Protocols

```text
HTTP
SOCKS5
```

### Proxy Management

The system can:

* Validate proxies before use
* Track proxy performance
* Prioritize reliable proxies
* Detect failed connections
* Disable unhealthy proxies
* Reuse healthy proxies
* Export working proxies

Working proxies can be written to:

```text
working_proxies.txt
```

---

## Browser Headers

Requests can include browser-style properties such as:

```text
User-Agent
Super-Properties
Timezone
Locale
```

These values are used to provide a more consistent request profile during processing.

---

## Automatic Resume

The application can preserve previously processed usernames.

This allows subsequent sessions to avoid unnecessary duplicate checks.

Typical result states include:

```text
Available
Taken
Checked
```

Previously processed entries can therefore be skipped depending on the current session configuration.

---

## Output

### Result Files

```text
available.txt
taken.txt
checked.txt
working_proxies.txt
```

### Session Directory

```text
results/
```

Example:

```text
results/
├── session-2026-09-11/
│   ├── available.txt
│   ├── taken.txt
│   ├── checked.txt
│   └── session.log
│
└── sniper.log
```

The exact output structure may vary depending on the current application configuration.

---

## Logging

Logs contain useful runtime information such as:

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
[05:42:18] [SUCCESS] Username available: example
[05:42:19] [INFO] Checked username: another
[05:42:20] [WARN] Proxy temporarily disabled
```

Quiet mode can reduce console output to essential status information:

```text
AVAILABLE
TAKEN
```

---

## Performance

The checker uses concurrent workers to process multiple requests in parallel.

Performance depends on:

* Number of workers
* Proxy quality
* Network latency
* Request delay
* Endpoint response time
* System resources

A balanced configuration is generally preferable to simply increasing concurrency.

### Recommended Balanced Profile

```text
Threads      25 - 40
Delay        0.05 - 0.15 seconds
Quiet Mode   Enabled
Proxy Check  Enabled
```

Actual performance may vary depending on network and proxy conditions.

---

## Project Structure

```text
Discord-Username-Generator-Checker/
│
├── image/
│   ├── nguyenphanno.png
│   └── logs.png
│
├── results/
│
├── main.go
├── go.mod
├── go.sum
├── proxies.txt
├── usernames.txt
├── available.txt
├── taken.txt
├── checked.txt
├── working_proxies.txt
└── README.md
```

---

## Platform Support

| Platform | Support   |
| -------- | --------- |
| Windows  | Supported |
| Linux    | Supported |
| macOS    | Supported |

Go's cross-platform build system makes it possible to compile the application for multiple operating systems.

---

## Disclaimer

This project is provided for educational and research purposes.

The application only checks username availability through publicly accessible endpoints and does not provide account access, credential collection, or account automation.

Users are responsible for complying with Discord's Terms of Service, applicable API policies, and local laws.

Avoid excessive request rates and use appropriate delays and proxy configurations.

---

## Author

<div align="center">

<img
src="https://github.com/nguyenphanno.png"
width="120"
height="120"
alt="nguyenphanno"
/>

<br><br>

### nguyenphanno

<br>

<a href="https://github.com/nguyenphanno">
  <img
    src="https://img.shields.io/badge/GitHub-nguyenphanno-18181B?style=for-the-badge&logo=github&logoColor=white"
    alt="GitHub"
  />
</a>

<a href="https://discord.com/users/1469216989158309928">
  <img
    src="https://img.shields.io/badge/Discord-Presence-5865F2?style=for-the-badge&logo=discord&logoColor=white"
    alt="Discord"
  />
</a>

</div>

---

<div align="center">

## Support

If this project is useful to you, consider starring the repository.

<br>

<a href="https://github.com/nguyenphanno/Discord-Username-Generator-Checker">
  <img
    src="https://img.shields.io/badge/Star%20Repository-e3b341?style=for-the-badge&logo=github&logoColor=white"
    alt="Star Repository"
  />
</a>

<br><br>

<img
src="https://img.shields.io/github/last-commit/nguyenphanno/Discord-Username-Generator-Checker?style=flat-square&color=8B5CF6"
alt="Last Commit"
/>

<br><br>

Made by <strong>nguyenphanno</strong>

</div>
