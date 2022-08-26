package pdf

import (
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/spf13/viper"
	"io"
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

	pdfFile, err := filesystem.Get().Create(path)
	if err != nil {
		return
	}

	defer util.Ignore(pdfFile.Close)

	var readers = make([]io.Reader, len(chapter.Pages))
	for i, page := range chapter.Pages {
		readers[i] = page
	}

	err = imagesToPDF(pdfFile, readers)
	return
}

// imagesToPDF will convert images to PDF and write to w
func imagesToPDF(w io.Writer, images []io.Reader) error {
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

	for _, r := range images {
		indRef, err := pdfcpu.NewPageForImage(ctx.XRefTable, r, pagesIndRef, imp)

		if err != nil {
			if viper.GetBool(constant.FormatsSkipUnsupportedImages) {
				continue
			}

			return err
		}

		if err = pdfcpu.AppendPageTree(indRef, 1, pagesDict); err != nil {
			return err
		}

		ctx.PageCount++
	}

	if conf.ValidationMode != pdfcpu.ValidationNone {
		if err = api.ValidateContext(ctx); err != nil {
			return err
		}
	}

	if err = api.WriteContext(ctx, w); err != nil {
		return err
	}

	log.Stats.Printf("XRefTable:\n%s\n", ctx)

	return nil
}
