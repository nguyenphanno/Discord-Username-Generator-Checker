package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/net/proxy"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	red     = "\033[38;2;255;85;85m"
	green   = "\033[38;2;80;250;123m"
	yellow  = "\033[38;2;255;184;108m"
	blue    = "\033[38;2;139;233;253m"
	magenta = "\033[38;2;255;121;198m"
	cyan    = "\033[38;2;139;233;253m"
	white   = "\033[38;2;248;248;242m"
	gray    = "\033[38;2;98;114;164m"
	purple  = "\033[38;2;189;147;249m"
	orange  = "\033[38;2;255;170;0m"
	version = "Discord<->Username<->Generator<->Checker"
)

var chromeProfiles = []struct {
	ua      string
	version string
	secChUa string
	platform string
}{
	{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36", "128.0.0.0", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`, `"Windows"`},
	{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36", "129.0.0.0", `"Chromium";v="129", "Not;A=Brand";v="24", "Google Chrome";v="129"`, `"Windows"`},
	{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36", "130.0.0.0", `"Chromium";v="130", "Not;A=Brand";v="24", "Google Chrome";v="130"`, `"Windows"`},
	{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36", "131.0.0.0", `"Chromium";v="131", "Not;A=Brand";v="24", "Google Chrome";v="131"`, `"Windows"`},
	{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36", "128.0.0.0", `"Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome";v="128"`, `"macOS"`},
	{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36", "129.0.0.0", `"Chromium";v="129", "Not;A=Brand";v="24", "Google Chrome";v="129"`, `"macOS"`},
	{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36", "130.0.0.0", `"Chromium";v="130", "Not;A=Brand";v="24", "Google Chrome";v="130"`, `"macOS"`},
	{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36", "129.0.0.0", `"Chromium";v="129", "Not;A=Brand";v="24", "Google Chrome";v="129"`, `"Linux"`},
	{"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:130.0) Gecko/20100101 Firefox/130.0", "130.0", "", `"Windows"`},
	{"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:131.0) Gecko/20100101 Firefox/131.0", "131.0", "", `"Windows"`},
	{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:130.0) Gecko/20100101 Firefox/130.0", "130.0", "", `"macOS"`},
}

var timezones = []string{
	"America/New_York", "America/Los_Angeles", "America/Chicago", "America/Denver",
	"Europe/London", "Europe/Paris", "Europe/Berlin", "Europe/Moscow", "Europe/Amsterdam",
	"Asia/Bangkok", "Asia/Singapore", "Asia/Tokyo", "Asia/Ho_Chi_Minh", "Asia/Jakarta",
	"Asia/Shanghai", "Asia/Seoul", "Australia/Sydney", "Australia/Melbourne", "UTC",
	"Pacific/Auckland", "America/Sao_Paulo", "Asia/Dubai", "Asia/Kolkata",
}

var locales = []string{"en-US", "en-GB", "en-AU", "fr-FR", "de-DE", "es-ES", "pt-BR", "ja-JP", "ko-KR", "vi-VN", "th-TH", "id-ID"}

type RuntimeConfig struct {
	Mode              string
	Threads           int
	Delay             float64
	UsernamesToGen    int
	UsernameMinLength int
	UsernameMaxLength int
	ProxyCheck        bool
	ProxyCheckAmount  int
	MaxFailCount      int
	MinProxyScore     int
	SmartQuiet        bool
}

var logMu sync.Mutex
var smartQuiet bool

func ts() string { return time.Now().Format("15:04:05") }

func logInfo(format string, args ...any) {
	logMu.Lock()
	defer logMu.Unlock()
	fmt.Printf("%s%s%s  %sINFO %s %s\n", gray, ts(), reset, cyan+bold, reset, fmt.Sprintf(format, args...))
}

func logSuccess(format string, args ...any) {
	logMu.Lock()
	defer logMu.Unlock()
	fmt.Printf("%s%s%s  %sHIT  %s %s\n", gray, ts(), reset, green+bold, reset, fmt.Sprintf(format, args...))
}

func logFail(format string, args ...any) {
	logMu.Lock()
	defer logMu.Unlock()
	fmt.Printf("%s%s%s  %sTAKE %s %s\n", gray, ts(), reset, red+bold, reset, fmt.Sprintf(format, args...))
}

func logError(format string, args ...any) {
	if smartQuiet {
		return
	}
	logMu.Lock()
	defer logMu.Unlock()
	fmt.Printf("%s%s%s  %sERR  %s %s\n", gray, ts(), reset, yellow+bold, reset, fmt.Sprintf(format, args...))
}

func logWarn(format string, args ...any) {
	logMu.Lock()
	defer logMu.Unlock()
	fmt.Printf("%s%s%s  %sWARN %s %s\n", gray, ts(), reset, magenta+bold, reset, fmt.Sprintf(format, args...))
}

func logSection(title string) {
	logMu.Lock()
	defer logMu.Unlock()
	fmt.Printf("\n%s%s%s\n", purple, strings.Repeat("─", 78), reset)
	fmt.Printf("%s  %s%s%s\n", purple, bold+white, title, reset)
	fmt.Printf("%s%s%s\n\n", purple, strings.Repeat("─", 78), reset)
}

func printBanner(proxyCount int, mode string) {
	fmt.Printf(`%s
                                                                                                    ▄                            ▄                          
                                                                                                 ▄██                          ▄██  ▀█▄                     
       ▄                                                                                         ███                          ███  ▐██▌                    
 ▄█▌  ▐█▓▄  ▄███████▀ ▄██████▄ ■▄█ ▄█▄▄   ▄█ ▄█▄▄    ▄███▄▄▄  ▄█▄▄ ▄█▄▄   ▄██████▄      ▄██████▄ █▓█ ▄█▄▄  ▄██████▄  ▄██████▄ █▓█ ▄██▀  ▄██████▄ ■▄█ ▄█▄▄  
▐▒▓    ▓▒▓▌▐▓▓▄      ▓▒▓  ▀▓▒▓▌▐▒▓▀  ▓▒▓▌▓▒▓▀▀▓▒▓▌ ▄▒▓  ▀▓▒▓▌▓▒▓▀▒▓▀▓▒▓▌ ▓▒▓  ▀▓▒▓▌    ▓▒▓  ▀▓▒▓▌▓▒▓▀▀▓▒▓▌▓▒▓  ▀▓▒▓▌▓▒▓  ▀▓▒▓▌▓▒▓▀▀▓▒▓▌▓▒▓  ▀▓▒▓▌▐▒▓▀  ▓▒▓▌
▒░▒▌   ▐░▒▌ ▀▒▒▒▒▒▒▒▄▒░▒▄▄▄▒▀▀ ▐░▒   ▐▒░▒▒░▒  ▐▒░▒ ▒░▒   ▐▒░▒▒░▒ ░░ ▐▒░▒ ▒░▒▄▄▄▒▀▀     ▒░▒   ▀▀  ▒░▒  ▐▒░▒▒░▒▄▄▄▒▀▀ ▒░▒   ▀▀  ▒░▒  ▐▒░▒▒░▒▄▄▄▒▀▀ ▐░▒   ▐▒░▒
▀█░░▄  ░▀░  ▄  ▀▀▀░░▌▐ ░▄    ▄▌▐ ░       ░ ░   ░ ░▌░ ░ ▄▀ ░ ░░ ░  █  ░ ░▌▐ ░▄    ▄▌    ▐ ░▄    ▄▌░ ░   ░ ░▐ ░▄    ▄▌▐ ░▄    ▄▌░ ░   ░ ░▐ ░▄    ▄▌▐ ░       
  ▀███░▒▀  ▀▀      ▀  ▀     ▀▀ ▐                   ▀  ▀    ▀              ▀     ▀▀      ▀     ▀▀           ▀     ▀▀  ▀     ▀▀           ▀     ▀▀ ▐         
                                              ▐ ▀                   ▐ ▀                               ▐  ▀                         ▐  ▀                    
                                              ▀                     ▀                                  ▀                            ▀                      
                                                                                                     ▀                            ▀
%s`, purple, reset)
	fmt.Printf("%s  %s%s%s  •  Proxies: %s%d%s  •  Mode: %s%s%s\n\n", gray, bold+cyan, version, reset, orange, proxyCount, reset, purple, strings.ToUpper(mode), reset)
}

func ask(prompt string, def string) string {
	fmt.Printf("%s%s%s [%s]: ", cyan, prompt, reset, def)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		return def
	}
	return text
}

func askInt(prompt string, def int) int {
	for {
		s := ask(prompt, strconv.Itoa(def))
		v, err := strconv.Atoi(s)
		if err == nil && v >= 0 {
			return v
		}
		fmt.Println(red + "Invalid number" + reset)
	}
}

func askFloat(prompt string, def float64) float64 {
	for {
		s := ask(prompt, fmt.Sprintf("%.2f", def))
		v, err := strconv.ParseFloat(s, 64)
		if err == nil && v >= 0 {
			return v
		}
		fmt.Println(red + "Invalid number" + reset)
	}
}

func askBool(prompt string, def bool) bool {
	d := "n"
	if def {
		d = "y"
	}
	s := strings.ToLower(ask(prompt+" (y/n)", d))
	return s == "y" || s == "yes" || s == "true" || s == "1"
}

func interactiveConfig() RuntimeConfig {
	fmt.Println(purple + "═══════════════ INTERACTIVE SETUP ═══════════════" + reset)
	mode := strings.ToLower(ask("Mode (live/check/generate/both)", "live"))
	threads := askInt("Threads", 35)
	delay := askFloat("Delay (seconds)", 0.08)
	target := askInt("Target usernames (for live/generate)", 5000)
	minLen := askInt("Username min length", 3)
	maxLen := askInt("Username max length", 8)
	if minLen > maxLen {
		minLen, maxLen = maxLen, minLen
	}
	proxyCheck := askBool("Health-check proxies first?", true)
	checkAmount := 150
	if proxyCheck {
		checkAmount = askInt("How many proxies to health-check", 200)
	}
	smartQ := askBool("Smart Quiet (only show HIT + TAKE)?", true)
	fmt.Println()
	return RuntimeConfig{
		Mode:              mode,
		Threads:           threads,
		Delay:             delay,
		UsernamesToGen:    target,
		UsernameMinLength: minLen,
		UsernameMaxLength: maxLen,
		ProxyCheck:        proxyCheck,
		ProxyCheckAmount:  checkAmount,
		MaxFailCount:      3,
		MinProxyScore:     2,
		SmartQuiet:        smartQ,
	}
}

type ProxyPool struct {
	mu         sync.RWMutex
	proxies    []string
	failCount  map[string]int
	scores     map[string]int
	successes  map[string]int
	rateUntil  map[string]time.Time
	disabled   map[string]bool
	working    map[string]bool
	schemes    map[string]string
	maxFail    int
	minScore   int
	totalReq   atomic.Int64
	successReq atomic.Int64
	elite      []string
	eliteMu    sync.RWMutex
}

func newProxyPool(proxies []string, maxFail, minScore int) *ProxyPool {
	p := &ProxyPool{
		proxies:   proxies,
		failCount: make(map[string]int),
		scores:    make(map[string]int),
		successes: make(map[string]int),
		rateUntil: make(map[string]time.Time),
		disabled:  make(map[string]bool),
		working:   make(map[string]bool),
		schemes:   make(map[string]string),
		maxFail:   maxFail,
		minScore:  minScore,
		elite:     make([]string, 0),
	}
	for _, px := range readLines("working_proxies.txt") {
		p.working[px] = true
		p.scores[px] = 9
		p.successes[px] = 4
		p.elite = append(p.elite, px)
	}
	rand.Shuffle(len(p.proxies), func(i, j int) {
		p.proxies[i], p.proxies[j] = p.proxies[j], p.proxies[i]
	})
	return p
}

func (p *ProxyPool) schemeOf(px string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.schemes[px]
}

func (p *ProxyPool) setScheme(px, scheme string) {
	if px == "" {
		return
	}
	p.mu.Lock()
	p.schemes[px] = scheme
	p.mu.Unlock()
}

func (p *ProxyPool) pick() (string, bool) {
	now := time.Now()
	p.eliteMu.RLock()
	eliteCopy := make([]string, len(p.elite))
	copy(eliteCopy, p.elite)
	p.eliteMu.RUnlock()

	p.mu.RLock()
	var active, good, eliteActive []string
	for _, px := range eliteCopy {
		if p.disabled[px] {
			continue
		}
		if t, ok := p.rateUntil[px]; ok && t.After(now) {
			continue
		}
		eliteActive = append(eliteActive, px)
	}
	for _, px := range p.proxies {
		if p.disabled[px] {
			continue
		}
		if t, ok := p.rateUntil[px]; ok && t.After(now) {
			continue
		}
		active = append(active, px)
		if p.scores[px] >= p.minScore {
			good = append(good, px)
		}
	}
	p.mu.RUnlock()

	if len(eliteActive) > 0 && rand.Float32() < 0.78 {
		return eliteActive[rand.Intn(len(eliteActive))], true
	}
	if len(good) > 0 {
		return good[rand.Intn(len(good))], true
	}
	if len(active) == 0 {
		return "", false
	}
	return active[rand.Intn(len(active))], true
}

func (p *ProxyPool) reportFailure(px string) {
	if px == "" {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.failCount[px]++
	p.scores[px] = max(p.scores[px]-3, -12)
	if p.failCount[px] >= p.maxFail {
		p.disabled[px] = true
		p.eliteMu.Lock()
		for i, e := range p.elite {
			if e == px {
				p.elite = append(p.elite[:i], p.elite[i+1:]...)
				break
			}
		}
		p.eliteMu.Unlock()
	}
}

func (p *ProxyPool) reportSuccess(px string) {
	if px == "" {
		return
	}
	p.mu.Lock()
	p.scores[px] = min(p.scores[px]+2, 16)
	p.successes[px]++
	delete(p.failCount, px)
	already := p.working[px]
	if !already {
		p.working[px] = true
	}
	p.mu.Unlock()
	p.successReq.Add(1)
	if !already {
		appendLine("working_proxies.txt", px)
	}
	if p.successes[px] >= 2 {
		p.eliteMu.Lock()
		found := false
		for _, e := range p.elite {
			if e == px {
				found = true
				break
			}
		}
		if !found {
			p.elite = append(p.elite, px)
		}
		p.eliteMu.Unlock()
	}
}

func (p *ProxyPool) rateLimit(px string, d time.Duration) {
	if px == "" {
		return
	}
	if d > 25*time.Minute {
		d = 25 * time.Minute
	}
	p.mu.Lock()
	p.rateUntil[px] = time.Now().Add(d)
	p.scores[px] = max(p.scores[px]-1, -6)
	p.mu.Unlock()
}

func (p *ProxyPool) disable(px string) {
	if px == "" {
		return
	}
	p.mu.Lock()
	p.disabled[px] = true
	p.mu.Unlock()
}

func (p *ProxyPool) activeCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	n := 0
	for _, px := range p.proxies {
		if !p.disabled[px] {
			n++
		}
	}
	return n
}

func (p *ProxyPool) workingCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.working)
}

func (p *ProxyPool) eliteCount() int {
	p.eliteMu.RLock()
	defer p.eliteMu.RUnlock()
	return len(p.elite)
}

func appendLine(path, line string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(line + "\n")
}

func readLines(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var out []string
	for _, ln := range strings.Split(string(data), "\n") {
		ln = strings.TrimSpace(ln)
		if ln != "" {
			out = append(out, ln)
		}
	}
	return out
}

func loadSkipSet() map[string]struct{} {
	skip := make(map[string]struct{})
	for _, f := range []string{"checked.txt", "taken.txt", "available.txt"} {
		for _, u := range readLines(f) {
			skip[strings.ToLower(strings.TrimSpace(u))] = struct{}{}
		}
	}
	resultsDir := "results"
	entries, err := os.ReadDir(resultsDir)
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if strings.HasPrefix(name, "available_") || strings.HasPrefix(name, "taken_") || strings.HasPrefix(name, "checked_") || strings.HasPrefix(name, "sniper_") {
				for _, u := range readLines(filepath.Join(resultsDir, name)) {
					parts := strings.SplitN(u, " | ", 2)
					skip[strings.ToLower(strings.TrimSpace(parts[0]))] = struct{}{}
				}
			}
		}
	}
	return skip
}

type lineWriter struct {
	ch chan string
}

func newLineWriter(path string) *lineWriter {
	w := &lineWriter{ch: make(chan string, 4096)}
	go func() {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			for range w.ch {
			}
			return
		}
		defer f.Close()
		buf := make([]string, 0, 256)
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		flush := func() {
			if len(buf) == 0 {
				return
			}
			for _, l := range buf {
				f.WriteString(l + "\n")
			}
			f.Sync()
			buf = buf[:0]
		}
		for {
			select {
			case l, ok := <-w.ch:
				if !ok {
					flush()
					return
				}
				buf = append(buf, l)
				if len(buf) >= 256 {
					flush()
				}
			case <-ticker.C:
				flush()
			}
		}
	}()
	return w
}

func (w *lineWriter) write(s string) { w.ch <- s }
func (w *lineWriter) close()         { close(w.ch) }

func validUsername(u string) bool {
	n := len(u)
	if n < 2 || n > 32 {
		return false
	}
	edge := func(c byte) bool {
		return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_'
	}
	valid := func(c byte) bool { return edge(c) || c == '.' }
	if !edge(u[0]) || !edge(u[n-1]) {
		return false
	}
	for i := 0; i < n; i++ {
		if !valid(u[i]) {
			return false
		}
	}
	return !strings.Contains(u, "..")
}

const firstChars = "abcdefghijklmnopqrstuvwxyz0123456789_"
const restChars = firstChars + "."

func genUsername(length int) string {
	b := make([]byte, length)
	for {
		b[0] = firstChars[rand.Intn(len(firstChars))]
		for i := 1; i < length; i++ {
			b[i] = restChars[rand.Intn(len(restChars))]
		}
		u := string(b)
		if validUsername(u) {
			return u
		}
	}
}

func randomUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func buildSuperProperties(profile struct {
	ua       string
	version  string
	secChUa  string
	platform string
}) string {
	osName := "Windows"
	if strings.Contains(profile.platform, "mac") {
		osName = "Mac OS X"
	} else if strings.Contains(profile.platform, "Linux") {
		osName = "Linux"
	}
	props := map[string]any{
		"os":                       osName,
		"browser":                  "Chrome",
		"device":                   "",
		"system_locale":            locales[rand.Intn(len(locales))],
		"has_client_mods":          false,
		"browser_user_agent":       profile.ua,
		"browser_version":          profile.version,
		"os_version":               "10",
		"referrer":                 "",
		"referring_domain":         "",
		"referrer_current":         "",
		"referring_domain_current": "",
		"release_channel":          "stable",
		"client_build_number":      398000 + rand.Intn(28000),
		"client_event_source":      nil,
		"client_launch_id":         randomUUID(),
		"client_app_state":         "focused",
		"client_heartbeat_session_id": randomUUID(),
	}
	b, _ := json.Marshal(props)
	return base64.StdEncoding.EncodeToString(b)
}

func shortProxy(px string) string {
	if px == "" {
		return "—"
	}
	u, err := url.Parse(px)
	if err != nil {
		return px
	}
	host := u.Host
	if host == "" {
		host = px
	}
	if len(host) > 22 {
		return host[:19] + "..."
	}
	return host
}

func detectScheme(px string) string {
	if px == "" {
		return "—"
	}
	u, err := url.Parse(px)
	if err == nil && u.Scheme != "" {
		return strings.ToLower(u.Scheme)
	}
	return "http"
}

const apiURL = "https://discord.com/api/v9/unique-username/username-attempt-unauthed"
const pingURL = "https://discord.com/api/v9/ping"

type Checker struct {
	cfg        RuntimeConfig
	pool       *ProxyPool
	transports sync.Map
	stats      struct {
		available atomic.Int64
		taken     atomic.Int64
		failed    atomic.Int64
		checked   atomic.Int64
		rateLimit atomic.Int64
	}
	startTime  time.Time
	resultsDir string
	sessionTag string
}

func newChecker(cfg RuntimeConfig, pool *ProxyPool) *Checker {
	tag := time.Now().Format("20060102_150405")
	dir := "results"
	os.MkdirAll(dir, 0755)
	return &Checker{
		cfg:        cfg,
		pool:       pool,
		startTime:  time.Now(),
		resultsDir: dir,
		sessionTag: tag,
	}
}

type apiResp struct {
	status     int
	taken      bool
	validJSON  bool
	retryAfter float64
	body       string
	latency    time.Duration
}

func (c *Checker) doRequest(rt http.RoundTripper, payload []byte) (*apiResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	profile := chromeProfiles[rand.Intn(len(chromeProfiles))]
	locale := locales[rand.Intn(len(locales))]
	tz := timezones[rand.Intn(len(timezones))]

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", profile.ua)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", locale+",en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("Origin", "https://discord.com")
	req.Header.Set("Referer", "https://discord.com/register")
	req.Header.Set("X-Super-Properties", buildSuperProperties(profile))
	req.Header.Set("X-Discord-Locale", locale)
	req.Header.Set("X-Discord-Timezone", tz)
	req.Header.Set("X-Debug-Options", "bugReporterEnabled")
	if profile.secChUa != "" {
		req.Header.Set("Sec-Ch-Ua", profile.secChUa)
		req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
		req.Header.Set("Sec-Ch-Ua-Platform", profile.platform)
	}
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	start := time.Now()
	client := &http.Client{Transport: rt, Timeout: 9 * time.Second}
	resp, err := client.Do(req)
	latency := time.Since(start)
	if err != nil {
		return &apiResp{latency: latency}, err
	}
	defer resp.Body.Close()
	out := &apiResp{status: resp.StatusCode, latency: latency}
	buf, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	out.body = string(buf)

	if out.status == 200 {
		var j struct {
			Taken bool `json:"taken"`
		}
		if err := json.Unmarshal(buf, &j); err == nil {
			out.taken = j.Taken
			out.validJSON = true
		}
	}
	if out.status == 429 {
		var j struct {
			RetryAfter float64 `json:"retry_after"`
		}
		if err := json.Unmarshal(buf, &j); err == nil && j.RetryAfter > 0 {
			out.retryAfter = j.RetryAfter
		} else {
			out.retryAfter = 5
		}
	}
	return out, nil
}

func (c *Checker) transportFor(raw, scheme string) (http.RoundTripper, error) {
	key := scheme + "|" + raw
	if v, ok := c.transports.Load(key); ok {
		return v.(*http.Transport), nil
	}
	t := &http.Transport{
		MaxIdleConns:        c.cfg.Threads * 6,
		MaxIdleConnsPerHost: c.cfg.Threads * 4,
		MaxConnsPerHost:     c.cfg.Threads * 4,
		IdleConnTimeout:     60 * time.Second,
		TLSHandshakeTimeout: 6 * time.Second,
		DialContext:         (&net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		ForceAttemptHTTP2:   true,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	if scheme == "socks5" || scheme == "socks5h" {
		host := raw
		if i := strings.Index(host, "://"); i >= 0 {
			host = host[i+3:]
		}
		d, err := proxy.SOCKS5("tcp", host, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}
		t.Proxy = nil
		t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			type dialRes struct {
				conn net.Conn
				err  error
			}
			ch := make(chan dialRes, 1)
			go func() {
				cn, er := d.Dial(network, addr)
				ch <- dialRes{cn, er}
			}()
			select {
			case r := <-ch:
				return r.conn, r.err
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	} else {
		norm := raw
		if !strings.Contains(norm, "://") {
			norm = "http://" + norm
		}
		u, err := url.Parse(norm)
		if err != nil {
			return nil, err
		}
		t.Proxy = http.ProxyURL(u)
	}
	actual, _ := c.transports.LoadOrStore(key, t)
	return actual.(*http.Transport), nil
}

func (c *Checker) requestViaProxy(raw string, payload []byte) (*apiResp, string, error) {
	scheme := c.pool.schemeOf(raw)
	if scheme == "" {
		scheme = detectScheme(raw)
		if scheme == "" {
			scheme = "http"
		}
	}
	rt, err := c.transportFor(raw, scheme)
	if err != nil {
		return nil, scheme, err
	}
	res, err := c.doRequest(rt, payload)
	if err == nil {
		if c.pool.schemeOf(raw) == "" {
			c.pool.setScheme(raw, scheme)
		}
		return res, scheme, nil
	}
	if c.pool.schemeOf(raw) == "" {
		rt2, err2 := c.transportFor(raw, "socks5")
		if err2 == nil {
			if res2, err3 := c.doRequest(rt2, payload); err3 == nil {
				c.pool.setScheme(raw, "socks5")
				return res2, "socks5", nil
			}
		}
		c.pool.setScheme(raw, "http")
	}
	return res, scheme, err
}

type checkResult struct {
	status   int
	errMsg   string
	proxy    string
	scheme   string
	latency  time.Duration
	httpCode int
	respBody string
}

func (c *Checker) check(username string) checkResult {
	if !validUsername(username) {
		return checkResult{status: 0, errMsg: "invalid"}
	}
	payload := []byte(fmt.Sprintf(`{"username":%q}`, username))
	var lastErr string
	var lastProxy string
	var lastScheme string
	var lastLatency time.Duration
	var lastCode int
	var lastBody string

	for attempt := 0; attempt < 4; attempt++ {
		px, ok := c.pool.pick()
		if !ok {
			return checkResult{status: 0, errMsg: "no proxies", proxy: lastProxy, scheme: lastScheme, latency: lastLatency}
		}
		lastProxy = px
		c.pool.totalReq.Add(1)
		res, scheme, err := c.requestViaProxy(px, payload)
		lastScheme = scheme
		if res != nil {
			lastLatency = res.latency
			lastCode = res.status
			lastBody = res.body
		}
		if err != nil {
			lastErr = "timeout"
			c.pool.reportFailure(px)
			time.Sleep(time.Duration(80+attempt*100) * time.Millisecond)
			continue
		}
		switch res.status {
		case 200:
			if res.validJSON {
				c.pool.reportSuccess(px)
				if res.taken {
					return checkResult{status: 2, proxy: px, scheme: scheme, latency: res.latency, httpCode: 200, respBody: res.body}
				}
				return checkResult{status: 1, proxy: px, scheme: scheme, latency: res.latency, httpCode: 200, respBody: res.body}
			}
			lastErr = "bad body"
			c.pool.reportFailure(px)
			time.Sleep(60 * time.Millisecond)
			continue
		case 429:
			c.stats.rateLimit.Add(1)
			lastErr = "rate-limit"
			c.pool.reportFailure(px)
			c.pool.rateLimit(px, time.Duration(res.retryAfter*float64(time.Second)))
			time.Sleep(180 * time.Millisecond)
			continue
		case 400:
			c.pool.reportFailure(px)
			return checkResult{status: 0, errMsg: "bad req", proxy: px, scheme: scheme, latency: res.latency, httpCode: 400, respBody: res.body}
		default:
			lastErr = fmt.Sprintf("http%d", res.status)
			c.pool.reportFailure(px)
			time.Sleep(time.Duration(60+attempt*80) * time.Millisecond)
		}
	}
	return checkResult{status: 0, errMsg: lastErr, proxy: lastProxy, scheme: lastScheme, latency: lastLatency, httpCode: lastCode, respBody: lastBody}
}

func (c *Checker) healthCheck(amount int) int {
	pxs := c.pool.proxies
	if len(pxs) > amount {
		pxs = pxs[:amount]
	}
	logInfo("Health-checking %s%d%s proxies...", cyan, len(pxs), reset)
	var wg sync.WaitGroup
	var okCount atomic.Int64
	sem := make(chan struct{}, 80)
	start := time.Now()
	for _, px := range pxs {
		wg.Add(1)
		go func(px string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			req, err := http.NewRequestWithContext(ctx, "GET", pingURL, nil)
			if err != nil {
				c.pool.disable(px)
				return
			}
			schemes := []string{c.pool.schemeOf(px)}
			if schemes[0] == "" {
				schemes = []string{"http", "socks5"}
			}
			ok := false
			for _, sc := range schemes {
				rt, err := c.transportFor(px, sc)
				if err != nil {
					continue
				}
				client := &http.Client{Transport: rt, Timeout: 5 * time.Second}
				if resp, err := client.Do(req); err == nil {
					resp.Body.Close()
					if c.pool.schemeOf(px) == "" {
						c.pool.setScheme(px, sc)
					}
					ok = true
					break
				}
			}
			if !ok {
				c.pool.disable(px)
			} else {
				okCount.Add(1)
			}
		}(px)
	}
	wg.Wait()
	alive := int(okCount.Load())
	elapsed := time.Since(start).Seconds()
	logInfo("Proxy health finished in %.1fs → %s%d%s alive / %d tested  •  Elite: %s%d%s  •  Active: %s%d%s",
		elapsed, green, alive, reset, len(pxs), orange, c.pool.eliteCount(), reset, cyan, c.pool.activeCount(), reset)
	if alive == 0 {
		logFail("No working proxies found")
	}
	return alive
}

func (c *Checker) runGenerate() int {
	minL, maxL := c.cfg.UsernameMinLength, c.cfg.UsernameMaxLength
	seen := make(map[string]struct{}, c.cfg.UsernamesToGen)
	jobs := make(chan int, c.cfg.Threads*8)
	results := make(chan string, c.cfg.UsernamesToGen)
	var wg sync.WaitGroup
	workers := min(c.cfg.Threads, 32)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := range jobs {
				results <- genUsername(n)
			}
		}()
	}
	go func() {
		for i := 0; i < c.cfg.UsernamesToGen; i++ {
			jobs <- rand.Intn(maxL-minL+1) + minL
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()
	count := 0
	for u := range results {
		if _, dup := seen[u]; !dup {
			seen[u] = struct{}{}
			count++
			if count <= 10 || count%100 == 0 {
				logSuccess("Generated  %s%s%s", bold+white, u, reset)
			}
		}
	}
	list := make([]string, 0, len(seen))
	for u := range seen {
		list = append(list, u)
	}
	sort.Strings(list)
	os.WriteFile("usernames.txt", []byte(strings.Join(list, "\n")), 0644)
	logInfo("Saved %s%d%s unique usernames → usernames.txt", green, len(list), reset)
	return len(list)
}

type result struct {
	username string
	res      checkResult
}

func formatLatency(d time.Duration) string {
	if d <= 0 {
		return "—"
	}
	return fmt.Sprintf("%dms", d.Milliseconds())
}

func shortBody(body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return ""
	}
	body = strings.ReplaceAll(body, "\n", " ")
	if len(body) > 55 {
		return body[:52] + "..."
	}
	return body
}

func (c *Checker) runCheck(usernames []string) {
	skip := loadSkipSet()
	seen := make(map[string]struct{}, len(usernames))
	toCheck := make([]string, 0, len(usernames))
	for _, u := range usernames {
		ul := strings.ToLower(strings.TrimSpace(u))
		if _, dup := seen[ul]; dup {
			continue
		}
		seen[ul] = struct{}{}
		if !validUsername(u) {
			continue
		}
		if _, done := skip[ul]; done {
			continue
		}
		toCheck = append(toCheck, u)
	}
	total := len(toCheck)
	logSection(fmt.Sprintf("CHECK MODE  •  %d usernames  •  %d threads  •  Elite %d  •  Active %d", total, c.cfg.Threads, c.pool.eliteCount(), c.pool.activeCount()))
	if total == 0 {
		logWarn("Nothing left to check")
		return
	}

	wChecked := newLineWriter("checked.txt")
	wAvail := newLineWriter("available.txt")
	wTaken := newLineWriter("taken.txt")
	wSniper := newLineWriter(filepath.Join(c.resultsDir, "sniper_"+c.sessionTag+".txt"))
	wAvailSess := newLineWriter(filepath.Join(c.resultsDir, "available_"+c.sessionTag+".txt"))
	wTakenSess := newLineWriter(filepath.Join(c.resultsDir, "taken_"+c.sessionTag+".txt"))
	defer wChecked.close()
	defer wAvail.close()
	defer wTaken.close()
	defer wSniper.close()
	defer wAvailSess.close()
	defer wTakenSess.close()

	jobs := make(chan string, c.cfg.Threads*6)
	results := make(chan result, c.cfg.Threads*6)
	var wg sync.WaitGroup
	delay := time.Duration(c.cfg.Delay * float64(time.Second))

	for i := 0; i < c.cfg.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for u := range jobs {
				res := c.check(u)
				if delay > 0 {
					time.Sleep(delay)
				}
				results <- result{u, res}
			}
		}()
	}
	go func() {
		for _, u := range toCheck {
			jobs <- u
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	start := time.Now()
	for r := range results {
		n := c.stats.checked.Add(1)
		wChecked.write(r.username)
		px := shortProxy(r.res.proxy)
		scheme := r.res.scheme
		if scheme == "" {
			scheme = "—"
		}
		lat := formatLatency(r.res.latency)
		bodyStr := shortBody(r.res.respBody)
		nowStr := time.Now().Format("2006-01-02 15:04:05")

		switch r.res.status {
		case 1:
			c.stats.available.Add(1)
			logSuccess("%-16s AVAILABLE  %-7s %-20s %-5s %s  [%d/%d]",
				r.username, scheme, px, lat, bodyStr, n, total)
			wAvail.write(r.username)
			wAvailSess.write(r.username)
			wSniper.write(fmt.Sprintf("%s | %s | proxy=%s | scheme=%s | latency=%s", r.username, nowStr, r.res.proxy, scheme, lat))
		case 2:
			c.stats.taken.Add(1)
			logFail("%-16s taken       %-7s %-20s %-5s %s  [%d/%d]",
				r.username, scheme, px, lat, bodyStr, n, total)
			wTaken.write(r.username)
			wTakenSess.write(r.username)
		default:
			c.stats.failed.Add(1)
			logError("%-16s %-10s %-7s %-20s %-5s %s  [%d/%d]",
				r.username, r.res.errMsg, scheme, px, lat, bodyStr, n, total)
		}
	}
	elapsed := time.Since(start).Seconds()
	logSection("SUMMARY")
	logInfo("Finished in %s%.1fs%s", cyan, elapsed, reset)
	logInfo("Available : %s%d%s", green, c.stats.available.Load(), reset)
	logInfo("Taken     : %s%d%s", red, c.stats.taken.Load(), reset)
	logInfo("Failed    : %s%d%s", yellow, c.stats.failed.Load(), reset)
	logInfo("RateLimit : %s%d%s", magenta, c.stats.rateLimit.Load(), reset)
	logInfo("Speed     : %s%.1f%s req/s", cyan, float64(c.stats.checked.Load())/elapsed, reset)
	logInfo("Elite     : %s%d%s  •  Working: %s%d%s", orange, c.pool.eliteCount(), reset, green, c.pool.workingCount(), reset)
	logInfo("Sniper    : %sresults/sniper_%s.txt%s", cyan, c.sessionTag, reset)
}

func (c *Checker) runLive(target int) {
	logSection(fmt.Sprintf("LIVE MODE  •  Target %d  •  %d threads  •  Elite %d  •  Active %d",
		target, c.cfg.Threads, c.pool.eliteCount(), c.pool.activeCount()))

	jobs := make(chan string, c.cfg.Threads*6)
	results := make(chan result, c.cfg.Threads*6)
	stop := make(chan struct{})
	var once sync.Once
	var wg sync.WaitGroup

	wChecked := newLineWriter("checked.txt")
	wAvail := newLineWriter("available.txt")
	wTaken := newLineWriter("taken.txt")
	wSniper := newLineWriter(filepath.Join(c.resultsDir, "sniper_"+c.sessionTag+".txt"))
	wAvailSess := newLineWriter(filepath.Join(c.resultsDir, "available_"+c.sessionTag+".txt"))
	wTakenSess := newLineWriter(filepath.Join(c.resultsDir, "taken_"+c.sessionTag+".txt"))
	defer wChecked.close()
	defer wAvail.close()
	defer wTaken.close()
	defer wSniper.close()
	defer wAvailSess.close()
	defer wTakenSess.close()

	var produced atomic.Int64
	seen := sync.Map{}
	skip := loadSkipSet()

	go func() {
		defer close(jobs)
		for produced.Load() < int64(target) {
			select {
			case <-stop:
				return
			default:
			}
			minL, maxL := c.cfg.UsernameMinLength, c.cfg.UsernameMaxLength
			u := genUsername(rand.Intn(maxL-minL+1) + minL)
			ul := strings.ToLower(u)
			if _, dup := seen.LoadOrStore(ul, struct{}{}); dup {
				continue
			}
			if _, sk := skip[ul]; sk {
				continue
			}
			produced.Add(1)
			jobs <- u
		}
	}()

	delay := time.Duration(c.cfg.Delay * float64(time.Second))
	for i := 0; i < c.cfg.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for u := range jobs {
				res := c.check(u)
				if delay > 0 {
					time.Sleep(delay)
				}
				results <- result{u, res}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	start := time.Now()
	go func() {
		<-sig
		once.Do(func() { close(stop) })
		logWarn("Interrupted by user")
	}()

	for r := range results {
		n := c.stats.checked.Add(1)
		wChecked.write(r.username)
		px := shortProxy(r.res.proxy)
		scheme := r.res.scheme
		if scheme == "" {
			scheme = "—"
		}
		lat := formatLatency(r.res.latency)
		bodyStr := shortBody(r.res.respBody)
		nowStr := time.Now().Format("2006-01-02 15:04:05")

		switch r.res.status {
		case 1:
			c.stats.available.Add(1)
			logSuccess("%-16s AVAILABLE  %-7s %-20s %-5s %s  [%d/%d]",
				r.username, scheme, px, lat, bodyStr, n, target)
			wAvail.write(r.username)
			wAvailSess.write(r.username)
			wSniper.write(fmt.Sprintf("%s | %s | proxy=%s | scheme=%s | latency=%s", r.username, nowStr, r.res.proxy, scheme, lat))
		case 2:
			c.stats.taken.Add(1)
			logFail("%-16s taken       %-7s %-20s %-5s %s  [%d/%d]",
				r.username, scheme, px, lat, bodyStr, n, target)
			wTaken.write(r.username)
			wTakenSess.write(r.username)
		default:
			c.stats.failed.Add(1)
			logError("%-16s %-10s %-7s %-20s %-5s %s  [%d/%d]",
				r.username, r.res.errMsg, scheme, px, lat, bodyStr, n, target)
		}
	}
	elapsed := time.Since(start).Seconds()
	logSection("SUMMARY")
	logInfo("Finished in %s%.1fs%s", cyan, elapsed, reset)
	logInfo("Available : %s%d%s", green, c.stats.available.Load(), reset)
	logInfo("Taken     : %s%d%s", red, c.stats.taken.Load(), reset)
	logInfo("Failed    : %s%d%s", yellow, c.stats.failed.Load(), reset)
	logInfo("RateLimit : %s%d%s", magenta, c.stats.rateLimit.Load(), reset)
	logInfo("Speed     : %s%.1f%s req/s", cyan, float64(c.stats.checked.Load())/elapsed, reset)
	logInfo("Elite     : %s%d%s  •  Working: %s%d%s", orange, c.pool.eliteCount(), reset, green, c.pool.workingCount(), reset)
	logInfo("Sniper    : %sresults/sniper_%s.txt%s", cyan, c.sessionTag, reset)
}

func (c *Checker) selfTest() {
	logInfo("Self-testing Discord API...")
	res := c.check("discord")
	if res.status == 0 {
		logWarn("Self-test failed: %s • %s/%s • %s", res.errMsg, res.scheme, shortProxy(res.proxy), formatLatency(res.latency))
	} else {
		logSuccess("Self-test OK • %s/%s • %s • %s", res.scheme, shortProxy(res.proxy), formatLatency(res.latency), shortBody(res.respBody))
	}
}

func loadProxyLines(path string) []string {
	lines := readLines(path)
	out := make([]string, 0, len(lines))
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" || strings.ContainsAny(ln, " \t") {
			continue
		}
		if !strings.Contains(ln, "://") {
			ln = "http://" + ln
		}
		scheme := strings.ToLower(strings.SplitN(ln, "://", 2)[0])
		switch scheme {
		case "http", "https", "socks4", "socks4a", "socks5", "socks5h":
			out = append(out, ln)
		}
	}
	return out
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	rand.Seed(time.Now().UnixNano())
	proxies := loadProxyLines("proxies.txt")
	printBanner(len(proxies), "SETUP")
	if len(proxies) == 0 {
		logFail("No proxies found in proxies.txt")
		return
	}

	cfg := interactiveConfig()
	smartQuiet = cfg.SmartQuiet

	pool := newProxyPool(proxies, cfg.MaxFailCount, cfg.MinProxyScore)
	c := newChecker(cfg, pool)

	logInfo("Proxies: %s%d%s  •  Mode: %s%s%s  •  Threads: %d  •  SmartQuiet: %v  •  Elite: %d",
		cyan, len(proxies), reset, purple, cfg.Mode, reset, cfg.Threads, smartQuiet, pool.eliteCount())

	if cfg.ProxyCheck {
		alive := c.healthCheck(cfg.ProxyCheckAmount)
		if alive == 0 {
			logFail("Zero working proxies. Exiting.")
			return
		}
	}

	switch strings.ToLower(cfg.Mode) {
	case "generate":
		c.runGenerate()
	case "check":
		usernames := readLines("usernames.txt")
		if len(usernames) == 0 {
			fmt.Printf("%s[?]%s Enter username: ", cyan, reset)
			var u string
			fmt.Scanln(&u)
			u = strings.TrimSpace(u)
			if u != "" {
				usernames = []string{u}
			}
		}
		if len(usernames) > 0 {
			c.selfTest()
			c.runCheck(usernames)
		}
	case "both":
		c.runGenerate()
		c.selfTest()
		c.runCheck(readLines("usernames.txt"))
	case "live":
		c.selfTest()
		c.runLive(max(cfg.UsernamesToGen, 1))
	default:
		logFail("Invalid mode")
	}
}
