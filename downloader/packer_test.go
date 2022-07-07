package downloader

import (
	"archive/zip"
	"bytes"
	"github.com/metafates/mangal/common"
	"github.com/metafates/mangal/config"
	"github.com/metafates/mangal/filesystem"
	"github.com/metafates/mangal/scraper"
	"github.com/metafates/mangal/util"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/api"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/afero"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

var testImages []*bytes.Buffer

func init() {
	filesystem.Set(afero.NewMemMapFs())
	config.Initialize("", false)
}

func init() {
	filesystem.Set(afero.NewOsFs())
	defer filesystem.Set(afero.NewMemMapFs())

	const testImagesPath = "../assets/testing_assets"

	// open testing images
	dir, err := afero.ReadDir(filesystem.Get(), testImagesPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range dir {
		if file.IsDir() {
			continue
		}

		contents, err := afero.ReadFile(filesystem.Get(), filepath.Join(testImagesPath, file.Name()))

		if err != nil {
			log.Fatal(err)
		}

		testImages = append(testImages, bytes.NewBuffer(contents))
	}
}

func TestPackToPlain(t *testing.T) {
	Convey("Given "+strconv.Itoa(len(testImages))+" images", t, func() {
		Convey("When packToPlain is called", func() {
			path, err := PackToPlain(testImages, "test", nil)

			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("It should return a path", func() {
				So(path, ShouldNotBeEmpty)

				Convey("That points to a folder", func() {
					isDir, _ := afero.IsDir(filesystem.Get(), path)
					So(isDir, ShouldBeTrue)

					Convey("With the correct name", func() {
						So(path, ShouldEqual, "test")
					})

					Convey("With the correct number of files", func() {
						files, _ := afero.ReadDir(filesystem.Get(), path)
						So(len(files), ShouldEqual, len(testImages))

						Convey("Where each file is a jpg image", func() {
							for _, file := range files {
								isPng := filepath.Ext(file.Name()) == ".jpg"
								So(isPng, ShouldBeTrue)
							}

							Convey("And not empty", func() {
								for _, file := range files {
									isEmpty, _ := afero.IsEmpty(filesystem.Get(), filepath.Join(path, file.Name()))
									So(isEmpty, ShouldBeFalse)
								}
							})
						})
					})

				})
			})

		})
	})
}

func TestPackToZip(t *testing.T) {
	Convey("Given "+strconv.Itoa(len(testImages))+" images", t, func() {
		Convey("When packToZip is called", func() {
			path, err := PackToZip(testImages, "test", nil)

			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("It should return a path", func() {
				So(path, ShouldNotBeEmpty)

				Convey("That points to a zip file", func() {
					isZip := filepath.Ext(path) == ".zip"
					So(isZip, ShouldBeTrue)

					Convey("With the correct name", func() {
						So(filepath.Base(path), ShouldEqual, "test.zip")
					})

					Convey("With the correct number of files", func() {
						fileInfo, err := filesystem.Get().Stat(path)

						if err != nil {
							log.Fatal(err)
						}

						file, err := filesystem.Get().OpenFile(path, os.O_RDONLY, os.FileMode(0644))

						if err != nil {
							log.Fatal(err)
						}

						reader, err := zip.NewReader(file, fileInfo.Size())

						if err != nil {
							log.Fatal(err)
						}

						So(len(reader.File), ShouldEqual, len(testImages))

						Convey("Where each file is a jpg image", func() {
							for _, file := range reader.File {
								isPng := filepath.Ext(file.Name) == ".jpg"
								So(isPng, ShouldBeTrue)
							}

							Convey("And not empty", func() {
								for _, file := range reader.File {
									isEmpty, _ := afero.IsEmpty(filesystem.Get(), filepath.Join(path, file.Name))
									So(isEmpty, ShouldBeFalse)
								}
							})
						})
					})
				})
			})
		})
	})
}

func TestGenerateComicInfo(t *testing.T) {
	Convey("Given sample packer context", t, func() {
		context := &PackerContext{
			Manga: &scraper.URL{
				Address: "https://example.com",
				Info:    "Test manga",
				Index:   0,
			},
			Chapter: &scraper.URL{
				Address: "https://example.com",
				Info:    "Test chapter",
				Index:   0,
			},
		}

		Convey("When generateComicInfo is called", func() {
			xml := generateComicInfo(context)

			Convey("It should return a non-empty string", func() {
				So(xml, ShouldNotBeEmpty)

				Convey("That contains the correct information", func() {
					So(xml, ShouldContainSubstring, "<Series>Test manga</Series>")
					So(xml, ShouldContainSubstring, "<Title>Test chapter</Title>")
				})
			})
		})
	})
}

func TestPackToCBZ(t *testing.T) {
	Convey("Given "+strconv.Itoa(len(testImages))+" images", t, func() {
		Convey("When packToCBZ is called with ComicInfo.xml option", func() {
			config.UserConfig.Formats.Default = common.CBZ
			config.UserConfig.Formats.Comicinfo = true

			path, err := PackToCBZ(testImages, "test", &PackerContext{
				Manga: &scraper.URL{
					Address: "https://example.com",
					Info:    "Example manga",
					Index:   0,
				},
				Chapter: &scraper.URL{
					Address: "https://example.com",
					Info:    "Example chapter",
					Index:   42,
				},
			})

			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("It should return a path", func() {
				So(path, ShouldNotBeEmpty)

				Convey("That points to a cbz file", func() {
					isCbz := filepath.Ext(path) == ".cbz"
					So(isCbz, ShouldBeTrue)

					Convey("With the correct name", func() {
						So(filepath.Base(path), ShouldEqual, "test.cbz")
					})

					Convey("With the correct number of files", func() {
						fileInfo, err := filesystem.Get().Stat(path)

						if err != nil {
							log.Fatal(err)
						}

						file, err := filesystem.Get().OpenFile(path, os.O_RDONLY, os.FileMode(0644))

						if err != nil {
							log.Fatal(err)
						}

						reader, err := zip.NewReader(file, fileInfo.Size())

						if err != nil {
							log.Fatal(err)
						}

						// +1 for the comicinfo.xml file
						So(len(reader.File), ShouldEqual, len(testImages)+1)

						Convey("Where exactly 1 ComicInfo.xml file exists", func() {
							comicinfo, exists := util.Find(reader.File, func(file *zip.File) bool {
								return file.Name == "ComicInfo.xml"
							})

							So(exists, ShouldBeTrue)
							comicinfoCounter := 0
							for _, file := range reader.File {
								if file.Name == "ComicInfo.xml" {
									comicinfoCounter++
								}
							}

							So(comicinfoCounter, ShouldEqual, 1)

							Convey("And it's not empty", func() {
								isEmpty := comicinfo.UncompressedSize64 == 0
								So(isEmpty, ShouldBeFalse)
							})

							Convey("And other files are jpg images", func() {
								for _, file := range reader.File {
									// skip comicinfo.xml
									if file.Name == "ComicInfo.xml" {
										continue
									}

									isPng := filepath.Ext(file.Name) == ".jpg"
									So(isPng, ShouldBeTrue)
								}

								Convey("And they aren't empty", func() {
									for _, file := range reader.File {
										isEmpty, _ := afero.IsEmpty(filesystem.Get(), filepath.Join(path, file.Name))
										So(isEmpty, ShouldBeFalse)
									}
								})
							})
						})
					})
				})
			})
		})
	})
}

func TestPackToPDF(t *testing.T) {
	Convey("Given "+strconv.Itoa(len(testImages))+" images", t, func() {
		Convey("When PackToPDF is called", func() {
			path, err := PackToPDF(testImages, "test", nil)

			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("It should return a path", func() {
				So(path, ShouldNotBeEmpty)

				Convey("That points to a pdf file", func() {
					isPdf := filepath.Ext(path) == ".pdf"
					So(isPdf, ShouldBeTrue)

					Convey("With the correct name", func() {
						So(filepath.Base(path), ShouldEqual, "test.pdf")
					})

					Convey("With the correct number of pages", func() {
						file, err := filesystem.Get().OpenFile(path, os.O_RDONLY, os.FileMode(0644))

						if err != nil {
							log.Fatal(err)
						}

						images, err := pdfcpu.ListImages(file, []string{}, nil)

						if err != nil {
							log.Fatal(err)
						}

						So(len(images), ShouldEqual, len(testImages))
					})
				})
			})
		})
	})
}

func TestPackToEpub(t *testing.T) {
	Convey("Given "+strconv.Itoa(len(testImages))+" images", t, func() {
		Convey("When PackToEpub is called", func() {
			path, err := PackToEpub(testImages, "test", nil)

			Convey("It should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("It should return a path", func() {
				So(path, ShouldNotBeEmpty)

				Convey("That points to an epub file", func() {
					isEpub := filepath.Ext(path) == ".epub"
					So(isEpub, ShouldBeTrue)

					Convey("With the correct name", func() {
						So(filepath.Base(path), ShouldEqual, "test.epub")

						Convey("And it's not empty", func() {
							err = EpubFile.Write(path)
							if err != nil {
								t.Fatal(err)
							}

							isEmpty, _ := afero.IsEmpty(filesystem.Get(), path)
							So(isEmpty, ShouldBeFalse)
						})
					})
				})
			})
		})
	})
}
