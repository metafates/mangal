 # Note

This package is an improved version of [x/image/tiff](https://github.com/golang/image/tree/master/tiff) featuring:

* Read support for CCITT Group3/4 compressed images using [x/image/ccitt](https://github.com/golang/image/tree/master/ccitt)
* Read/write support for LZW compressed images using [github.com/hhrutter/lzw](https://github.com/hhrutter/lzw)
* Read/write support for the CMYK color model.


## Background

Working on [pdfcpu](https://github.com/pdfcpu/pdfcpu) (a PDF processor) created a need for processing TIFF files and LZW compression in details beyond the standard library.

1) CCITT compression for monochrome images was the first need. This is being addressed as part of ongoing work on [x/image/ccitt](https://github.com/golang/image/tree/master/ccitt).

2) As stated in this [golang proposal](https://github.com/golang/go/issues/25409) Go LZW implementations are spread out over the standard library at [compress/lzw](https://github.com/golang/go/tree/master/src/compress/lzw) and [x/image/tiff/lzw](https://github.com/golang/image/tree/master/tiff/lzw). As of Go 1.12 [compress/lzw](https://github.com/golang/go/tree/master/src/compress/lzw) works reliably for GIF only. This is also the reason the TIFF package at [x/image/tiff](https://github.com/golang/image/tree/master/tiff) provides its own lzw implementation for compression. With PDF there is a third variant of lzw needed for reading/writing lzw compressed PDF object streams and processing embedded TIFF images.
[github.com/hhrutter/lzw](https://github.com/hhrutter/lzw) fills this gap. It is an extended version of [compress/lzw](https://github.com/golang/go/tree/master/src/compress/lzw) supporting GIF, PDF and TIFF.

3) The PDF specification defines a CMYK color space. This is currently not supported at [x/image/tiff](https://github.com/golang/image/tree/master/tiff).

## Goal

An improved version of [x/image/tiff](https://github.com/golang/image/tree/master/tiff) with full read/write support for CCITT, LZW compression and the CMYK color model.
