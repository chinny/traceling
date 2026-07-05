// Package render rasterizes PDF pages via MuPDF (go-fitz).
package render

import (
	"fmt"
	"image"

	"github.com/gen2brain/go-fitz"
)

// Page is a rasterized PDF page plus its size in PDF points.
type Page struct {
	Image    *image.RGBA
	WidthPt  float64
	HeightPt float64
}

// Document wraps an open PDF.
type Document struct {
	doc *fitz.Document
}

func Open(path string) (*Document, error) {
	doc, err := fitz.New(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	return &Document{doc: doc}, nil
}

func (d *Document) NumPages() int { return d.doc.NumPage() }

// Page rasterizes page i (0-based) at the given DPI.
func (d *Document) Page(i, dpi int) (*Page, error) {
	bounds, err := d.doc.Bound(i)
	if err != nil {
		return nil, fmt.Errorf("page %d bounds: %w", i+1, err)
	}
	img, err := d.doc.ImageDPI(i, float64(dpi))
	if err != nil {
		return nil, fmt.Errorf("render page %d: %w", i+1, err)
	}
	// Bound reports the page box at 72 DPI, i.e. in points.
	return &Page{
		Image:    img,
		WidthPt:  float64(bounds.Dx()),
		HeightPt: float64(bounds.Dy()),
	}, nil
}

func (d *Document) Close() error { return d.doc.Close() }
