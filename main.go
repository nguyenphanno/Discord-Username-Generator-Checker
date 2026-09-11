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
	"unicode"

	"golang.org/x/net/proxy"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	purple  = "\033[38;2;189;147;249m"
	purple2 = "\033[38;2;167;120;255m"
	cyan    = "\033[38;2;139;233;253m"
	cyan2   = "\033[38;2;100;210;255m"
	white   = "\033[38;2;248;248;242m"
	gray    = "\033[38;2;120;130;160m"
	green   = "\033[38;2;80;250;123m"
	red     = "\033[38;2;255;85;85m"
	yellow  = "\033[38;2;255;184;108m"
	orange  = "\033[38;2;255;170;0m"
	magenta = "\033[38;2;255;121;198m"
	version = "Discord<->Username<->Generator<->Checker"
)

var chromeProfiles = []struct {
	ua       string
	version  string
	secChUa  string
	platform string
}{
	{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36", "131.0.0.0", `"Chromium";v="131", "Not;A=Brand";v="24", "Google Chrome";v="131"`, `"Windows"`},
	{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36", "130.0.0.0", `"Chromium";v="130", "Not;A=Brand";v="24", "Google Chrome";v="130"`, `"Windows"`},
	{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36", "131.0.0.0", `"Chromium";v="131", "Not;A=Brand";v="24", "Google Chrome";v="131"`, `"macOS"`},
	{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36", "130.0.0.0", `"Chromium";v="130", "Not;A=Brand";v="24", "Google Chrome";v="130"`, `"Linux"`},
	{"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:132.0) Gecko/20100101 Firefox/132.0", "132.0", "", `"Windows"`},
}

var timezones = []string{
	"America/New_York", "America/Los_Angeles", "Europe/London", "Europe/Paris",
	"Asia/Bangkok", "Asia/Singapore", "Asia/Tokyo", "Asia/Ho_Chi_Minh", "UTC",
}

var locales = []string{"en-US", "en-GB", "vi-VN", "th-TH", "ja-JP", "fr-FR"}

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
	fmt.Printf("%s%s%s  %sINFO %s %s\n", gray, ts(), reset, purple+bold, reset, fmt.Sprintf(format, args...))
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
	fmt.Printf("%s  %s%s%s  •  Proxies: %s%d%s  •  Mode: %s%s%s\n\n",
		gray, bold+purple, version, reset, cyan, proxyCount, reset, purple, strings.ToUpper(mode), reset)
}

func ask(prompt, def string) string {
	fmt.Printf("%s%s%s [%s%s%s]: ", purple, prompt, reset, cyan, def, reset)
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
	mode := strings.ToLower(ask("Mode (live/check/generate/both)", "live"))
	threads := askInt("Threads (recommended 40-80)", 50)
	delay := askFloat("Delay seconds (0 = fastest)", 0.00)
	target := askInt("Target usernames", 5000)
	minLen := askInt("Username min length", 3)
	maxLen := askInt("Username max length", 8)
	if minLen > maxLen {
		minLen, maxLen = maxLen, minLen
	}
	proxyCheck := askBool("Health-check proxies first?", true)
	checkAmount := 200
	if proxyCheck {
		checkAmount = askInt("Proxies to health-check", 200)
	}
	smartQ := askBool("Smart Quiet (only HIT + TAKE)?", true)
	fmt.Println()
	return RuntimeConfig{
		Mode: mode, Threads: threads, Delay: delay, UsernamesToGen: target,
		UsernameMinLength: minLen, UsernameMaxLength: maxLen,
		ProxyCheck: proxyCheck, ProxyCheckAmount: checkAmount,
		MaxFailCount: 3, MinProxyScore: 1, SmartQuiet: smartQ,
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
	idx        atomic.Uint64
}

func newProxyPool(proxies []string, maxFail, minScore int) *ProxyPool {
	p := &ProxyPool{
		proxies: proxies, failCount: make(map[string]int), scores: make(map[string]int),
		successes: make(map[string]int), rateUntil: make(map[string]time.Time),
		disabled: make(map[string]bool), working: make(map[string]bool),
		schemes: make(map[string]string), maxFail: maxFail, minScore: minScore,
		elite: make([]string, 0),
	}
	for _, px := range readLines("working_proxies.txt") {
		p.working[px] = true
		p.scores[px] = 10
		p.successes[px] = 5
		p.elite = append(p.elite, px)
	}
	rand.Shuffle(len(p.proxies), func(i, j int) { p.proxies[i], p.proxies[j] = p.proxies[j], p.proxies[i] })
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
	eliteLen := len(p.elite)
	p.eliteMu.RUnlock()

	if eliteLen > 0 && rand.Float32() < 0.75 {
		p.eliteMu.RLock()
		for i := 0; i < eliteLen; i++ {
			idx := int(p.idx.Add(1)) % eliteLen
			px := p.elite[idx]
			p.eliteMu.RUnlock()
			p.mu.RLock()
			disabled := p.disabled[px]
			limited := false
			if t, ok := p.rateUntil[px]; ok && t.After(now) {
				limited = true
			}
			p.mu.RUnlock()
			if !disabled && !limited {
				return px, true
			}
			p.eliteMu.RLock()
		}
		p.eliteMu.RUnlock()
	}

	p.mu.RLock()
	defer p.mu.RUnlock()
	n := len(p.proxies)
	if n == 0 {
		return "", false
	}
	start := int(p.idx.Add(1)) % n
	for i := 0; i < n; i++ {
		px := p.proxies[(start+i)%n]
		if p.disabled[px] {
			continue
		}
		if t, ok := p.rateUntil[px]; ok && t.After(now) {
			continue
		}
		return px, true
	}
	return "", false
}

func (p *ProxyPool) reportFailure(px string) {
	if px == "" {
		return
	}
	p.mu.Lock()
	p.failCount[px]++
	p.scores[px] = max(p.scores[px]-2, -10)
	if p.failCount[px] >= p.maxFail {
		p.disabled[px] = true
	}
	p.mu.Unlock()
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
	succ := p.successes[px]
	p.mu.Unlock()
	p.successReq.Add(1)
	if !already {
		appendLine("working_proxies.txt", px)
	}
	if succ >= 2 {
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
	if d > 20*time.Minute {
		d = 20 * time.Minute
	}
	p.mu.Lock()
	p.rateUntil[px] = time.Now().Add(d)
	p.mu.Unlock()
}

func (p *ProxyPool) disable(px string) {
	if px == "" {
		return
	}
	p.mu.Lock()
	p.disabled[px] = true
	p.mu.Unlock()
	p.eliteMu.Lock()
	for i, e := range p.elite {
		if e == px {
			p.elite = append(p.elite[:i], p.elite[i+1:]...)
			break
		}
	}
	p.eliteMu.Unlock()
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
	entries, err := os.ReadDir("results")
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if strings.HasPrefix(name, "available_") || strings.HasPrefix(name, "taken_") || strings.HasPrefix(name, "sniper_") {
				for _, u := range readLines(filepath.Join("results", name)) {
					parts := strings.SplitN(u, " | ", 2)
					skip[strings.ToLower(strings.TrimSpace(parts[0]))] = struct{}{}
				}
			}
		}
	}
	return skip
}

type lineWriter struct{ ch chan string }

func newLineWriter(path string) *lineWriter {
	w := &lineWriter{ch: make(chan string, 8192)}
	go func() {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			for range w.ch {
			}
			return
		}
		defer f.Close()
		buf := make([]string, 0, 512)
		ticker := time.NewTicker(150 * time.Millisecond)
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
				if len(buf) >= 512 {
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
	u = strings.TrimSpace(strings.ToLower(u))
	n := len(u)
	if n < 2 || n > 32 {
		return false
	}
	for _, f := range []string{"discord", "everyone", "here", "```", "@", "#", ":"} {
		if strings.Contains(u, f) {
			return false
		}
	}
	for _, r := range u {
		if !(unicode.IsLetter(r) && r < unicode.MaxASCII) && !unicode.IsDigit(r) && r != '_' && r != '.' {
			return false
		}
		if unicode.IsLetter(r) && unicode.IsUpper(r) {
			return false
		}
	}
	if u[0] == '.' || u[n-1] == '.' || strings.Contains(u, "..") {
		return false
	}
	return true
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
	ua, version, secChUa, platform string
}) string {
	osName := "Windows"
	if strings.Contains(profile.platform, "mac") {
		osName = "Mac OS X"
	} else if strings.Contains(profile.platform, "Linux") {
		osName = "Linux"
	}
	props := map[string]any{
		"os": osName, "browser": "Chrome", "device": "",
		"system_locale": locales[rand.Intn(len(locales))], "has_client_mods": false,
		"browser_user_agent": profile.ua, "browser_version": profile.version,
		"os_version": "10", "referrer": "", "referring_domain": "",
		"referrer_current": "", "referring_domain_current": "",
		"release_channel": "stable", "client_build_number": 410000 + rand.Intn(20000),
		"client_event_source": nil, "client_launch_id": randomUUID(),
		"client_app_state": "focused", "client_heartbeat_session_id": randomUUID(),
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
		return "http"
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
		available, taken, failed, checked, rateLimit atomic.Int64
	}
	resultsDir, sessionTag string
}

func newChecker(cfg RuntimeConfig, pool *ProxyPool) *Checker {
	tag := time.Now().Format("20060102_150405")
	os.MkdirAll("results", 0755)
	return &Checker{cfg: cfg, pool: pool, resultsDir: "results", sessionTag: tag}
}

type apiResp struct {
	status     int
	taken      bool
	validJSON  bool
	retryAfter float64
	body       string
	latency    time.Duration
}

func buildPayload(username string) ([]byte, error) {
	return json.Marshal(map[string]string{"username": username})
}

func (c *Checker) doRequest(rt http.RoundTripper, payload []byte) (*apiResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
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
	req.Header.Set("Accept-Language", locale+",en;q=0.8")
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

	start := time.Now()
	client := &http.Client{Transport: rt, Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	latency := time.Since(start)
	if err != nil {
		return &apiResp{latency: latency}, err
	}
	defer resp.Body.Close()
	out := &apiResp{status: resp.StatusCode, latency: latency}
	buf, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	out.body = string(buf)

	if out.status == 200 {
		var j struct {
			Taken bool `json:"taken"`
		}
		if json.Unmarshal(buf, &j) == nil {
			out.taken = j.Taken
			out.validJSON = true
		}
	}
	if out.status == 429 {
		var j struct {
			RetryAfter float64 `json:"retry_after"`
		}
		if json.Unmarshal(buf, &j) == nil && j.RetryAfter > 0 {
			out.retryAfter = j.RetryAfter
		} else {
			out.retryAfter = 4
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
		MaxIdleConns: 2000, MaxIdleConnsPerHost: 100, MaxConnsPerHost: 100,
		IdleConnTimeout: 45 * time.Second, TLSHandshakeTimeout: 4 * time.Second,
		ForceAttemptHTTP2: true,
		DialContext:       (&net.Dialer{Timeout: 4 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
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
			return d.Dial(network, addr)
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
	if c.pool.schemeOf(raw) == "" && scheme != "socks5" {
		if rt2, err2 := c.transportFor(raw, "socks5"); err2 == nil {
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
	status                 int
	errMsg, proxy, scheme  string
	latency                time.Duration
	httpCode               int
	respBody               string
}

func (c *Checker) check(username string) checkResult {
	username = strings.ToLower(strings.TrimSpace(username))
	if !validUsername(username) {
		return checkResult{status: 0, errMsg: "invalid"}
	}
	payload, err := buildPayload(username)
	if err != nil {
		return checkResult{status: 0, errMsg: "payload error"}
	}

	var lastErr, lastProxy, lastScheme string
	var lastLatency time.Duration
	var lastCode int
	var lastBody string

	maxAttempts := 8

	for attempt := 0; attempt < maxAttempts; attempt++ {
		px, ok := c.pool.pick()
		if !ok {
			time.Sleep(250 * time.Millisecond)
			px, ok = c.pool.pick()
			if !ok {
				return checkResult{status: 0, errMsg: "no proxies", proxy: lastProxy, scheme: lastScheme, latency: lastLatency}
			}
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
			continue

		case 400:
			c.pool.disable(px)
			return checkResult{status: 0, errMsg: "bad request", proxy: px, scheme: scheme, latency: res.latency, httpCode: 400, respBody: res.body}

		case 429:
			c.stats.rateLimit.Add(1)
			lastErr = "rate-limit"
			c.pool.reportFailure(px)
			c.pool.rateLimit(px, time.Duration(res.retryAfter*float64(time.Second)))
			continue

		default:
			lastErr = fmt.Sprintf("http%d", res.status)
			c.pool.reportFailure(px)
			continue
		}
	}
	return checkResult{status: 0, errMsg: lastErr, proxy: lastProxy, scheme: lastScheme, latency: lastLatency, httpCode: lastCode, respBody: lastBody}
}

func (c *Checker) healthCheck(amount int) int {
	pxs := c.pool.proxies
	if len(pxs) > amount {
		pxs = pxs[:amount]
	}
	logInfo("Health-checking %s%d%s proxies concurrently...", cyan, len(pxs), reset)
	var wg sync.WaitGroup
	var okCount atomic.Int64
	sem := make(chan struct{}, 120)
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
			ok := false
			for _, sc := range []string{"http", "socks5"} {
				rt, err := c.transportFor(px, sc)
				if err != nil {
					continue
				}
				client := &http.Client{Transport: rt, Timeout: 5 * time.Second}
				if resp, err := client.Do(req); err == nil {
					resp.Body.Close()
					c.pool.setScheme(px, sc)
					ok = true
					break
				}
			}
			if ok {
				okCount.Add(1)
			} else {
				c.pool.disable(px)
			}
		}(px)
	}
	wg.Wait()
	alive := int(okCount.Load())
	logInfo("Health done in %.1fs → %s%d%s alive  •  Elite: %s%d%s  •  Active: %s%d%s",
		time.Since(start).Seconds(), green, alive, reset, cyan, c.pool.eliteCount(), reset, purple, c.pool.activeCount(), reset)
	return alive
}

func (c *Checker) runGenerate() int {
	minL, maxL := c.cfg.UsernameMinLength, c.cfg.UsernameMaxLength
	seen := make(map[string]struct{}, c.cfg.UsernamesToGen)
	jobs := make(chan int, c.cfg.Threads*4)
	results := make(chan string, c.cfg.UsernamesToGen)
	var wg sync.WaitGroup
	for i := 0; i < min(c.cfg.Threads, 40); i++ {
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
	for u := range results {
		seen[u] = struct{}{}
	}
	list := make([]string, 0, len(seen))
	for u := range seen {
		list = append(list, u)
	}
	sort.Strings(list)
	os.WriteFile("usernames.txt", []byte(strings.Join(list, "\n")), 0644)
	logInfo("Saved %s%d%s usernames → usernames.txt", green, len(list), reset)
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
	body = strings.TrimSpace(strings.ReplaceAll(body, "\n", " "))
	if body == "" {
		return ""
	}
	if len(body) > 50 {
		return body[:47] + "..."
	}
	return body
}

func (c *Checker) runCheck(usernames []string) {
	skip := loadSkipSet()
	seen := make(map[string]struct{})
	toCheck := make([]string, 0, len(usernames))
	for _, u := range usernames {
		ul := strings.ToLower(strings.TrimSpace(u))
		if _, ok := seen[ul]; ok || !validUsername(u) {
			continue
		}
		if _, ok := skip[ul]; ok {
			continue
		}
		seen[ul] = struct{}{}
		toCheck = append(toCheck, u)
	}
	total := len(toCheck)
	logSection(fmt.Sprintf("CHECK  •  %d usernames  •  %d threads  •  Elite %d", total, c.cfg.Threads, c.pool.eliteCount()))
	if total == 0 {
		logWarn("Nothing left to check")
		return
	}

	wChecked := newLineWriter("checked.txt")
	wAvail := newLineWriter("available.txt")
	wTaken := newLineWriter("taken.txt")
	wSniper := newLineWriter(filepath.Join(c.resultsDir, "sniper_"+c.sessionTag+".txt"))
	defer wChecked.close()
	defer wAvail.close()
	defer wTaken.close()
	defer wSniper.close()

	jobs := make(chan string, c.cfg.Threads*4)
	results := make(chan result, c.cfg.Threads*4)
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
		body := shortBody(r.res.respBody)
		now := time.Now().Format("2006-01-02 15:04:05")

		switch r.res.status {
		case 1:
			c.stats.available.Add(1)
			logSuccess("%-16s AVAILABLE  %-7s %-20s %-5s %s  [%d/%d]", r.username, scheme, px, lat, body, n, total)
			wAvail.write(r.username)
			wSniper.write(fmt.Sprintf("%s | %s | proxy=%s | scheme=%s | latency=%s", r.username, now, r.res.proxy, scheme, lat))
		case 2:
			c.stats.taken.Add(1)
			logFail("%-16s taken       %-7s %-20s %-5s %s  [%d/%d]", r.username, scheme, px, lat, body, n, total)
			wTaken.write(r.username)
		default:
			c.stats.failed.Add(1)
			logError("%-16s %-10s %-7s %-20s %-5s %s  [%d/%d]", r.username, r.res.errMsg, scheme, px, lat, body, n, total)
		}
	}
	elapsed := time.Since(start).Seconds()
	logSection("SUMMARY")
	logInfo("Time: %.1fs  |  Available: %s%d%s  |  Taken: %s%d%s  |  Failed: %d  |  Speed: %.1f req/s",
		elapsed, green, c.stats.available.Load(), reset, red, c.stats.taken.Load(), reset, c.stats.failed.Load(), float64(c.stats.checked.Load())/elapsed)
	logInfo("Elite: %d  |  Working: %d  |  Sniper: results/sniper_%s.txt", c.pool.eliteCount(), c.pool.workingCount(), c.sessionTag)
}

func (c *Checker) runLive(target int) {
	logSection(fmt.Sprintf("LIVE  •  Target %d  •  %d threads  •  Elite %d", target, c.cfg.Threads, c.pool.eliteCount()))

	jobs := make(chan string, c.cfg.Threads*4)
	results := make(chan result, c.cfg.Threads*4)
	stop := make(chan struct{})
	var once sync.Once
	var wg sync.WaitGroup

	wChecked := newLineWriter("checked.txt")
	wAvail := newLineWriter("available.txt")
	wTaken := newLineWriter("taken.txt")
	wSniper := newLineWriter(filepath.Join(c.resultsDir, "sniper_"+c.sessionTag+".txt"))
	defer wChecked.close()
	defer wAvail.close()
	defer wTaken.close()
	defer wSniper.close()

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
			u := genUsername(rand.Intn(c.cfg.UsernameMaxLength-c.cfg.UsernameMinLength+1) + c.cfg.UsernameMinLength)
			ul := strings.ToLower(u)
			if _, ok := seen.LoadOrStore(ul, struct{}{}); ok {
				continue
			}
			if _, ok := skip[ul]; ok {
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
		body := shortBody(r.res.respBody)
		now := time.Now().Format("2006-01-02 15:04:05")

		switch r.res.status {
		case 1:
			c.stats.available.Add(1)
			logSuccess("%-16s AVAILABLE  %-7s %-20s %-5s %s  [%d/%d]", r.username, scheme, px, lat, body, n, target)
			wAvail.write(r.username)
			wSniper.write(fmt.Sprintf("%s | %s | proxy=%s | scheme=%s | latency=%s", r.username, now, r.res.proxy, scheme, lat))
		case 2:
			c.stats.taken.Add(1)
			logFail("%-16s taken       %-7s %-20s %-5s %s  [%d/%d]", r.username, scheme, px, lat, body, n, target)
			wTaken.write(r.username)
		default:
			c.stats.failed.Add(1)
			logError("%-16s %-10s %-7s %-20s %-5s %s  [%d/%d]", r.username, r.res.errMsg, scheme, px, lat, body, n, target)
		}
	}
	elapsed := time.Since(start).Seconds()
	logSection("SUMMARY")
	logInfo("Time: %.1fs  |  Available: %s%d%s  |  Taken: %s%d%s  |  Failed: %d  |  Speed: %.1f req/s",
		elapsed, green, c.stats.available.Load(), reset, red, c.stats.taken.Load(), reset, c.stats.failed.Load(), float64(c.stats.checked.Load())/elapsed)
	logInfo("Elite: %d  |  Working: %d  |  Sniper: results/sniper_%s.txt", c.pool.eliteCount(), c.pool.workingCount(), c.sessionTag)
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
		if ln == "" {
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

	logInfo("Proxies: %s%d%s  |  Threads: %s%d%s  |  Delay: %.2f  |  SmartQuiet: %v  |  Elite: %s%d%s",
		cyan, len(proxies), reset, cyan, cfg.Threads, reset, cfg.Delay, smartQuiet, cyan, pool.eliteCount(), reset)

	if cfg.ProxyCheck {
		if c.healthCheck(cfg.ProxyCheckAmount) == 0 {
			logFail("No working proxies. Exiting.")
			return
		}
	}

	switch strings.ToLower(cfg.Mode) {
	case "generate":
		c.runGenerate()
	case "check":
		usernames := readLines("usernames.txt")
		if len(usernames) == 0 {
			fmt.Printf("%sUsername: %s", purple, reset)
			var u string
			fmt.Scanln(&u)
			if strings.TrimSpace(u) != "" {
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
