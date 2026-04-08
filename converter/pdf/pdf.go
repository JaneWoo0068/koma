package pdf

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/key"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/spf13/viper"
	"io"

	_ "golang.org/x/image/webp"
)

type PDF struct{}

func New() *PDF {
	return &PDF{}
}

func (*PDF) Save(chapter *source.Chapter) (string, error) {
	return save(chapter, false)
}

func (*PDF) SaveTemp(chapter *source.Chapter) (string, error) {
	return save(chapter, true)
}

func save(chapter *source.Chapter, temp bool) (path string, err error) {
	path, err = chapter.Path(temp)
	if err != nil {
		return
	}

	file, err := filesystem.Api().Create(path)
	if err != nil {
		return
	}

	defer util.Ignore(file.Close)

	err = pagesToPDF(file, chapter.Pages)
	return
}

// pagesToPDF will convert images to PDF and write to w
func pagesToPDF(w io.Writer, pages []*source.Page) error {
	conf := pdfcpu.NewDefaultConfiguration()
	conf.Cmd = pdfcpu.IMPORTIMAGES
	imp := pdfcpu.DefaultImportConfig()

	var (
		ctx *pdfcpu.Context
		err error
	)

	ctx, err = pdfcpu.CreateContextWithXRefTable(conf, imp.PageDim)
	if err != nil {
		return err
	}

	pagesIndRef, err := ctx.Pages()
	if err != nil {
		return err
	}

	// This is the page tree root.
	pagesDict, err := ctx.DereferenceDict(*pagesIndRef)
	if err != nil {
		return err
	}

	for _, r := range pages {
		// Read the page contents so we can decode image dimensions
		// and then create a fresh reader for the PDF import.
		data := r.Contents.Bytes()

		// Decode image dimensions to set the page size to match the image,
		// preventing tall webtoon/manhwa pages from being clipped to A4.
		imgCfg, _, err := image.DecodeConfig(bytes.NewReader(data))
		if err == nil {
			imp.PageDim = &pdfcpu.Dim{
				Width:  float64(imgCfg.Width),
				Height: float64(imgCfg.Height),
			}
		}

		indRef, err := pdfcpu.NewPageForImage(ctx.XRefTable, bytes.NewReader(data), pagesIndRef, imp)

		if err != nil {
			if viper.GetBool(key.FormatsSkipUnsupportedImages) {
				continue
			}

			return err
		}

		if err = pdfcpu.AppendPageTree(indRef, 1, pagesDict); err != nil {
			return err
		}

		ctx.PageCount++
	}

	if err = api.WriteContext(ctx, w); err != nil {
		return err
	}

	log.Stats.Printf("XRefTable:\n%s\n", ctx)

	return nil
}
