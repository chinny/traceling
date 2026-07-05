package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/chinny/traceling/internal/lighten"
	"github.com/chinny/traceling/internal/pdfwriter"
	"github.com/chinny/traceling/internal/render"
)

var (
	flagOutput string
	flagGrey   int
	flagDPI    int
	flagForce  bool
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "traceling INPUT.pdf",
		Short: "Turn coloring-page PDFs into light-grey tracing pages",
		Long: `traceling rasterizes a PDF and remaps its blacks to a light grey,
producing pages kids can trace over with a marker. White stays white;
everything in between scales linearly.`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args[0])
		},
	}

	root.Flags().StringVarP(&flagOutput, "output", "o", "", "Output PDF path (default: INPUT-traceable.pdf)")
	root.Flags().IntVarP(&flagGrey, "grey", "g", 210, "Grey level for black lines, 0-254 (lower = darker)")
	root.Flags().IntVarP(&flagDPI, "dpi", "d", 300, "Render resolution")
	root.Flags().BoolVarP(&flagForce, "force", "f", false, "Overwrite the output file if it exists")

	root.AddCommand(newVersionCmd())
	return root
}

func run(input string) error {
	if flagGrey < 0 || flagGrey > 254 {
		return fmt.Errorf("--grey must be 0-254, got %d", flagGrey)
	}
	if flagDPI < 72 || flagDPI > 1200 {
		return fmt.Errorf("--dpi must be 72-1200, got %d", flagDPI)
	}

	output := flagOutput
	if output == "" {
		ext := filepath.Ext(input)
		output = strings.TrimSuffix(input, ext) + "-traceable.pdf"
	}
	if !flagForce {
		if _, err := os.Stat(output); err == nil {
			return fmt.Errorf("%s already exists (use --force to overwrite)", output)
		}
	}

	doc, err := render.Open(input)
	if err != nil {
		return err
	}
	defer doc.Close()

	var pdf pdfwriter.Writer
	n := doc.NumPages()
	for i := 0; i < n; i++ {
		page, err := doc.Page(i, flagDPI)
		if err != nil {
			return err
		}
		grey := lighten.ToTraceable(page.Image, uint8(flagGrey))
		pdf.AddPage(grey, page.WidthPt, page.HeightPt)
		fmt.Printf("page %d/%d done\n", i+1, n)
	}

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	if _, err := pdf.WriteTo(f); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	fmt.Printf("wrote %s (%d pages, grey %d, %d dpi)\n", output, n, flagGrey, flagDPI)
	return nil
}
