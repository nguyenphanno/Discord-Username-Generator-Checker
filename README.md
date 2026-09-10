<div align="center">

<img src="https://capsule-render.vercel.app/api?type=waving&color=0:09090b,45:18181b,75:312e81,100:5865F2&height=220&section=header&text=Discord%20Username%20Tool&fontSize=44&fontColor=ffffff&fontAlignY=38&desc=Generator%20%26%20Checker&descAlignY=60&descSize=20&animation=fadeIn" width="100%" />

# Discord Username Generator & Checker

### High-performance · Multi-threaded · Proxy-powered

<p>
  <img src="https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Windows%20%7C%20Linux%20%7C%20macOS-Platform-18181B?style=for-the-badge" />
  <img src="https://img.shields.io/badge/HTTP%20%7C%20SOCKS5-Proxy-F97316?style=for-the-badge" />
</p>

<p>
  <img src="https://img.shields.io/github/stars/nguyenphanno/Discord-Username-Generator-Checker?style=for-the-badge&logo=github&logoColor=white&color=e3b341" />
  <img src="https://img.shields.io/github/forks/nguyenphanno/Discord-Username-Generator-Checker?style=for-the-badge&logo=github&logoColor=white&color=5865F2" />
  <img src="https://img.shields.io/github/issues/nguyenphanno/Discord-Username-Generator-Checker?style=for-the-badge&color=EF4444" />
  <img src="https://img.shields.io/github/license/nguyenphanno/Discord-Username-Generator-Checker?style=for-the-badge&color=22C55E" />
</p>

<br>

<img src="https://count.getloli.com/@nguyenphanno?name=nguyenphanno&theme=booru-qualityhentais&padding=1&offset=0&align=center&scale=1&pixelated=0&darkmode=auto" alt="Profile Views" />

<br><br>

<a href="https://github.com/nguyenphanno/Discord-Username-Generator-Checker">
<img src="https://img.shields.io/badge/View%20Repository-18181B?style=for-the-badge&logo=github&logoColor=white" />
</a>

</div>

---

<div align="center">

**Generate usernames · Check availability · Manage proxies · Export results**

</div>

## Preview

<div align="center">

<img
src="https://raw.githubusercontent.com/nguyenphanno/Discord-Username-Generator-Checker/main/image/nguyenphanno.png"
width="900"
alt="Main Preview"
/>

<br><br>

<img
src="https://raw.githubusercontent.com/nguyenphanno/Discord-Username-Generator-Checker/main/image/logs.png"
width="900"
alt="Logs Preview"
/>

</div>

---

## Features

<div align="center">

`⚡ Multi-threaded`  
`🤫 Smart Quiet`  
`🌐 Elite Proxy`  
`🔄 Auto Resume`

<br>

`📝 Sniper Logging`  
`🕵 Browser Fingerprinting`  
`🖥 Interactive Setup`

<br>

`💾 Session Export`  
`🛡 Proxy Health Check`

</div>

### High-performance checking

Multi-threaded username checking designed for fast concurrent processing.

### Smart Quiet Mode

Only displays the important results:

```text
AVAILABLE
TAKEN
```

### Elite Proxy System

Automatic proxy scoring and prioritization with health-checking and automatic disabling of low-quality proxies.

### Automatic Resume

Previously checked, taken and available usernames are automatically skipped.

### Sniper-ready Logging

Logs include timestamp and proxy information.

### Browser Fingerprinting

Uses realistic browser-related properties:

```text
User-Agent
Super-Properties
Timezone
Locale
```

### Session Export

Results are organized into session-based output and logs.

---

## Modes

### `live`

Continuously generate and check usernames.

### `check`

Check usernames from:

```text
usernames.txt
```

### `generate`

Generate usernames only.

### `both`

Generate usernames and then check them.

<br>

<div align="center">

<img src="https://img.shields.io/badge/LIVE-Continuous-8B5CF6?style=for-the-badge" />
<img src="https://img.shields.io/badge/CHECK-File%20Based-5865F2?style=for-the-badge" />
<img src="https://img.shields.io/badge/GENERATE-Generate-00ADD8?style=for-the-badge" />
<img src="https://img.shields.io/badge/BOTH-Combined-e3b341?style=for-the-badge" />

</div>

---

## Installation

### Clone

```bash
git clone https://github.com/nguyenphanno/Discord-Username-Generator-Checker.git
cd Discord-Username-Generator-Checker
```

### Dependencies

```bash
go mod tidy
```

### Proxies

Create:

```text
proxies.txt
```

Supported:

```text
http://ip:port
socks5://ip:port
```

### Run

```bash
go run main.go
```

---

## Usage

The tool starts with an interactive console setup.

Configure:

```text
Mode
Number of threads
Delay between requests
Target amount of usernames
Username length range
Smart Quiet mode
Proxy health-check settings
```

The configuration is selected directly from the console before the session starts.

---

## Proxy System

<div align="center">

<img src="https://img.shields.io/badge/HTTP-SUPPORTED-F97316?style=for-the-badge" />
<img src="https://img.shields.io/badge/SOCKS5-SUPPORTED-8B5CF6?style=for-the-badge" />

<br><br>

<b>Advanced Elite Proxy System</b>

<br>

<sub>
Automatic scoring · Prioritization · Health-checking · Auto-disabling
</sub>

</div>

### Supported formats

```text
http://ip:port
socks5://ip:port
```

### Proxy management

```text
Load
 ↓
Health Check
 ↓
Score
 ↓
Prioritize
 ↓
Use
 ↓
Monitor
```

Low-quality proxies can be automatically disabled, while high-quality working proxies can be exported to:

```text
working_proxies.txt
```

---

## Output

### Results

```text
available.txt
taken.txt
checked.txt
working_proxies.txt
```

### Sessions

```text
results/
```

Contains session logs and sniper logs.

<div align="center">

<img src="https://img.shields.io/badge/AVAILABLE-available.txt-22C55E?style=for-the-badge" />
<img src="https://img.shields.io/badge/TAKEN-taken.txt-EF4444?style=for-the-badge" />
<img src="https://img.shields.io/badge/CHECKED-checked.txt-5865F2?style=for-the-badge" />
<img src="https://img.shields.io/badge/WORKING%20PROXIES-working__proxies.txt-F97316?style=for-the-badge" />

</div>

---

## Recommended Settings

<div align="center">

### Balanced

```text
Threads      25 – 40
Delay        0.05 – 0.15 seconds
Smart Quiet  Enabled
Proxy Check  Enabled
```

</div>

---

## Notes

> **Username Availability**
> This tool only checks username availability.

> **Public Endpoint**
> Uses Discord’s public unauthenticated endpoint.

> **Proxy Quality**
> High-quality proxies are strongly recommended for best performance.

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
<img src="https://img.shields.io/badge/GitHub-nguyenphanno-18181B?style=for-the-badge&logo=github&logoColor=white" />
</a>

<a href="https://discord.com/users/1469216989158309928">
<img src="https://img.shields.io/badge/Discord-Presence-5865F2?style=for-the-badge&logo=discord&logoColor=white" />
</a>

</div>

---

<div align="center">

<img src="https://capsule-render.vercel.app/api?type=rect&color=5865F2&height=3&section=footer" width="100%" />

<br><br>

## Support the Project

If you find this project useful, consider giving it a star.

<br>

<a href="https://github.com/nguyenphanno/Discord-Username-Generator-Checker">
<img src="https://img.shields.io/badge/☆%20Star%20Repository-e3b341?style=for-the-badge&logo=github&logoColor=white" />
</a>

<br><br>

<img src="https://img.shields.io/github/last-commit/nguyenphanno/Discord-Username-Generator-Checker?style=flat-square&color=8B5CF6" />

<br><br>

Made by <b>nguyenphanno</b>

<br><br>

<img src="https://capsule-render.vercel.app/api?type=waving&color=0:5865F2,40:312e81,70:18181b,100:09090b&height=120&section=footer" width="100%" />

</div>
