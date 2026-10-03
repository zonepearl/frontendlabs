package site

import (
	"bytes"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// slugIDs generates heading anchors exactly like the wiki's Python builders
// (lowercase, runs of non-alphanumerics -> "-", "-1", "-2" for duplicates),
// so every existing cross-guide link "#chapter-21-..." keeps working.
type slugIDs struct{ seen map[string]int }

var nonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func (s *slugIDs) Generate(value []byte, _ ast.NodeKind) []byte {
	base := strings.ToLower(string(value))
	base = strings.ReplaceAll(base, "§", "sec")
	base = strings.Trim(nonAlnum.ReplaceAllString(base, "-"), "-")
	if base == "" {
		base = "section"
	}
	n := s.seen[base]
	s.seen[base] = n + 1
	if n == 0 {
		return []byte(base)
	}
	return []byte(base + "-" + strconv.Itoa(n))
}

func (s *slugIDs) Put(value []byte) { s.seen[string(value)]++ }

// linkRewriter maps links between guides' .md files to the site's URLs.
type linkRewriter struct {
	doc   *Doc
	byMD  map[string]string // docs-relative .md path -> site URL
	byOrg map[string]string // docs-relative original .html path -> site URL
}

var scheme = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*:`)

func (lr *linkRewriter) Transform(doc *ast.Document, _ text.Reader, _ parser.Context) {
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if l, ok := n.(*ast.Link); ok && entering {
			l.Destination = []byte(lr.rewrite(string(l.Destination)))
		}
		return ast.WalkContinue, nil
	})
}

// byDir finds the guide that owns documents under dir (e.g. the Go plan's
// detailed-90-day-plan/ folder).
func (lr *linkRewriter) byDir(dir string) (string, bool) {
	for src, url := range lr.byMD {
		if strings.HasPrefix(src, dir+"/") {
			return url[:strings.Index(url, "/")+1], true
		}
	}
	return "", false
}

func (lr *linkRewriter) rewrite(dest string) string {
	if dest == "" || strings.HasPrefix(dest, "#") || scheme.MatchString(dest) {
		return dest
	}
	p, frag, hasFrag := strings.Cut(dest, "#")
	target := path.Clean(path.Join(path.Dir(lr.doc.Source), p))
	url, ok := lr.byMD[target]
	if !ok {
		url, ok = lr.byOrg[target]
	}
	if !ok && path.Ext(target) == "" { // a directory: its README page, or the guide whose pages live there
		if url, ok = lr.byMD[target+"/README.md"]; !ok {
			url, ok = lr.byDir(target)
		}
	}
	if !ok {
		if strings.HasSuffix(p, ".md") { // a document the site doesn't render: link the raw copy
			url = "docs/" + target
		} else {
			return dest
		}
	}
	out := relURL(lr.doc.URL, url)
	if hasFrag {
		out += "#" + frag
	}
	return out
}

// relURL returns a link from a page at site path `from` to site path `to`,
// using only relative paths so the site works from any base URL or from disk.
func relURL(from, to string) string {
	if to == from {
		return "./"
	}
	return strings.Repeat("../", strings.Count(from, "/")) + to
}

// renderMarkdown converts one document to HTML with GFM tables, task lists,
// autolinks, raw HTML (e.g. <details>), and wiki-compatible heading IDs.
func renderMarkdown(src []byte, lr *linkRewriter) (string, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithASTTransformers(util.Prioritized(lr, 100)),
		),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	var buf bytes.Buffer
	ctx := parser.NewContext(parser.WithIDs(&slugIDs{seen: map[string]int{}}))
	if err := md.Convert(src, &buf, parser.WithContext(ctx)); err != nil {
		return "", err
	}
	return buf.String(), nil
}
