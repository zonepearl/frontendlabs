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
		// HomeLayout picks the home page: "classic" (categories, home.html) or
		// "spotlight" (AI spotlight + topic sections, home-spotlight.html).
		HomeLayout string `json:"home_layout"`
	} `json:"site"`
	Categories []Category `json:"categories"` // classic layout: sections and previous/next order
	Guides     []*Guide   `json:"guides"`

	// Spotlight layout only.
	Sections  []Section  `json:"sections"`  // topic sections, each listing its guides in display order
	Path      []string   `json:"path"`      // the course's reading order; previous/next follows it
	Spotlight *Spotlight `json:"spotlight"` // the featured guide shown first

	dir string // directory containing guides.json
}

type Category struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Note  string `json:"note"`
}

// Section is a topic group on the spotlight home page.
type Section struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Note   string   `json:"note"`
	Guides []string `json:"guides"` // slugs, in display (and previous/next) order
}

// Spotlight features one guide at the top of the home page, with deep links
// into it and a list of topics that are planned but not written yet.
type Spotlight struct {
	Guide   string     `json:"guide"`
	Eyebrow string     `json:"eyebrow"`
	Links   []SpotLink `json:"links"`
	Coming  []string   `json:"coming"`
}

type SpotLink struct {
	Title  string `json:"title"`
	Anchor string `json:"anchor"` // heading id inside the spotlight guide
}

type Guide struct {
	Slug      string   `json:"slug"`
	Title     string   `json:"title"`
	Short     string   `json:"short"` // compact name for the course-path ribbon
	Tag       string   `json:"tag"`   // optional note shown on the spotlight layout's card
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
	if err := m.validateLayout(seen); err != nil {
		return nil, fmt.Errorf("%s: %w", file, err)
	}
	return &m, nil
}

// Spotlit reports whether the home page uses the spotlight layout.
func (m *Manifest) Spotlit() bool { return m.Site.HomeLayout == "spotlight" }

// validateLayout makes sure the spotlight layout cannot silently drop a guide
// from the home page or link to one that does not exist.
func (m *Manifest) validateLayout(exists map[string]bool) error {
	switch m.Site.HomeLayout {
	case "", "classic":
		return nil
	case "spotlight":
	default:
		return fmt.Errorf("site.home_layout %q: want \"classic\" or \"spotlight\"", m.Site.HomeLayout)
	}
	placed := map[string]string{}
	for _, sec := range m.Sections {
		for _, slug := range sec.Guides {
			if !exists[slug] {
				return fmt.Errorf("section %q lists unknown guide %q", sec.ID, slug)
			}
			if prev, dup := placed[slug]; dup {
				return fmt.Errorf("guide %q is in both section %q and %q", slug, prev, sec.ID)
			}
			placed[slug] = sec.ID
		}
	}
	for _, g := range m.Guides {
		if placed[g.Slug] == "" {
			return fmt.Errorf("guide %q is in no section, so it would be missing from the home page", g.Slug)
		}
	}
	for _, slug := range m.Path {
		if !exists[slug] {
			return fmt.Errorf("path lists unknown guide %q", slug)
		}
	}
	if m.Spotlight == nil || !exists[m.Spotlight.Guide] {
		return fmt.Errorf("spotlight layout needs spotlight.guide set to an existing guide")
	}
	return nil
}

// SectionOf returns the spotlight section containing the guide, or nil.
func (m *Manifest) SectionOf(slug string) *Section {
	for i := range m.Sections {
		for _, s := range m.Sections[i].Guides {
			if s == slug {
				return &m.Sections[i]
			}
		}
	}
	return nil
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
