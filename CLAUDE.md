# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```sh
make build       # build → ./bin/traceling (injects version/commit/date via ldflags)
make test        # go test ./...
make vet         # go vet ./...
make fmt         # gofmt -s -w .
make install     # go install to $GOBIN

# Run a single package's tests
go test ./internal/lighten/...
go test ./internal/pdfwriter/...
```

## Architecture

The CLI (`cmd/traceling/`) is a thin Cobra shell. All logic lives in `internal/`:

```
internal/
  render/     — PDF → *image.RGBA pages via MuPDF (go-fitz; default cgo build
                statically links the bundled MuPDF — with CGO_ENABLED=0 it
                falls back to purego and dlopens a system libmupdf at runtime);
                also reports page size in PDF points via Bound()
  lighten/    — tonal remap: grayscale conversion + LUT so black → grey level,
                white → white, linear in between
  pdfwriter/  — minimal dependency-free PDF writer: one full-bleed 8-bit
                DeviceGray FlateDecode image per page, hand-built xref table
```

**Data flow:** `render.Open` → per page `render.Document.Page(i, dpi)` →
`lighten.ToTraceable(img, grey)` → `pdfwriter.Writer.AddPage(gray, wPt, hPt)` →
`WriteTo(file)`.

**Gotchas:**
- go-fitz v1.28's plain `Image()` misrenders (tiny bounds); always use `ImageDPI`.
- `pdfwriter` packs rows tightly because `image.Gray.Stride` can exceed width.
- The round-trip test re-renders written PDFs through go-fitz to catch
  structural PDF bugs.

**Module path:** `github.com/chinny/traceling`

## Releases

Push a `v*` tag → `.github/workflows/release.yml` builds native binaries on a
runner matrix (linux amd64/arm64, darwin amd64/arm64, windows amd64) and
publishes archives + checksums as a GitHub Release. Builds are native per
platform because the cgo MuPDF link can't cross-compile without a matching
C toolchain.
