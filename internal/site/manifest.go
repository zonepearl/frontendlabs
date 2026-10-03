// Package site turns the Markdown guides in docs/ into a book-style static site.
package site

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Manifest is guides.json: which documents exist, and how they're presented.
type Manifest struct {
	Site struct {
		Title    string `json:"title"`
		Tagline  string `json:"tagline"`
		BaseURL  string `json:"base_url"`  // e.g. https://frontendlabs.xyz -- canonical URLs, sitemap, social cards
		WikiRoot string `json:"wiki_root"` // where the original guides live, relative to the manifest
	} `json:"site"`
	Categories []Category `json:"categories"`
	Guides     []*Guide   `json:"guides"`

	dir string // directory containing guides.json
}

type Category struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Note  string `json:"note"`
}

type Guide struct {
	Slug      string   `json:"slug"`
	Title     string   `json:"title"`
	Icon      string   `json:"icon"`
	Category  string   `json:"category"`
	Step      string   `json:"step"`
	Level     string   `json:"level"`
	Summary   string   `json:"summary"`
	Source    string   `json:"source"`    // Markdown path, relative to docs/ (and to the wiki root)
	Originals []string `json:"originals"` // the wiki's own HTML editions, copied as-is
	Pages     []*Page  `json:"pages"`     // extra documents rendered as sub-pages of this guide
}

type Page struct {
	Slug   string `json:"slug"`
	Source string `json:"source"`
}

// Doc is one renderable Markdown document and its place in the site.
type Doc struct {
	Guide  *Guide
	Page   *Page  // nil for the guide's main document
	Source string // path relative to docs/
	URL    string // site path, always ending in "/" (e.g. "tcp-ip/" or "go-plan/week5/")
}

func LoadManifest(file string) (*Manifest, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	m.dir = filepath.Dir(file)
	seen := map[string]bool{}
	for _, g := range m.Guides {
		if seen[g.Slug] {
			return nil, fmt.Errorf("duplicate guide slug %q", g.Slug)
		}
		seen[g.Slug] = true
	}
	return &m, nil
}

// Abs turns a site path ("tcp-ip/") into an absolute URL on the production domain.
func (m *Manifest) Abs(p string) string {
	return strings.TrimRight(m.Site.BaseURL, "/") + "/" + strings.TrimLeft(p, "/")
}

// WikiRoot is the absolute path of the original guides.
func (m *Manifest) WikiRoot() string { return filepath.Join(m.dir, m.Site.WikiRoot) }

// Docs lists every Markdown document the site renders.
func (m *Manifest) Docs() []*Doc {
	var out []*Doc
	for _, g := range m.Guides {
		out = append(out, &Doc{Guide: g, Source: path.Clean(g.Source), URL: g.Slug + "/"})
		for _, p := range g.Pages {
			out = append(out, &Doc{Guide: g, Page: p, Source: path.Clean(p.Source), URL: g.Slug + "/" + p.Slug + "/"})
		}
	}
	return out
}

// Files lists every file (Markdown and original HTML) that sync copies.
func (m *Manifest) Files() []string {
	var out []string
	for _, d := range m.Docs() {
		out = append(out, d.Source)
	}
	for _, g := range m.Guides {
		out = append(out, g.Originals...)
	}
	return out
}

func (m *Manifest) Guide(slug string) *Guide {
	for _, g := range m.Guides {
		if g.Slug == slug {
			return g
		}
	}
	return nil
}
