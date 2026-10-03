package sensitive

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/url"
	"strings"

	"github.com/anthonynsimon/bild/transform"
	"github.com/gen2brain/h265/heic"
	"github.com/gen2brain/vpx/webp"
	"github.com/opentreehole/go-common"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"

	"treehole_next/config"
)

const (
	yidunMaxRequestBodyBytes   = 10 * 1024 * 1024
	yidunRequestSafetyMargin   = 1 * 1024 * 1024
	yidunTargetRequestBodySize = yidunMaxRequestBodyBytes - yidunRequestSafetyMargin
	maxDownloadedImageBytes    = 64 * 1024 * 1024
	maxDecodedImagePixels      = 50_000_000
	maxImageDimension          = 4096
	minImageDimension          = 512
	requestFixedOverhead       = 64 * 1024
)

var errImageTooLarge = common.BadRequest("图片过大或格式无法压缩，请压缩后重试")

type yidunImageRequest struct {
	Name string `json:"name"`
	Type int    `json:"type"`
	Data string `json:"data"`
}

func prepareImageForYidun(imageData []byte) (string, error) {
	originalBase64 := base64.StdEncoding.EncodeToString(imageData)
	config, format, err := image.DecodeConfig(bytes.NewReader(imageData))
	// ICO is accepted unchanged, but has no local decoder.
	if bytes.HasPrefix(imageData, []byte{0, 0, 1, 0}) {
		if len(imageData) >= 6 && binary.LittleEndian.Uint16(imageData[4:]) > 0 &&
			6+int(binary.LittleEndian.Uint16(imageData[4:]))*16 <= len(imageData) &&
			yidunRequestBodySize(originalBase64) <= yidunTargetRequestBodySize {
			return originalBase64, nil
		}
		return "", errImageTooLarge
	}
	if err != nil || config.Width <= 0 || config.Height <= 0 || config.Width > maxDecodedImagePixels/config.Height {
		return "", errImageTooLarge
	}
	// TIFF: inspect the first IFD to distinguish DNG and additional pages.
	multiplePages := false
	if format == "tiff" {
		var order binary.ByteOrder = binary.LittleEndian
		if imageData[0] == 'M' {
			order = binary.BigEndian
		}
		offset := uint64(order.Uint32(imageData[4:8]))
		if offset+2 > uint64(len(imageData)) {
			return "", errImageTooLarge
		}
		next := offset + 2 + uint64(order.Uint16(imageData[offset:]))*12
		if next+4 > uint64(len(imageData)) {
			return "", errImageTooLarge
		}
		for entry := offset + 2; entry < next; entry += 12 {
			if order.Uint16(imageData[entry:]) == 50706 { // DNGVersion
				return "", errImageTooLarge
			}
		}
		multiplePages = order.Uint32(imageData[next:]) != 0
	}
	if yidunRequestBodySize(originalBase64) <= yidunTargetRequestBodySize {
		return originalBase64, nil
	}
	var source image.Image
	var animation *webp.WEBP
	switch format {
	// GIF: restore complete canvases before resizing local frames.
	case "gif":
		decoded, err := gif.DecodeAll(bytes.NewReader(imageData))
		if err != nil {
			return "", errImageTooLarge
		}
		animation = &webp.WEBP{Image: gifFrames(decoded), Delay: make([]int, len(decoded.Delay)), LoopCount: webpLoopCount(decoded.LoopCount)}
		for i, delay := range decoded.Delay {
			animation.Delay[i] = delay * 10
		}
	// WebP: DecodeAll composites animation frames and preserves timing.
	case "webp":
		animation, err = webp.DecodeAll(bytes.NewReader(imageData), webp.Options{FrameSizeLimit: maxDecodedImagePixels, Threads: 1})
	// HEIC: apply orientation before converting stills or sequences.
	case "heic":
		decoded, err := heic.DecodeAll(bytes.NewReader(imageData), heic.Options{AutoRotate: true, FrameSizeLimit: maxDecodedImagePixels, Threads: 1})
		if err != nil {
			return "", errImageTooLarge
		}
		animation = &webp.WEBP{Image: decoded.Image, Delay: make([]int, len(decoded.Delay)), LoopCount: webpLoopCount(decoded.LoopCount)}
		for i, delay := range decoded.Delay {
			animation.Delay[i] = int(delay*1000 + 0.5)
			if delay > 0 {
				animation.Delay[i] = max(1, animation.Delay[i])
			}
		}
	// TIFF: do not silently discard additional pages.
	case "tiff":
		if multiplePages {
			return "", errImageTooLarge
		}
		source, err = tiff.Decode(bytes.NewReader(imageData))
	// JPEG / JFIF
	case "jpeg":
		source, err = jpeg.Decode(bytes.NewReader(imageData))
	// PNG
	case "png":
		source, err = png.Decode(bytes.NewReader(imageData))
	// BMP
	case "bmp":
		source, err = bmp.Decode(bytes.NewReader(imageData))
	default:
		return "", errImageTooLarge
	}
	if err != nil {
		return "", errImageTooLarge
	}
	if animation != nil {
		if len(animation.Image) > 1 {
			return compressAnimation(animation)
		}
		source = animation.Image[0]
	}

	opaque := isOpaque(source)
	longestSide := max(source.Bounds().Dx(), source.Bounds().Dy())
	for dimension := min(maxImageDimension, longestSide); ; dimension = max(minImageDimension, int(float64(dimension)*0.8)) {
		candidate := source
		if longestSide > dimension {
			width, height := fitImageSize(source.Bounds(), dimension)
			candidate = transform.Resize(source, width, height, transform.Lanczos)
		}

		qualities := []int{85, 75, 65, 55, 45, 40}
		if !opaque {
			qualities = []int{0}
		}
		for _, quality := range qualities {
			var output bytes.Buffer
			if opaque {
				err = jpeg.Encode(&output, candidate, &jpeg.Options{Quality: quality})
			} else {
				err = png.Encode(&output, candidate)
			}
			if err != nil {
				return "", err
			}
			encodedBase64 := base64.StdEncoding.EncodeToString(output.Bytes())
			if yidunRequestBodySize(encodedBase64) <= yidunTargetRequestBodySize {
				return encodedBase64, nil
			}
		}

		if dimension <= minImageDimension {
			return "", errImageTooLarge
		}
	}
}

func gifFrames(animation *gif.GIF) []image.Image {
	canvas := image.NewRGBA(image.Rect(0, 0, animation.Config.Width, animation.Config.Height))
	background := color.Color(color.Transparent)
	if palette, ok := animation.Config.ColorModel.(color.Palette); ok && int(animation.BackgroundIndex) < len(palette) {
		background = palette[animation.BackgroundIndex]
	}
	for _, entry := range animation.Image[0].Palette {
		if _, _, _, alpha := entry.RGBA(); alpha == 0 {
			background = color.Transparent
			break
		}
	}
	draw.Draw(canvas, canvas.Bounds(), image.NewUniform(background), image.Point{}, draw.Src)
	frames := make([]image.Image, len(animation.Image))
	for i, frame := range animation.Image {
		disposal := animation.Disposal[i]
		var previous *image.RGBA
		if disposal == gif.DisposalPrevious {
			previous = image.NewRGBA(canvas.Bounds())
			copy(previous.Pix, canvas.Pix)
		}
		draw.Draw(canvas, frame.Bounds(), frame, frame.Bounds().Min, draw.Over)
		snapshot := image.NewRGBA(canvas.Bounds())
		copy(snapshot.Pix, canvas.Pix)
		frames[i] = snapshot
		switch disposal {
		case gif.DisposalBackground:
			draw.Draw(canvas, frame.Bounds(), image.NewUniform(background), image.Point{}, draw.Src)
		case gif.DisposalPrevious:
			canvas = previous
		}
	}
	return frames
}

// GIF and HEIC count repeats; WebP counts total plays, with zero meaning forever.
func webpLoopCount(repeats int) int {
	if repeats < 0 {
		return 1
	}
	if repeats > 0 {
		return repeats + 1
	}
	return 0
}

func compressAnimation(animation *webp.WEBP) (string, error) {
	// WebP stores loop counts in 16 bits and frame durations in 24 bits.
	if animation.LoopCount < 0 || animation.LoopCount > 65535 {
		return "", errImageTooLarge
	}
	for _, delay := range animation.Delay {
		if delay < 0 || delay > 0xffffff {
			return "", errImageTooLarge
		}
	}
	longestSide := max(animation.Image[0].Bounds().Dx(), animation.Image[0].Bounds().Dy())
	for dimension := min(1024, longestSide); ; dimension = max(512, dimension/2) {
		candidate := *animation
		candidate.Image = make([]image.Image, len(animation.Image))
		for i, frame := range animation.Image {
			candidate.Image[i] = frame
			if longestSide > dimension {
				width, height := fitImageSize(frame.Bounds(), dimension)
				candidate.Image[i] = transform.Resize(frame, width, height, transform.Lanczos)
			}
		}
		for _, quality := range []int{75, 50} {
			var output bytes.Buffer
			if err := webp.EncodeAll(&output, &candidate, webp.EncodeOptions{Quality: quality, Method: 0, Threads: 1}); err != nil {
				return "", errImageTooLarge
			}
			encoded := base64.StdEncoding.EncodeToString(output.Bytes())
			if yidunRequestBodySize(encoded) <= yidunTargetRequestBodySize {
				return encoded, nil
			}
		}
		if dimension <= 512 {
			return "", errImageTooLarge
		}
	}
}

func isOpaque(source image.Image) bool {
	for y := source.Bounds().Min.Y; y < source.Bounds().Max.Y; y++ {
		for x := source.Bounds().Min.X; x < source.Bounds().Max.X; x++ {
			_, _, _, alpha := source.At(x, y).RGBA()
			if alpha != 0xffff {
				return false
			}
		}
	}
	return true
}

func fitImageSize(bounds image.Rectangle, maxDimension int) (int, int) {
	width, height := bounds.Dx(), bounds.Dy()
	if width >= height {
		return maxDimension, max(1, int(float64(height)*float64(maxDimension)/float64(width)))
	}
	return max(1, int(float64(width)*float64(maxDimension)/float64(height))), maxDimension
}

func yidunRequestBodySize(base64Image string) int {
	images, _ := json.Marshal([]yidunImageRequest{{
		Name: "image_Image",
		Type: 2,
	}})
	form := url.Values{
		"businessId": {config.Config.YiDunBusinessIdImage},
		"images":     {string(images)},
		"nonce":      {"00000000-0000-0000-0000-000000000000"},
		"secretId":   {"secretId"},
		"signature":  {"signature"},
		"timestamp":  {"0000000000000"},
		"version":    {"v5.2"},
	}
	// Base64 needs no JSON escaping; only these characters expand in a form.
	escaped := strings.Count(base64Image, "+") + strings.Count(base64Image, "/") + strings.Count(base64Image, "=")
	return len(form.Encode()) + len(base64Image) + 2*escaped + requestFixedOverhead
}
