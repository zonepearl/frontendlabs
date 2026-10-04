package site

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// Book is one rendered document, structured for the reader.
type Book struct {
	Title    string
	Subtitle string
	Body     string // HTML: chapters wrapped in <section class="chapter">
	TOC      []TOCEntry
	Words    int
}

type TOCEntry struct {
	Level   int // 1 = part, 2 = chapter
	ID      string
	Text    string
	Minutes int
	URL     string // site path of the page holding this entry (set by the builder)
}

// Minutes is the whole document's reading time at 220 words per minute.
func (b *Book) Minutes() int { return b.Words / 220 }

func (b *Book) Chapters() int {
	n := 0
	for _, e := range b.TOC {
		if e.Level == 2 {
			n++
		}
	}
	return n
}

var (
	reHeading = regexp.MustCompile(`(?s)<h([1-6]) id="([^"]+)">(.*?)</h[1-6]>`)
	reTags    = regexp.MustCompile(`(?s)<[^>]+>`)
	reSpace   = regexp.MustCompile(`\s+`)
	reH1Title = regexp.MustCompile(`(?s)^\s*<h1 id="[^"]+">(.*?)</h1>\s*(<blockquote>.*?</blockquote>)?`)
	reNext    = regexp.MustCompile(`<h[1-6][ >]|<hr>`)
)

func plain(s string) string {
	return strings.TrimSpace(reSpace.ReplaceAllString(html.UnescapeString(reTags.ReplaceAllString(s, " ")), " "))
}

// The kinds of section the guides use, by heading text. Each becomes a
// visually distinct card (same rules as the wiki's own book builder).
var blockKinds = []struct {
	kind string
	re   *regexp.Regexp
}{
	{"tldr", regexp.MustCompile(`^(in one sentence|job$|key takeaways)`)},
	{"problem", regexp.MustCompile(`^(the problem|where we are|the setting|in our story|why this matters|why it exists)`)},
	{"idea", regexp.MustCompile(`^(the idea|the analogy|the insight|the claim|why it matters|the simple version)`)},
	{"how", regexp.MustCompile(`^(how it actually works|how it works|the math|the formula|the mechanism)`)},
	{"example", regexp.MustCompile(`^(worked example|fully worked example|a worked example|real-world example|real-world scenario|worked micro-example|war story)`)},
	{"practice", regexp.MustCompile(`^(practice|try it yourself|experiments|hands-on|build it|step \d)`)},
	{"confusion", regexp.MustCompile(`^(common confusions|common mistakes|what goes wrong|if something goes wrong|when .* goes wrong|what breaks)`)},
	{"check", regexp.MustCompile(`^(check yourself|milestone check|end of part)`)},
	{"reading", regexp.MustCompile(`^(further reading|read next)`)},
	{"series", regexp.MustCompile(`^(across the series|go deeper in the series)`)},
}

func classify(heading string) string {
	h := strings.ToLower(plain(heading))
	for _, k := range blockKinds {
		if k.re.MatchString(h) {
			return k.kind
		}
	}
	return ""
}

// buildBook turns rendered Markdown into the reader's structure.
func buildBook(body string) *Book {
	b := &Book{}

	// Title page: the first H1 and the blockquote that introduces it.
	if m := reH1Title.FindStringSubmatchIndex(body); m != nil {
		b.Title = plain(body[m[2]:m[3]])
		if m[4] >= 0 {
			b.Subtitle = firstSentences(plain(body[m[4]:m[5]]), 240)
		}
		body = body[m[1]:]
	}

	body = addGitHubAliases(body)
	body = wrapBlocks(body)
	body = strings.ReplaceAll(body, "<table>", `<div class="tbl"><table>`)
	body = strings.ReplaceAll(body, "</table>", "</table></div>")

	// Contents (parts and chapters) and reading time per chapter.
	locs := reHeading.FindAllStringSubmatchIndex(body, -1)
	for i, m := range locs {
		level := int(body[m[2]] - '0')
		end := len(body)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		words := len(strings.Fields(plain(body[m[1]:end])))
		b.Words += words
		if level <= 2 {
			b.TOC = append(b.TOC, TOCEntry{Level: level, ID: body[m[4]:m[5]], Text: plain(body[m[6]:m[7]])})
		}
		// a chapter's time includes all of its sub-sections
		for j := len(b.TOC) - 1; j >= 0; j-- {
			if b.TOC[j].Level == 2 {
				b.TOC[j].Minutes += words
				break
			}
			if b.TOC[j].Level == 1 {
				break
			}
		}
	}
	for i := range b.TOC {
		b.TOC[i].Minutes = max(1, (b.TOC[i].Minutes+110)/220) // words -> minutes at 220 wpm
	}

	b.Body = wrapChapters(body)
	return b
}

// wrapBlocks puts each recognised h3 and its content (up to the next heading
// or horizontal rule) into <section class="blk blk-KIND">.
func wrapBlocks(body string) string {
	var out strings.Builder
	rest := body
	for {
		m := reHeading.FindStringSubmatchIndex(rest)
		for m != nil && rest[m[2]] != '3' { // only h3s start blocks
			out.WriteString(rest[:m[1]])
			rest = rest[m[1]:]
			m = reHeading.FindStringSubmatchIndex(rest)
		}
		if m == nil {
			out.WriteString(rest)
			return out.String()
		}
		kind := classify(rest[m[6]:m[7]])
		out.WriteString(rest[:m[0]])
		if kind == "" {
			out.WriteString(rest[m[0]:m[1]])
			rest = rest[m[1]:]
			continue
		}
		heading := fmt.Sprintf(`<h3 id="%s" class="k-%s">%s</h3>`, rest[m[4]:m[5]], kind, rest[m[6]:m[7]])
		after := rest[m[1]:]
		end := len(after)
		if n := reNext.FindStringIndex(after); n != nil {
			end = n[0]
		}
		fmt.Fprintf(&out, `<section class="blk blk-%s">%s%s</section>`, kind, heading, after[:end])
		rest = after[end:]
	}
}

// wrapChapters wraps every h2 and what follows it in a chapter section.
func wrapChapters(body string) string {
	var out strings.Builder
	parts := strings.Split(body, `<h2 id="`)
	out.WriteString(parts[0])
	for _, p := range parts[1:] {
		id := p[:strings.IndexByte(p, '"')]
		// A part heading (h1) inside this chunk ends the chapter.
		inner, tail := p, ""
		if i := strings.Index(p, "<h1 id="); i >= 0 {
			inner, tail = p[:i], p[i:]
		}
		fmt.Fprintf(&out, `<section class="chapter" id="ch-%s" data-ch="%s"><h2 id="%s</section>%s`, id, id, inner, tail)
	}
	return out.String()
}

func firstSentences(s string, limit int) string {
	var out string
	for _, sent := range regexp.MustCompile(`(?:[.!?])\s+`).Split(s, -1) {
		sent = strings.TrimSpace(sent)
		if sent == "" {
			continue
		}
		if out != "" && len(out)+len(sent) > limit {
			break
		}
		if out != "" {
			out += ". "
		}
		out += strings.TrimRight(sent, ".")
	}
	if out != "" && !strings.HasSuffix(out, "?") && !strings.HasSuffix(out, "!") {
		out += "."
	}
	return out
}

var reGHStrip = regexp.MustCompile(`[^\p{L}\p{N} _-]`)

// githubSlug is the anchor GitHub generates for a heading: lowercase, drop
// punctuation, spaces -> "-" (so "A & B" becomes "a--b").
func githubSlug(text string) string {
	return strings.ReplaceAll(reGHStrip.ReplaceAllString(strings.ToLower(text), ""), " ", "-")
}

// addGitHubAliases adds an empty anchor with the GitHub-style ID inside any
// heading whose GitHub slug differs from its own ID, so links written for
// GitHub's rendering ("#3-latency--load-testing") also work here.
func addGitHubAliases(body string) string {
	ids := map[string]bool{}
	for _, m := range reHeading.FindAllStringSubmatch(body, -1) {
		ids[m[2]] = true
	}
	return reHeading.ReplaceAllStringFunc(body, func(h string) string {
		m := reHeading.FindStringSubmatch(h)
		gh := githubSlug(plain(m[3]))
		if gh == "" || gh == m[2] || ids[gh] {
			return h
		}
		ids[gh] = true
		open := fmt.Sprintf(`<h%s id="%s">`, m[1], m[2])
		return open + `<span class="anchor-alias" id="` + gh + `"></span>` + h[len(open):]
	})
}

// ChapterPage is one chapter rendered as its own page (/<doc>/<chapter-id>/).
type ChapterPage struct {
	ID      string // the chapter heading's id; also the URL segment
	Title   string
	Part    string // title of the part (h1) the chapter belongs to, if any
	HTML    string // the part's introduction (first chapter of a part only) + the chapter section
	Words   int
	Minutes int
	Desc    string // first paragraph, trimmed, for <meta name="description">
	URL     string // site path, set by the builder
}

var (
	reChapterStart = regexp.MustCompile(`<section class="chapter" id="ch-`)
	rePartH1       = regexp.MustCompile(`(?s)<h1 id="[^"]+">(.*?)</h1>`)
	reFirstPara    = regexp.MustCompile(`(?s)<(?:p|li)>(.*?)</(?:p|li)>`) // opening prose: paragraphs and list items
	reChapterH2    = regexp.MustCompile(`(?s)^(<section class="chapter" id="ch-[^"]+" data-ch="[^"]+">)<h2 id="([^"]+)">(.*?)</h2>`)
)

// Split cuts the book body into its front matter (everything before the first
// part or chapter) and one page per chapter. A part's heading and introduction
// go on the page of the part's first chapter, so nothing is lost; a part with
// no chapters becomes a page of its own.
func (b *Book) Split() (front string, pages []*ChapterPage) {
	body := b.Body
	locs := reChapterStart.FindAllStringIndex(body, -1)
	if len(locs) == 0 {
		return body, nil
	}
	front = body[:locs[0][0]]
	prelude := ""
	if i := strings.Index(front, "<h1 id="); i >= 0 { // a part starts before the first chapter
		front, prelude = front[:i], front[i:]
	}
	part := ""
	if m := rePartH1.FindAllStringSubmatch(prelude, -1); m != nil {
		part = plain(m[len(m)-1][1])
	}
	for i, l := range locs {
		end := len(body)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		chunk := body[l[0]:end]
		// the chapter section ends at "</section>" followed by the next part's heading, if any
		section, tail := chunk, ""
		if j := strings.Index(chunk, "</section><h1 id="); j >= 0 {
			section, tail = chunk[:j+len("</section>")], chunk[j+len("</section>"):]
		}
		m := reChapterH2.FindStringSubmatch(section)
		if m == nil {
			continue
		}
		p := &ChapterPage{ID: m[2], Title: plain(m[3]), Part: part}
		// the chapter title is the page's h1
		section = m[1] + `<h1 id="` + m[2] + `" class="chapter-title">` + m[3] + `</h1>` + section[len(m[0]):]
		p.HTML = prelude + section
		p.Words = len(strings.Fields(plain(p.HTML)))
		p.Minutes = max(1, (p.Words+110)/220)
		p.Desc = describe(section)
		pages = append(pages, p)
		prelude = tail
		if pm := rePartH1.FindAllStringSubmatch(tail, -1); pm != nil {
			part = plain(pm[len(pm)-1][1])
		}
	}
	if strings.TrimSpace(prelude) != "" { // trailing parts without chapters
		if m := regexp.MustCompile(`<h1 id="([^"]+)">(.*?)</h1>`).FindStringSubmatch(prelude); m != nil {
			p := &ChapterPage{ID: m[1], Title: plain(m[2]), Part: plain(m[2]), HTML: prelude}
			p.Words = len(strings.Fields(plain(prelude)))
			p.Minutes = max(1, (p.Words+110)/220)
			pages = append(pages, p)
		} else if len(pages) > 0 {
			pages[len(pages)-1].HTML += prelude
		}
	}
	return front, pages
}

// IDs lists every element id in an HTML fragment (headings, aliases, anything).
func IDs(html string) []string {
	var out []string
	for _, m := range reIDAttr.FindAllStringSubmatch(html, -1) {
		out = append(out, m[1])
	}
	return out
}

// trimWords shortens s to at most limit bytes at a word boundary, adding "…".
func trimWords(s string, limit int) string {
	if len(s) <= limit {
		return s
	}
	cut := strings.LastIndex(s[:limit], " ")
	if cut < limit/2 {
		cut = limit
	}
	return strings.TrimRight(s[:cut], " ,;:—-") + "…"
}

// describe builds a meta description from a chapter's opening paragraphs:
// enough text to fill the ~155 characters search results show.
func describe(section string) string {
	desc := ""
	for _, pm := range reFirstPara.FindAllStringSubmatch(section, 8) {
		t := plain(pm[1])
		if t == "" {
			continue
		}
		if desc != "" {
			desc += " "
		}
		desc += t
		if len(desc) >= 120 {
			break
		}
	}
	return trimWords(reSpacePunct.ReplaceAllString(desc, "$1"), 155)
}

// reSpacePunct removes the space that stripping inline tags leaves before punctuation ("salt ," -> "salt,").
var reSpacePunct = regexp.MustCompile(` ([,.;:!?)])`)
