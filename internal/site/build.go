package site

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"
)

// Builder renders docs/ into a static site in out/.
type Builder struct {
	M       *Manifest
	DocsDir string
	OutDir  string
	Assets  fs.FS // contains templates/ and static/
	Live    bool  // inject the live-reload client (serve mode)

	mu     sync.Mutex
	tmpl   *template.Template
	metas  map[string]*docMeta        // by docs-relative source path
	appIDs map[string]map[string]bool // element ids on each rendered page, by site URL
}

// docMeta is what the home page and search index need from each document.
type docMeta struct {
	Title    string
	Chapters int
	Words    int
	Labs     int
	TOC      []TOCEntry
}

func (b *Builder) docMaps() (byMD, byOrg map[string]string) {
	byMD, byOrg = map[string]string{}, map[string]string{}
	for _, d := range b.M.Docs() {
		byMD[d.Source] = d.URL
	}
	for _, g := range b.M.Guides {
		for _, o := range g.Originals {
			byOrg[path.Clean(o)] = g.Slug + "/"
		}
	}
	return
}

func (b *Builder) loadTemplates() error {
	t, err := template.New("").Funcs(template.FuncMap{
		"rel":     relURL,
		"minutes": formatMinutes,
		"base":    path.Base,
		"lower":   strings.ToLower,
		"year":    func() int { return time.Now().Year() },
	}).ParseFS(b.Assets, "templates/*.html")
	if err != nil {
		return err
	}
	b.tmpl = t
	return nil
}

// BuildAll renders every document, the home page, and copies assets.
func (b *Builder) BuildAll() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	start := time.Now()
	if err := b.loadTemplates(); err != nil {
		return err
	}
	b.metas = map[string]*docMeta{}
	b.appIDs = map[string]map[string]bool{}
	if err := b.copyAssets(); err != nil {
		return err
	}
	for _, d := range b.M.Docs() {
		if err := b.buildDoc(d); err != nil {
			return fmt.Errorf("%s: %w", d.Source, err)
		}
	}
	if err := b.copyOriginals(); err != nil {
		return err
	}
	if err := b.buildHome(); err != nil {
		return err
	}
	if err := b.writeHostingFiles(); err != nil {
		return err
	}
	fmt.Printf("built %d documents in %v -> %s\n", len(b.M.Docs()), time.Since(start).Round(time.Millisecond), b.OutDir)
	return nil
}

// Rebuild re-renders only the documents whose sources changed, then the
// home page (titles, chapter counts and reading time may have changed). It
// returns the site URLs that changed, for live reload.
func (b *Builder) Rebuild(changedSources []string) ([]string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	set := map[string]bool{}
	for _, s := range changedSources {
		set[path.Clean(filepath.ToSlash(s))] = true
	}
	var urls []string
	for _, d := range b.M.Docs() {
		if set[d.Source] {
			start := time.Now()
			if err := b.buildDoc(d); err != nil {
				return urls, fmt.Errorf("%s: %w", d.Source, err)
			}
			fmt.Printf("rebuilt %-48s in %v\n", d.Source, time.Since(start).Round(time.Millisecond))
			urls = append(urls, "/"+d.URL)
		}
	}
	if err := b.copyOriginals(); err != nil {
		return urls, err
	}
	if len(urls) > 0 {
		urls = append(urls, "/")
		return urls, b.buildHome()
	}
	return urls, nil
}

// RebuildTheme reloads templates and assets and re-renders everything.
func (b *Builder) RebuildTheme() error {
	return b.BuildAll()
}

type pageLink struct {
	Title, URL string
	Current    bool
}

type bookPage struct {
	Site      *Manifest
	Guide     *Guide
	Doc       *Doc
	Book      *Book
	Body      template.HTML
	PageURL   string
	Pages     []pageLink // for guides with sub-pages
	Originals []pageLink
	SourceURL string
	Prev      *Guide
	Next      *Guide
	Live      bool
	Built     string
}

func (b *Builder) buildDoc(d *Doc) error {
	src, err := os.ReadFile(filepath.Join(b.DocsDir, filepath.FromSlash(d.Source)))
	if err != nil {
		return err
	}
	byMD, byOrg := b.docMaps()
	htmlBody, err := renderMarkdown(src, &linkRewriter{doc: d, byMD: byMD, byOrg: byOrg})
	if err != nil {
		return err
	}
	book := buildBook(htmlBody)
	if book.Title == "" {
		book.Title = d.Guide.Title
	}
	b.metas[d.Source] = &docMeta{
		Title: book.Title, Chapters: book.Chapters(), Words: book.Words, TOC: book.TOC,
		Labs: countPrograms(src),
	}

	p := &bookPage{
		Site: b.M, Guide: d.Guide, Doc: d, Book: book, Body: template.HTML(book.Body),
		PageURL: d.URL, Live: b.Live, Built: time.Now().Format("2 Jan 2006 15:04"),
		SourceURL: relURL(d.URL, "docs/"+d.Source),
	}
	if len(d.Guide.Pages) > 0 {
		p.Pages = append(p.Pages, pageLink{"Overview", relURL(d.URL, d.Guide.Slug+"/"), d.Page == nil})
		for _, pg := range d.Guide.Pages {
			title := pg.Slug
			if m := b.metas[path.Clean(pg.Source)]; m != nil {
				title = m.Title
			} else if t := firstHeading(filepath.Join(b.DocsDir, pg.Source)); t != "" {
				title = t
			}
			p.Pages = append(p.Pages, pageLink{title, relURL(d.URL, d.Guide.Slug+"/"+pg.Slug+"/"), d.Page == pg})
		}
	}
	for _, o := range d.Guide.Originals {
		label := "Original book edition"
		if !strings.HasSuffix(o, "-book.html") {
			label = "Original HTML edition"
		}
		p.Originals = append(p.Originals, pageLink{label, relURL(d.URL, "originals/"+o), false})
	}
	p.Prev, p.Next = b.neighbours(d.Guide)
	var page strings.Builder
	if err := b.tmpl.ExecuteTemplate(&page, "book.html", p); err != nil {
		return err
	}
	b.appIDs[d.URL] = idSet(page.String())
	if err := writeFile(filepath.Join(b.OutDir, filepath.FromSlash(d.URL+"index.html")), []byte(page.String())); err != nil {
		return err
	}
	return b.copyFile(filepath.Join(b.DocsDir, d.Source), filepath.Join(b.OutDir, "docs", d.Source))
}

// neighbours returns the previous/next guide. In the classic layout that is the
// guide's category order (the series order for the systems series). In the
// spotlight layout, guides on the course path follow the path even though the
// home page shows them in different topic sections; other guides follow their
// section's order.
func (b *Builder) neighbours(g *Guide) (prev, next *Guide) {
	var order []string
	switch {
	case b.M.Spotlit() && slices.Contains(b.M.Path, g.Slug):
		order = b.M.Path
	case b.M.Spotlit():
		if sec := b.M.SectionOf(g.Slug); sec != nil {
			order = sec.Guides
		}
	default:
		for _, x := range b.M.Guides {
			if x.Category == g.Category {
				order = append(order, x.Slug)
			}
		}
	}
	if i := slices.Index(order, g.Slug); i >= 0 {
		if i > 0 {
			prev = b.M.Guide(order[i-1])
		}
		if i+1 < len(order) {
			next = b.M.Guide(order[i+1])
		}
	}
	return
}

type homeCard struct {
	*Guide
	URL      string
	Chapters int
	Minutes  int
	Labs     int
}

type homeSection struct {
	Category
	Cards []homeCard
}

// spotView is the spotlight hero: the featured guide, deep links into it, the
// other guides of its section, and topics that are planned but not written.
type spotView struct {
	Eyebrow string
	Card    homeCard
	Links   []struct{ Title, URL string }
	Extras  []homeCard
	Coming  []string
}

type homePage struct {
	Site     *Manifest
	Sections []homeSection // classic: categories. spotlight: topic sections after the spotlight
	Totals   struct{ Guides, Chapters, Hours, Labs int }
	PageURL  string
	Live     bool
	Built    string

	// Spotlight layout only.
	Spotlight *spotView
	Path      []homeCard // the course, in reading order
	Alongside []homeCard // guides read alongside the course (step "Alongside")
}

func (b *Builder) buildHome() error {
	hp := &homePage{Site: b.M, Live: b.Live, Built: time.Now().Format("2 Jan 2006 15:04")}
	words := 0
	cards := map[string]homeCard{}
	for _, g := range b.M.Guides {
		card := homeCard{Guide: g, URL: g.Slug + "/"}
		for _, d := range b.M.Docs() {
			if d.Guide == g {
				if m := b.metas[d.Source]; m != nil {
					card.Chapters += m.Chapters
					card.Minutes += m.Words / 220
					card.Labs += m.Labs
					words += m.Words
				}
			}
		}
		cards[g.Slug] = card
		hp.Totals.Guides++
		hp.Totals.Chapters += card.Chapters
		hp.Totals.Labs += card.Labs
	}
	hp.Totals.Hours = words / 220 / 60

	page := "home.html"
	if b.M.Spotlit() {
		page = "home-spotlight.html"
		sp := b.M.Spotlight
		hp.Spotlight = &spotView{Eyebrow: sp.Eyebrow, Card: cards[sp.Guide], Coming: sp.Coming}
		for _, l := range sp.Links {
			hp.Spotlight.Links = append(hp.Spotlight.Links, struct{ Title, URL string }{l.Title, sp.Guide + "/#" + l.Anchor})
		}
		home := b.M.SectionOf(sp.Guide)
		for _, sec := range b.M.Sections {
			if home != nil && sec.ID == home.ID {
				for _, slug := range sec.Guides {
					if slug != sp.Guide {
						hp.Spotlight.Extras = append(hp.Spotlight.Extras, cards[slug])
					}
				}
				continue
			}
			hs := homeSection{Category: Category{ID: sec.ID, Title: sec.Title, Note: sec.Note}}
			for _, slug := range sec.Guides {
				hs.Cards = append(hs.Cards, cards[slug])
			}
			hp.Sections = append(hp.Sections, hs)
		}
		for _, slug := range b.M.Path {
			hp.Path = append(hp.Path, cards[slug])
		}
		for _, g := range b.M.Guides {
			if g.Step == "Alongside" {
				hp.Alongside = append(hp.Alongside, cards[g.Slug])
			}
		}
	} else {
		for _, c := range b.M.Categories {
			sec := homeSection{Category: c}
			for _, g := range b.M.Guides {
				if g.Category == c.ID {
					sec.Cards = append(sec.Cards, cards[g.Slug])
				}
			}
			hp.Sections = append(hp.Sections, sec)
		}
	}
	if err := b.render(page, "index.html", hp); err != nil {
		return err
	}
	return b.writeSearchIndex()
}

// writeSearchIndex lists every guide, part, and chapter for the home search.
func (b *Builder) writeSearchIndex() error {
	type entry struct {
		T string `json:"t"` // title
		U string `json:"u"` // URL
		G string `json:"g"` // guide title
		I string `json:"i"` // icon
	}
	var idx []entry
	for _, d := range b.M.Docs() {
		m := b.metas[d.Source]
		if m == nil {
			continue
		}
		idx = append(idx, entry{m.Title, d.URL, d.Guide.Title, d.Guide.Icon})
		for _, e := range m.TOC {
			idx = append(idx, entry{e.Text, d.URL + "#" + e.ID, d.Guide.Title, d.Guide.Icon})
		}
	}
	js, _ := json.Marshal(idx)
	return writeFile(filepath.Join(b.OutDir, "search.json"), js)
}

// writeHostingFiles adds what a public static host expects: a sitemap with
// absolute URLs, robots.txt, a 404 page, and a CNAME for the custom domain.
func (b *Builder) writeHostingFiles() error {
	if b.M.Site.BaseURL == "" {
		return nil
	}
	var sm strings.Builder
	sm.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	today := time.Now().Format("2006-01-02")
	fmt.Fprintf(&sm, "  <url><loc>%s</loc><lastmod>%s</lastmod></url>\n", b.M.Abs(""), today)
	for _, d := range b.M.Docs() {
		lastmod := today
		if st, err := os.Stat(filepath.Join(b.DocsDir, d.Source)); err == nil {
			lastmod = st.ModTime().Format("2006-01-02")
		}
		fmt.Fprintf(&sm, "  <url><loc>%s</loc><lastmod>%s</lastmod></url>\n", b.M.Abs(d.URL), lastmod)
	}
	sm.WriteString("</urlset>\n")
	if err := writeFile(filepath.Join(b.OutDir, "sitemap.xml"), []byte(sm.String())); err != nil {
		return err
	}
	robots := "User-agent: *\nAllow: /\nDisallow: /docs/\nDisallow: /originals/\n\nSitemap: " + b.M.Abs("sitemap.xml") + "\n"
	if err := writeFile(filepath.Join(b.OutDir, "robots.txt"), []byte(robots)); err != nil {
		return err
	}
	host := strings.TrimPrefix(strings.TrimPrefix(b.M.Site.BaseURL, "https://"), "http://")
	if err := writeFile(filepath.Join(b.OutDir, "CNAME"), []byte(strings.Trim(host, "/")+"\n")); err != nil {
		return err
	}
	return b.render("404.html", "404.html", nil)
}

func (b *Builder) render(name, out string, data any) error {
	var sb strings.Builder
	if err := b.tmpl.ExecuteTemplate(&sb, name, data); err != nil {
		return err
	}
	return writeFile(filepath.Join(b.OutDir, filepath.FromSlash(out)), []byte(sb.String()))
}

func (b *Builder) copyAssets() error {
	return fs.WalkDir(b.Assets, "static", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := fs.ReadFile(b.Assets, p)
		if err != nil {
			return err
		}
		return writeFile(filepath.Join(b.OutDir, "assets", strings.TrimPrefix(p, "static/")), data)
	})
}

func (b *Builder) copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

func writeFile(p string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	tmp := p + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, p) // readers never see a half-written page
}

func firstHeading(file string) string {
	b, err := os.ReadFile(file)
	if err != nil {
		return ""
	}
	for _, line := range strings.SplitN(string(b), "\n", 20) {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(line[2:])
		}
	}
	return ""
}

var reGoBlock = regexp.MustCompile("(?ms)^```go\\n(.*?)^```")

// countPrograms counts complete Go programs (fenced go blocks with
// "package main") -- the runnable labs, not snippets or pointers to them.
func countPrograms(md []byte) int {
	n := 0
	for _, m := range reGoBlock.FindAllSubmatch(md, -1) {
		if strings.Contains(string(m[1]), "package main") {
			n++
		}
	}
	return n
}

func formatMinutes(m int) string {
	switch {
	case m < 60:
		return fmt.Sprintf("%d min", max(1, m))
	case m%60 == 0:
		return fmt.Sprintf("%d h", m/60)
	default:
		return fmt.Sprintf("%d h %d min", m/60, m%60)
	}
}
