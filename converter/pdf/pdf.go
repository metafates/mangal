package pdf

import (
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/constant"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/source"
	"github.com/metafates/mangal/util"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
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

func save(chapter *source.Chapter, temp bool) (string, error) {
	var (
		mangaDir string
		err      error
	)

	if temp {
		mangaDir, err = filesystem.Get().TempDir("", constant.TempPrefix)
	} else {
		mangaDir, err = prepareMangaDir(chapter.Manga)
	}

	if err != nil {
		return "", err
	}

	chapterPdf := filepath.Join(mangaDir, util.SanitizeFilename(chapter.FormattedName())+".pdf")
	pdfFile, err := filesystem.Get().Create(chapterPdf)
	if err != nil {
		return "", err
	}

	defer func(pdfFile afero.File) {
		_ = pdfFile.Close()
	}(pdfFile)

	var readers = make([]io.Reader, len(chapter.Pages))
	for i, page := range chapter.Pages {
		readers[i] = page
	}

	err = imagesToPDF(pdfFile, readers)
	if err != nil {
		return "", err
	}

	return chapterPdf, nil
}

// prepareMangaDir will create manga direcotry if it doesn't exist
func prepareMangaDir(manga *source.Manga) (mangaDir string, err error) {
	absDownloaderPath, err := filepath.Abs(viper.GetString(config.DownloaderPath))
	if err != nil {
		return "", err
	}

	if viper.GetBool(config.DownloaderCreateMangaDir) {
		mangaDir = filepath.Join(
			absDownloaderPath,
			util.SanitizeFilename(manga.Name),
		)
	} else {
		mangaDir = absDownloaderPath
	}

	if err = filesystem.Get().MkdirAll(mangaDir, os.ModePerm); err != nil {
		return "", err
	}

	return mangaDir, nil
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
			if viper.GetBool(config.FormatsSkipUnsupportedImages) {
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
