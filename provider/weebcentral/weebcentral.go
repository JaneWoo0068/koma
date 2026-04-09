package weebcentral

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/source"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

const (
	ID      = "weebcentral-built-in"
	Name    = "WeebCentral"
	baseURL = "https://weebcentral.com"
)

type Source struct {
	client *http.Client
}

func New() *Source {
	return &Source{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *Source) Name() string { return Name }
func (s *Source) ID() string   { return ID }

func (s *Source) newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", constant.UserAgent)
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Accept-Language", "en-US")
	return req, nil
}

func (s *Source) fetchDocument(req *http.Request) (*goquery.Document, error) {
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weebcentral: %s returned %d", req.URL, resp.StatusCode)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}

func (s *Source) Search(query string) ([]*source.Manga, error) {
	searchURL := fmt.Sprintf("%s/search/data?text=%s&sort=Best+Match&order=Descending&official=Any&display_mode=Full+Display",
		baseURL, url.QueryEscape(query))

	req, err := s.newRequest("GET", searchURL)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", baseURL+"/search")

	doc, err := s.fetchDocument(req)
	if err != nil {
		return nil, err
	}

	var mangas []*source.Manga
	doc.Find("article").Each(func(i int, sel *goquery.Selection) {
		link := sel.Find("a[href*='/series/']").First()
		href, exists := link.Attr("href")
		if !exists {
			return
		}

		name := strings.TrimSpace(sel.Find("a.link.link-hover").First().Text())
		if name == "" {
			// fallback: try tooltip data-tip
			name, _ = sel.Find("span.tooltip").First().Attr("data-tip")
			name = strings.TrimSpace(name)
		}
		if name == "" {
			return
		}

		coverURL := sel.Find("img[alt$='cover']").First().AttrOr("src", "")

		manga := &source.Manga{
			Name:   name,
			URL:    href,
			Index:  uint16(i),
			ID:     filepath.Base(href),
			Source: s,
		}
		manga.Metadata.Cover.ExtraLarge = coverURL

		mangas = append(mangas, manga)
	})

	return mangas, nil
}

func (s *Source) ChaptersOf(manga *source.Manga) ([]*source.Chapter, error) {
	req, err := s.newRequest("GET", manga.URL)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", baseURL)

	doc, err := s.fetchDocument(req)
	if err != nil {
		return nil, err
	}

	var chapters []*source.Chapter
	doc.Find("#chapter-list > div").Each(func(i int, sel *goquery.Selection) {
		link := sel.Find("a[href*='/chapters/']").First()
		href, exists := link.Attr("href")
		if !exists {
			return
		}

		name := strings.TrimSpace(link.Find("span.grow span").First().Text())
		if name == "" {
			name = fmt.Sprintf("Chapter %d", i+1)
		}

		chapter := &source.Chapter{
			Name:  name,
			URL:   href,
			ID:    filepath.Base(href),
			Manga: manga,
		}

		// Parse publish date from time element
		if datetime, exists := sel.Find("time").Attr("datetime"); exists {
			if t, err := time.Parse(time.RFC3339, datetime); err == nil {
				chapter.PublishDate.Year = t.Year()
				chapter.PublishDate.Month = int(t.Month())
				chapter.PublishDate.Day = t.Day()
			}
		}

		chapters = append(chapters, chapter)
	})

	// Chapters are in descending order (newest first), reverse to ascending
	for i, j := 0, len(chapters)-1; i < j; i, j = i+1, j-1 {
		chapters[i], chapters[j] = chapters[j], chapters[i]
	}

	// Set indices after reversing
	for i, ch := range chapters {
		ch.Index = uint16(i)
	}

	manga.Chapters = chapters
	return chapters, nil
}

func (s *Source) PagesOf(chapter *source.Chapter) ([]*source.Page, error) {
	imagesURL := fmt.Sprintf("%s/images?is_prev=False&current_page=1&reading_style=long_strip", chapter.URL)

	req, err := s.newRequest("GET", imagesURL)
	if err != nil {
		return nil, err
	}
	// HTMX headers required by the endpoint
	req.Header.Set("HX-Request", "true")
	req.Header.Set("HX-Current-URL", chapter.URL)
	req.Header.Set("Referer", chapter.URL)

	doc, err := s.fetchDocument(req)
	if err != nil {
		return nil, err
	}

	var pages []*source.Page
	doc.Find("img").Each(func(i int, sel *goquery.Selection) {
		src, exists := sel.Attr("src")
		if !exists || src == "" {
			return
		}
		// Skip broken image placeholder
		if strings.Contains(src, "broken_image") {
			return
		}

		ext := filepath.Ext(src)
		ext = strings.Split(ext, "?")[0]

		page := &source.Page{
			URL:       src,
			Index:     uint16(i),
			Extension: ext,
			Chapter:   chapter,
		}
		pages = append(pages, page)
	})

	chapter.Pages = pages
	return pages, nil
}
