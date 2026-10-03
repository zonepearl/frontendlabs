package site

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

// Watcher regenerates the site whenever its inputs change: the original
// guides in the wiki (synced into docs/), docs/ itself, and -- in dev mode --
// the templates and CSS.
type Watcher struct {
	B        *Builder
	SyncWiki bool   // also watch the wiki originals and copy changes into docs/
	WebDir   string // non-empty in dev mode: watch templates/static on disk
	Interval time.Duration
	OnChange func(urls []string) // called after each successful rebuild
}

func (w *Watcher) Run(ctx context.Context) {
	m := w.B.M
	wikiFP := takeFingerprint(m.WikiRoot(), m.Files())
	docsFP := takeFingerprint(w.B.DocsDir, m.Files())
	var webFP fingerprint
	if w.WebDir != "" {
		webFP = treeFingerprint(w.WebDir)
	}
	t := time.NewTicker(w.Interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
		if w.SyncWiki {
			if now := takeFingerprint(m.WikiRoot(), m.Files()); len(wikiFP.diff(now)) > 0 {
				wikiFP = now
				changed, err := Sync(m, w.B.DocsDir)
				if err != nil {
					log.Printf("sync: %v", err)
				}
				for _, c := range changed {
					fmt.Printf("synced  %s\n", c)
				}
			}
		}
		if now := takeFingerprint(w.B.DocsDir, m.Files()); len(docsFP.diff(now)) > 0 {
			changed := docsFP.diff(now)
			docsFP = now
			urls, err := w.B.Rebuild(changed)
			if err != nil {
				log.Printf("rebuild: %v", err) // keep serving the last good version
				continue
			}
			if w.OnChange != nil && len(urls) > 0 {
				w.OnChange(urls)
			}
		}
		if w.WebDir != "" {
			if now := treeFingerprint(w.WebDir); len(webFP.diff(now)) > 0 {
				webFP = now
				if err := w.B.RebuildTheme(); err != nil {
					log.Printf("theme: %v", err)
					continue
				}
				if w.OnChange != nil {
					w.OnChange([]string{"*"})
				}
			}
		}
	}
}

// hub fans live-reload events out to every connected browser tab.
type hub struct {
	mu      sync.Mutex
	clients map[chan string]bool
}

func (h *hub) broadcast(urls []string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.clients {
		for _, u := range urls {
			select {
			case c <- u:
			default: // a slow tab misses an event rather than blocking the build
			}
		}
	}
}

func (h *hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	rc := http.NewResponseController(w)
	c := make(chan string, 16)
	h.mu.Lock()
	h.clients[c] = true
	h.mu.Unlock()
	defer func() { h.mu.Lock(); delete(h.clients, c); h.mu.Unlock() }()
	fmt.Fprint(w, "retry: 1000\n\n")
	rc.Flush()
	ping := time.NewTicker(20 * time.Second)
	defer ping.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case u := <-c:
			fmt.Fprintf(w, "data: %s\n\n", u)
		case <-ping.C:
			fmt.Fprint(w, ": ping\n\n") // keeps proxies from closing an idle stream
		}
		rc.SetWriteDeadline(time.Now().Add(30 * time.Second))
		if err := rc.Flush(); err != nil {
			return
		}
	}
}

// Serve builds the site, serves it, and regenerates it on every change.
func Serve(ctx context.Context, w *Watcher, addr string) error {
	h := &hub{clients: map[chan string]bool{}}
	w.OnChange = func(urls []string) {
		slices.Sort(urls)
		fmt.Printf("live-reload %s\n", strings.Join(urls, " "))
		h.broadcast(urls)
	}
	go w.Run(ctx)

	mux := http.NewServeMux()
	mux.Handle("/__live", h)
	files := http.FileServer(http.Dir(w.B.OutDir))
	mux.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Cache-Control", "no-cache") // always revalidate: pages change under us
		files.ServeHTTP(rw, r)
	}))
	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() { <-ctx.Done(); srv.Shutdown(context.Background()) }()
	fmt.Printf("serving %s on http://localhost%s  (watching for changes; Ctrl-C to stop)\n", filepath.Base(w.B.OutDir), addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}
