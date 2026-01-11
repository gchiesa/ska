package tui

import (
	"bytes"
	_ "embed"
	"os"
	"strings"

	"github.com/BourgeoisBear/rasterm"
	"github.com/common-nighthawk/go-figure"
	"github.com/gchiesa/ska/internal/configuration"
)

//go:embed resources/logo.png
var logoData []byte

// TextBanner generate a text banner
func TextBanner() string {
	var builder strings.Builder

	figure.Write(&builder, figure.NewFigure(configuration.AppIdentifier, "doom", true))
	builder.WriteString("Your scaffolding buddy!\n")
	builder.WriteRune('\n')

	return builder.String()
}

// CanUseGraphic return true if the terminal is capable of graphical images
func CanUseGraphic() bool {
	return rasterm.IsKittyCapable() || rasterm.IsItermCapable()
}

// GraphicalBanner generate a banner from an image embedded in the app
func GraphicalBanner() error {
	var err error
	image := bytes.NewReader(logoData)
	if rasterm.IsKittyCapable() {
		opts := rasterm.KittyImgOpts{}
		err = rasterm.KittyCopyPNGInline(os.Stdout, image, opts)
	}
	if rasterm.IsItermCapable() {
		opts := rasterm.ItermImgOpts{}
		err = rasterm.ItermCopyFileInlineWithOptions(os.Stdout, image, opts)
	}
	if err != nil {
		return err
	}
	var builder strings.Builder
	builder.WriteRune('\n')
	println(builder.String())
	return nil
}
