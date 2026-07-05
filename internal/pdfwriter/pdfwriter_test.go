package pdfwriter

import (
	"bytes"
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/gen2brain/go-fitz"
)

func TestWriteToRejectsEmptyDocument(t *testing.T) {
	var w Writer
	if _, err := w.WriteTo(&bytes.Buffer{}); err == nil {
		t.Fatal("expected error for empty document")
	}
}

// TestRoundTrip writes a two-page PDF and re-renders it with MuPDF to
// verify page count, page size, and pixel values survive.
func TestRoundTrip(t *testing.T) {
	img := image.NewGray(image.Rect(0, 0, 4, 4))
	for i := range img.Pix {
		img.Pix[i] = 255 // white
	}
	img.Pix[0] = 210 // one light-grey pixel top-left

	var w Writer
	w.AddPage(img, 612, 792) // US Letter
	w.AddPage(img, 612, 792)

	path := filepath.Join(t.TempDir(), "out.pdf")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.WriteTo(f); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	doc, err := fitz.New(path)
	if err != nil {
		t.Fatalf("reopen written PDF: %v", err)
	}
	defer doc.Close()

	if got := doc.NumPage(); got != 2 {
		t.Fatalf("NumPage = %d, want 2", got)
	}
	bounds, err := doc.Bound(0)
	if err != nil {
		t.Fatal(err)
	}
	if bounds.Dx() != 612 || bounds.Dy() != 792 {
		t.Errorf("page size = %dx%d pt, want 612x792", bounds.Dx(), bounds.Dy())
	}

	// Render at 72 DPI and sample the corners of the stretched image.
	rendered, err := doc.ImageDPI(0, 72)
	if err != nil {
		t.Fatal(err)
	}
	rb := rendered.Bounds()
	tl := rendered.RGBAAt(rb.Min.X+2, rb.Min.Y+2)
	br := rendered.RGBAAt(rb.Max.X-3, rb.Max.Y-3)
	if tl.R < 190 || tl.R > 230 {
		t.Errorf("top-left pixel R = %d, want ~210 (light grey)", tl.R)
	}
	if br.R < 250 {
		t.Errorf("bottom-right pixel R = %d, want ~255 (white)", br.R)
	}
}
