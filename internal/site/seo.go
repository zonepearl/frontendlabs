package site

import (
	"encoding/json"
	"fmt"
	"html/template"
)

// Structured data (schema.org JSON-LD) for search engines: a Course for each
// guide, a TechArticle for each chapter, breadcrumbs everywhere, and the
// WebSite on the home page. json.Marshal escapes <, > and &, so the output is
// safe inside a <script> element.

type ld = map[string]any

func (b *Builder) publisher() ld {
	return ld{"@type": "Organization", "name": "Frontend Labs", "url": b.M.Abs(""),
		"logo": ld{"@type": "ImageObject", "url": b.M.Abs("assets/brand/logo-512.png")}}
}

func (b *Builder) breadcrumbs(items ...[2]string) ld {
	var list []ld
	for i, it := range items {
		list = append(list, ld{"@type": "ListItem", "position": i + 1, "name": it[0], "item": b.M.Abs(it[1])})
	}
	return ld{"@type": "BreadcrumbList", "itemListElement": list}
}

func toJS(v any) template.JS {
	js, _ := json.Marshal(v)
	return template.JS(js)
}

// structuredData describes a guide's landing page (c == nil) or one chapter page.
func (b *Builder) structuredData(pd *prepared, c *ChapterPage) template.JS {
	d, book := pd.doc, pd.book
	desc := book.Subtitle
	if desc == "" {
		desc = d.Guide.Summary
	}
	home := [2]string{"Frontend Labs", ""}
	guide := [2]string{d.Guide.Title, d.Guide.Slug + "/"}
	minutes := fmt.Sprintf("PT%dM", max(1, book.Minutes()))

	if c == nil {
		var main ld
		if d.Page == nil {
			main = ld{
				"@type": "Course", "name": book.Title, "description": desc, "url": b.M.Abs(d.URL),
				"provider": b.publisher(), "inLanguage": "en", "isAccessibleForFree": true,
				"educationalLevel": d.Guide.Level, "timeRequired": minutes,
				"offers":            ld{"@type": "Offer", "category": "Free", "price": 0, "priceCurrency": "USD"},
				"hasCourseInstance": ld{"@type": "CourseInstance", "courseMode": "Online", "courseWorkload": minutes},
			}
		} else {
			main = ld{"@type": "TechArticle", "headline": headline(book.Title), "description": desc, "url": b.M.Abs(d.URL),
				"inLanguage": "en", "author": b.publisher(), "publisher": b.publisher(),
				"image":    b.M.Abs("assets/brand/logo-512.png"),
				"isPartOf": ld{"@type": "Course", "name": d.Guide.Title, "url": b.M.Abs(d.Guide.Slug + "/")}}
		}
		crumbs := []([2]string){home, guide}
		if d.Page != nil {
			crumbs = append(crumbs, [2]string{book.Title, d.URL})
		}
		return toJS(ld{"@context": "https://schema.org", "@graph": []ld{main, b.breadcrumbs(crumbs...)}})
	}

	article := ld{
		"@type": "TechArticle", "headline": headline(c.Title), "url": b.M.Abs(c.URL),
		"mainEntityOfPage": b.M.Abs(c.URL), "inLanguage": "en",
		"author": b.publisher(), "publisher": b.publisher(),
		"image":     b.M.Abs("assets/brand/logo-512.png"),
		"wordCount": c.Words, "timeRequired": fmt.Sprintf("PT%dM", c.Minutes),
		"isPartOf": ld{"@type": "Course", "name": book.Title, "url": b.M.Abs(d.URL)},
	}
	if c.Desc != "" {
		article["description"] = c.Desc
	}
	if c.Part != "" {
		article["articleSection"] = c.Part
	}
	crumbs := []([2]string){home, guide}
	if d.Page != nil {
		crumbs = append(crumbs, [2]string{book.Title, d.URL})
	}
	crumbs = append(crumbs, [2]string{c.Title, c.URL})
	return toJS(ld{"@context": "https://schema.org", "@graph": []ld{article, b.breadcrumbs(crumbs...)}})
}

// homeStructuredData describes the site itself.
func (b *Builder) homeStructuredData() template.JS {
	return toJS(ld{"@context": "https://schema.org", "@graph": []ld{
		{"@type": "WebSite", "name": b.M.Site.Title, "url": b.M.Abs(""), "description": b.M.Site.Tagline,
			"inLanguage": "en", "publisher": b.publisher()},
		b.publisher(),
	}})
}

// headline keeps schema.org headlines within Google's 110-character limit.
func headline(s string) string {
	if len(s) <= 110 {
		return s
	}
	return trimWords(s, 109)
}
