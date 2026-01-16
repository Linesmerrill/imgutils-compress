// Package compress provides image compression utilities.
package compress

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"golang.org/x/image/draw"
)

// Quality levels for easy reference.
const (
	QualityLow    = 30
	QualityMedium = 60
	QualityHigh   = 80
	QualityBest   = 95
)

// Options configures compression behavior.
type Options struct {
	Quality    int  // JPEG quality (1-100)
	MaxWidth   int  // Maximum width (0 = no limit)
	MaxHeight  int  // Maximum height (0 = no limit)
	Progressive bool // Not used in standard library
}

// DefaultOptions returns sensible compression defaults.
func DefaultOptions() Options {
	return Options{
		Quality:   80,
		MaxWidth:  0,
		MaxHeight: 0,
	}
}

// Compress compresses an image with the given options.
func Compress(src image.Image, opts Options) image.Image {
	if opts.MaxWidth <= 0 && opts.MaxHeight <= 0 {
		return src
	}

	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Calculate new dimensions if limits are set
	newW, newH := w, h
	if opts.MaxWidth > 0 && w > opts.MaxWidth {
		ratio := float64(opts.MaxWidth) / float64(w)
		newW = opts.MaxWidth
		newH = int(float64(h) * ratio)
	}
	if opts.MaxHeight > 0 && newH > opts.MaxHeight {
		ratio := float64(opts.MaxHeight) / float64(newH)
		newH = opts.MaxHeight
		newW = int(float64(newW) * ratio)
	}

	if newW == w && newH == h {
		return src
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)
	return dst
}

// CompressJPEG compresses and encodes an image as JPEG.
func CompressJPEG(src image.Image, w io.Writer, opts Options) error {
	img := Compress(src, opts)
	quality := opts.Quality
	if quality <= 0 || quality > 100 {
		quality = 80
	}
	return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
}

// CompressPNG compresses and encodes an image as PNG.
// Note: PNG is lossless, so quality setting doesn't apply.
func CompressPNG(src image.Image, w io.Writer, opts Options) error {
	img := Compress(src, opts)
	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	return encoder.Encode(w, img)
}

// CompressFile compresses an image file.
func CompressFile(inputPath, outputPath string, opts Options) error {
	in, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer in.Close()

	src, _, err := image.Decode(in)
	if err != nil {
		return err
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	return CompressJPEG(src, out, opts)
}

// CompressToSize attempts to compress an image to fit within a target file size.
// It iteratively reduces quality until the target is met.
func CompressToSize(src image.Image, targetBytes int) ([]byte, error) {
	quality := 95
	minQuality := 10

	for quality >= minQuality {
		buf := &limitedBuffer{limit: targetBytes * 2}
		err := jpeg.Encode(buf, src, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, err
		}

		if buf.Len() <= targetBytes {
			return buf.Bytes(), nil
		}

		quality -= 5
	}

	// Return best effort at minimum quality
	buf := &limitedBuffer{limit: targetBytes * 2}
	err := jpeg.Encode(buf, src, &jpeg.Options{Quality: minQuality})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// limitedBuffer is a simple buffer for size checking.
type limitedBuffer struct {
	data  []byte
	limit int
}

func (b *limitedBuffer) Write(p []byte) (n int, err error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *limitedBuffer) Bytes() []byte {
	return b.data
}

func (b *limitedBuffer) Len() int {
	return len(b.data)
}
