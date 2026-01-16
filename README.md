# imgutils-compress

[![Go Reference](https://pkg.go.dev/badge/github.com/imgutils-org/imgutils-compress.svg)](https://pkg.go.dev/github.com/imgutils-org/imgutils-compress)
[![Go Report Card](https://goreportcard.com/badge/github.com/imgutils-org/imgutils-compress)](https://goreportcard.com/report/github.com/imgutils-org/imgutils-compress)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go library for compressing images with quality and size controls. Part of the [imgutils](https://github.com/imgutils-org) collection.

## Features

- Quality-based JPEG compression
- PNG compression with maximum compression level
- Maximum dimension constraints with aspect ratio preservation
- Target file size compression
- Predefined quality levels

## Installation

```bash
go get github.com/imgutils-org/imgutils-compress
```

## Quick Start

```go
package main

import (
    "log"

    "github.com/imgutils-org/imgutils-compress"
)

func main() {
    // Compress an image to 80% quality with max dimensions
    err := compress.CompressFile("large.jpg", "small.jpg", compress.Options{
        Quality:   80,
        MaxWidth:  1920,
        MaxHeight: 1080,
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

## Usage Examples

### Basic Compression

```go
// Open image
file, _ := os.Open("input.jpg")
src, _, _ := image.Decode(file)
file.Close()

// Compress to JPEG with 80% quality
out, _ := os.Create("compressed.jpg")
compress.CompressJPEG(src, out, compress.Options{
    Quality: compress.QualityHigh,
})
out.Close()
```

### Quality Levels

```go
// Predefined quality constants
compress.QualityLow    // 30 - Smallest file size
compress.QualityMedium // 60 - Balanced
compress.QualityHigh   // 80 - Good quality (default)
compress.QualityBest   // 95 - Near-lossless

// Use in options
opts := compress.Options{Quality: compress.QualityMedium}
```

### Constrain Dimensions

```go
// Compress and limit dimensions (maintains aspect ratio)
opts := compress.Options{
    Quality:   80,
    MaxWidth:  1920,
    MaxHeight: 1080,
}

compressed := compress.Compress(src, opts)
```

### Compress to Target File Size

```go
// Compress to fit within 500KB
data, err := compress.CompressToSize(src, 500*1024)
if err != nil {
    log.Fatal(err)
}

// Write the bytes to file
os.WriteFile("output.jpg", data, 0644)
```

### PNG Compression

```go
// PNG uses lossless compression (quality setting ignored)
out, _ := os.Create("compressed.png")
compress.CompressPNG(src, out, compress.Options{
    MaxWidth:  1920,
    MaxHeight: 1080,
})
out.Close()
```

### File-to-File Compression

```go
// Simple one-liner
compress.CompressFile("input.jpg", "output.jpg", compress.Options{
    Quality:   75,
    MaxWidth:  2048,
    MaxHeight: 2048,
})
```

## API Reference

### Constants

```go
const (
    QualityLow    = 30  // Smallest file size
    QualityMedium = 60  // Balanced
    QualityHigh   = 80  // Good quality
    QualityBest   = 95  // Near-lossless
)
```

### Types

#### Options

```go
type Options struct {
    Quality   int  // JPEG quality (1-100)
    MaxWidth  int  // Maximum width (0 = no limit)
    MaxHeight int  // Maximum height (0 = no limit)
}
```

### Functions

| Function | Description |
|----------|-------------|
| `Compress(src, opts)` | Resize image if exceeds max dimensions |
| `CompressJPEG(src, w, opts)` | Compress and encode as JPEG |
| `CompressPNG(src, w, opts)` | Compress and encode as PNG |
| `CompressFile(in, out, opts)` | Compress file to file |
| `CompressToSize(src, bytes)` | Compress to target file size |
| `DefaultOptions()` | Returns defaults (quality 80) |

## Quality vs File Size Guide

| Quality | Use Case | Typical Reduction |
|---------|----------|-------------------|
| 30 | Thumbnails, previews | 90%+ |
| 60 | Web images | 70-80% |
| 80 | General photos | 50-60% |
| 95 | Archival, printing | 20-30% |

## Requirements

- Go 1.16 or later

## Related Packages

- [imgutils-convert](https://github.com/imgutils-org/imgutils-convert) - Format conversion
- [imgutils-resize](https://github.com/imgutils-org/imgutils-resize) - Image resizing
- [imgutils-sdk](https://github.com/imgutils-org/imgutils-sdk) - Unified SDK

## License

MIT License - see [LICENSE](LICENSE) for details.
