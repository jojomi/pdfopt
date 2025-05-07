package pdfopt

import (
	"os"
	"os/exec"
	"strconv"
)

type PDFOpt struct {
	inputFile    string
	settingsType string
	resolution   int
}

func isGhostscriptAvailable() bool {
	_, err := exec.LookPath("gs")
	return err == nil
}

func NewPDFOpt(filename string) *PDFOpt {
	if !isGhostscriptAvailable() {
		panic("ghostscript (gs) is not available in system PATH")
	}

	return &PDFOpt{
		inputFile:    filename,
		settingsType: "/screen",
		resolution:   72,
	}
}

func (p *PDFOpt) ForEbook() *PDFOpt {
	p.settingsType = "/ebook"
	p.resolution = 150
	return p
}

func (p *PDFOpt) ForPrepress() *PDFOpt {
	p.settingsType = "/prepress"
	p.resolution = 300
	return p
}

func (p *PDFOpt) ForPrint() *PDFOpt {
	p.settingsType = "/printer"
	p.resolution = 300
	return p
}

func (p *PDFOpt) ForScreen() *PDFOpt {
	p.settingsType = "/screen"
	p.resolution = 72
	return p
}

func (p *PDFOpt) ImageDPI(dpi int) *PDFOpt {
	p.resolution = dpi
	return p
}

func (p *PDFOpt) Optimize(outputFilename string) error {
	return p.executeGhostscript(outputFilename, p.settingsType, p.resolution)
}

func (p *PDFOpt) OptimizeInplace() error {
	return p.executeGhostscriptInplace(p.settingsType, p.resolution)
}

func (p *PDFOpt) executeGhostscript(outputFile, settingsType string, resolution int) error {
	res := strconv.Itoa(resolution)
	cmd := exec.Command("gs",
		"-sDEVICE=pdfwrite",
		"-dPDFSETTINGS="+settingsType,

		// strive for PDF/A-1b compatibility
		"-sProcessColorModel=DeviceCMYK",
		/*
			"-dPDFA=1", // kills fonts
			//"-sColorConversionStrategy=RGB",
			"-dPDFACompatibilityPolicy=1",
		*/
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

		// embed fonts
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
