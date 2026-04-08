package source

import "io"

// Source is the interface that all sources must implement.
type Source interface {
	Name() string
	Search(query string) ([]*Manga, error)
	ChaptersOf(manga *Manga) ([]*Chapter, error)
	PagesOf(chapter *Chapter) ([]*Page, error)
	ID() string
}

// CloseSource closes the source if it implements io.Closer.
// This is used to clean up resources such as Lua states and headless browsers.
func CloseSource(src Source) {
	if closer, ok := src.(io.Closer); ok {
		_ = closer.Close()
	}
}
