package mangadex

import (
	"fmt"
	"github.com/darylhjd/mangodex"
	"github.com/metafates/mangal/key"
	"github.com/metafates/mangal/source"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"net/url"
	"strconv"
	"time"
)

func (m *Mangadex) ChaptersOf(manga *source.Manga) ([]*source.Chapter, error) {
	params := url.Values{}
	params.Set("limit", strconv.Itoa(500))
	ratings := []string{mangodex.Safe, mangodex.Suggestive}
	for _, rating := range ratings {
		params.Add("contentRating[]", rating)
	}

	if viper.GetBool(key.MangadexNSFW) {
		params.Add("contentRating[]", mangodex.Porn)
		params.Add("contentRating[]", mangodex.Erotica)
	}

	// scanlation group for the chapter
	params.Add("includes[]", mangodex.ScanlationGroupRel)
	params.Set("order[chapter]", "asc")

	type chapterCandidate struct {
		chapter mangodex.Chapter
		volume  string
	}

	var (
		candidates []chapterCandidate
		currOffset = 0
	)

	language := viper.GetString(key.MangadexLanguage)

	firstRequest := true
	for {
		if !firstRequest {
			if delay := viper.GetInt(key.MangadexRequestDelay); delay > 0 {
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
		}
		firstRequest = false

		params.Set("offset", strconv.Itoa(currOffset))
		list, err := m.client.Chapter.GetMangaChapters(manga.ID, params)
		if err != nil {
			return nil, err
		}

		for _, chapter := range list.Data {
			// Skip external chapters. Their pages cannot be downloaded.
			if chapter.Attributes.ExternalURL != nil && !viper.GetBool(key.MangadexShowUnavailableChapters) {
				continue
			}

			// skip chapters that are not in the current language
			if language != "any" && chapter.Attributes.TranslatedLanguage != language {
				continue
			}

			var volume string
			if chapter.Attributes.Volume != nil {
				volume = fmt.Sprintf("Vol.%s", *chapter.Attributes.Volume)
			}

			candidates = append(candidates, chapterCandidate{chapter: chapter, volume: volume})
		}
		currOffset += 500
		if currOffset >= list.Total {
			break
		}
	}

	// Deduplicate chapters from different scanlation groups (#162).
	// When multiple groups translate the same chapter, prefer the official group.
	seen := make(map[string]int) // dedupKey -> index in deduplicated slice
	var deduplicated []chapterCandidate

	for _, c := range candidates {
		dedupKey := c.volume + "|" + c.chapter.GetChapterNum()

		if idx, exists := seen[dedupKey]; exists {
			// Replace the existing entry only if this one is from an official group
			if isOfficialChapter(c.chapter) && !isOfficialChapter(deduplicated[idx].chapter) {
				deduplicated[idx] = c
			}
		} else {
			seen[dedupKey] = len(deduplicated)
			deduplicated = append(deduplicated, c)
		}
	}

	// Build the final chapter list from deduplicated candidates.
	chapters := make([]*source.Chapter, 0, len(deduplicated))
	for i, c := range deduplicated {
		name := c.chapter.GetTitle()
		if name == "" {
			name = fmt.Sprintf("Chapter %s", c.chapter.GetChapterNum())
		} else {
			name = fmt.Sprintf("Chapter %s - %s", c.chapter.GetChapterNum(), name)
		}

		chapter := &source.Chapter{
			Name:   name,
			Index:  uint16(i),
			ID:     c.chapter.ID,
			URL:    fmt.Sprintf("https://mangadex.org/chapter/%s", c.chapter.ID),
			Manga:  manga,
			Volume: c.volume,
		}

		// Parse the chapter publish date from MangaDex if available.
		if t, err := time.Parse(time.RFC3339, c.chapter.Attributes.PublishAt); err == nil {
			chapter.PublishDate = source.NewDate(t.Year(), int(t.Month()), t.Day())
		}

		chapters = append(chapters, chapter)
	}

	slices.SortFunc(chapters, func(a, b *source.Chapter) bool {
		return a.Index < b.Index
	})

	manga.Chapters = chapters
	return chapters, nil
}

// isOfficialChapter returns true if any scanlation group relationship
// on the chapter has Official set to true.
func isOfficialChapter(chapter mangodex.Chapter) bool {
	for _, rel := range chapter.Relationships {
		if rel.Type != mangodex.ScanlationGroupRel {
			continue
		}
		if attrs, ok := rel.Attributes.(*mangodex.ScanlationGroupAttributes); ok && attrs != nil {
			if attrs.Official {
				return true
			}
		}
	}
	return false
}
