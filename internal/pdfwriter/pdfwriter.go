// Package pdfwriter emits a minimal PDF: one full-bleed 8-bit grayscale
// image per page, flate-compressed. No external dependencies.
package pdfwriter

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"image"
	"io"
)

type page struct {
	img      *image.Gray
	widthPt  float64
	heightPt float64
}

// Writer accumulates pages and serializes them as a PDF document.
type Writer struct {
	pages []page
}

// AddPage appends img as a new page of the given size in PDF points
// (1 pt = 1/72 inch). The image is stretched to fill the page.
func (w *Writer) AddPage(img *image.Gray, widthPt, heightPt float64) {
	w.pages = append(w.pages, page{img, widthPt, heightPt})
}

// objects per page: page dict, content stream, image XObject
const objsPerPage = 3

// WriteTo serializes the document. It implements io.WriterTo.
func (w *Writer) WriteTo(out io.Writer) (int64, error) {
	if len(w.pages) == 0 {
		return 0, fmt.Errorf("pdfwriter: no pages added")
	}

	var buf bytes.Buffer
	offsets := make([]int, 2+len(w.pages)*objsPerPage+1) // 1-indexed

	buf.WriteString("%PDF-1.4\n%\xff\xff\xff\xff\n")

	offsets[1] = buf.Len()
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")

	offsets[2] = buf.Len()
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [")
	for i := range w.pages {
		fmt.Fprintf(&buf, " %d 0 R", pageObj(i))
	}
	fmt.Fprintf(&buf, " ] /Count %d >>\nendobj\n", len(w.pages))

	for i, p := range w.pages {
		content := fmt.Sprintf("q %.4f 0 0 %.4f 0 0 cm /Im0 Do Q", p.widthPt, p.heightPt)

		offsets[pageObj(i)] = buf.Len()
		fmt.Fprintf(&buf,
			"%d 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %.4f %.4f] "+
				"/Contents %d 0 R /Resources << /XObject << /Im0 %d 0 R >> >> >>\nendobj\n",
			pageObj(i), p.widthPt, p.heightPt, contentObj(i), imageObj(i))

		offsets[contentObj(i)] = buf.Len()
		fmt.Fprintf(&buf, "%d 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n",
			contentObj(i), len(content), content)

		samples, err := compressSamples(p.img)
		if err != nil {
			return 0, err
		}
		offsets[imageObj(i)] = buf.Len()
		fmt.Fprintf(&buf,
			"%d 0 obj\n<< /Type /XObject /Subtype /Image /Width %d /Height %d "+
				"/ColorSpace /DeviceGray /BitsPerComponent 8 /Filter /FlateDecode /Length %d >>\nstream\n",
			imageObj(i), p.img.Bounds().Dx(), p.img.Bounds().Dy(), len(samples))
		buf.Write(samples)
		buf.WriteString("\nendstream\nendobj\n")
	}

	xrefStart := buf.Len()
	numObjs := len(offsets) - 1
	fmt.Fprintf(&buf, "xref\n0 %d\n0000000000 65535 f \n", numObjs+1)
	for i := 1; i <= numObjs; i++ {
		fmt.Fprintf(&buf, "%010d 00000 n \n", offsets[i])
	}
	fmt.Fprintf(&buf, "trailer\n<< /Size %d /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n",
		numObjs+1, xrefStart)

	return buf.WriteTo(out)
}

func pageObj(i int) int    { return 3 + i*objsPerPage }
func contentObj(i int) int { return 4 + i*objsPerPage }
func imageObj(i int) int   { return 5 + i*objsPerPage }

// compressSamples flate-compresses the image's raw samples, packing rows
// tightly in case the stride exceeds the row width.
func compressSamples(img *image.Gray) ([]byte, error) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	var out bytes.Buffer
	zw := zlib.NewWriter(&out)
	for y := 0; y < h; y++ {
		if _, err := zw.Write(img.Pix[y*img.Stride : y*img.Stride+w]); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
