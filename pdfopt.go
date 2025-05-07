package pdfopt

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type PDFOpt struct {
	inputFile string
}

func isGhostscriptAvailable() bool {
	_, err := exec.LookPath("gs")
	return err == nil
}

func NewPDFOpt(filename string) (*PDFOpt, error) {
	if !isGhostscriptAvailable() {
		return nil, fmt.Errorf("ghostscript (gs) is not available in system PATH")
	}

	return &PDFOpt{
		inputFile: filename,
	}, nil
}

func (p *PDFOpt) ForEbook(outputFilename string) error {
	return p.executeGhostscript(outputFilename, "/ebook", 120)
}

func (p *PDFOpt) ForEbookInplace() error {
	return p.executeGhostscriptInplace("/ebook", 120)
}

func (p *PDFOpt) ForPrint(outputFilename string) error {
	return p.executeGhostscript(outputFilename, "/printer", 300)
}

func (p *PDFOpt) ForPrintInplace() error {
	return p.executeGhostscriptInplace("/printer", 300)
}

func (p *PDFOpt) ForScreen(outputFilename string) error {
	return p.executeGhostscript(outputFilename, "/screen", 72)
}

func (p *PDFOpt) ForScreenInplace() error {
	return p.executeGhostscriptInplace("/screen", 72)
}

func (p *PDFOpt) executeGhostscript(outputFile, settingsType string, resolution int) error {
	res := strconv.Itoa(resolution)
	cmd := exec.Command("gs",
		"-sDEVICE=pdfwrite",
		// strive for PDF/A-1b compatibility
		"-sProcessColorModel=DeviceCMYK",
		"-sPDFACompatibilityPolicy=1",
		"-sColorConversionStrategy=RGB",
		"-dPDFA=1",

		"-dPDFSETTINGS="+settingsType,
		"-dPDFACompatibilityPolicy=1",
		"-dCompatibilityLevel=1.4",

		// image related options
		"-dDownsampleColorImages=true",
		"-dDownsampleGrayImages=true",
		"-dDownsampleMonoImages=true",
		"-dColorImageResolution="+res,
		"-dGrayImageResolution="+res,
		"-dMonoImageResolution="+res,
		"-dColorImageDownsampleThreshold=1.0",
		"-dGrayImageDownsampleThreshold=1.0",
		"-dMonoImageDownsampleThreshold=1.0",

		// embed fonts far as needed
		"-dSubsetFonts=true",
		"-dEmbedAllFonts=true",

		"-dPreserveAnnots=true",
		"-dPreserveEPSInfo=true",
		"-dPreserveOPIComments=true",
		"-dPreserveOverprintSettings=true",
		"-dAutoFilterColorImages=false",
		"-dAutoFilterGrayImages=false",
		"-dTransparencyLevel=1",

		// silent processing
		"-q",

		// input and output file
		"-o", outputFile,
		p.inputFile,
	)

	return cmd.Run()
}

func (p *PDFOpt) executeGhostscriptInplace(settingsType string, resolution int) error {
	tempFile := p.inputFile + ".tmp"
	err := p.executeGhostscript(tempFile, settingsType, resolution)
	if err != nil {
		return err
	}

	if err := os.Remove(p.inputFile); err != nil {
		return err
	}

	return os.Rename(tempFile, p.inputFile)
}
