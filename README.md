# traceling

[![CI](https://github.com/chinny/traceling/actions/workflows/ci.yml/badge.svg)](https://github.com/chinny/traceling/actions/workflows/ci.yml)

Turn coloring-page PDFs into light-grey tracing pages for kids.

`traceling` rasterizes a PDF and remaps its blacks to a light grey, so kids can
trace the lines with a black marker. White stays white; everything in between
scales linearly.

## Install

Download a prebuilt binary for your platform from the
[latest release](https://github.com/chinny/traceling/releases/latest)
(linux amd64/arm64, macOS amd64/arm64, windows amd64), unpack it, and put
`traceling` on your `PATH`. Binaries are self-contained — MuPDF is statically
linked.

Or build from source (requires Go and a C toolchain):

```sh
make install   # or: go install github.com/chinny/traceling/cmd/traceling@latest
```

## Usage

```sh
traceling coloring-page.pdf                 # writes coloring-page-traceable.pdf
traceling in.pdf -o out.pdf                 # explicit output path
traceling in.pdf --grey 180                 # darker lines (0-254, default 210)
traceling in.pdf --dpi 600                  # higher print resolution (default 300)
traceling in.pdf -f                         # overwrite existing output
```

## How it works

1. Renders each page with MuPDF ([go-fitz](https://github.com/gen2brain/go-fitz)) at the requested DPI.
2. Converts to grayscale and remaps the tonal range: `out = grey + (255-grey) * in / 255`.
3. Writes a fresh PDF (one flate-compressed grayscale image per page) with the
   original page dimensions, so it prints at the original size.

## Development

```sh
make build   # → bin/traceling
make test
make vet
make fmt
```

## License

MIT
