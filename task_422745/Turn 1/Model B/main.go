package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/png"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"golang.org/x/image/transform"
)

func filenames(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var filenames []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filenames = append(filenames, filepath.Join(dir, file.Name()))
	}
	return filenames, nil
}

// LoadImage reads an image from a file.
func LoadImage(fn string) (image.Image, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var m image.Image
	switch ext := filepath.Ext(fn); ext {
	case ".jpg", ".jpeg":
		m, _, err = jpeg.Decode(f)
	case ".png":
		m, err = png.Decode(f)
	case ".gif":
		m, err = gif.Decode(f)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

// processRawImages reads image files from the input directory and processes them through the pipeline, 
// saving the resulting images to the output directory.
func processRawImages(inputDir string, outputDir string, closeCh chan struct{}, wrGlob chan []byte) {
	defer func() {
		close(wrGlob)
		select {
		case <-closeCh:
			return
		default:
		}
	}()

	imgFilenames, err := filenames(inputDir)
	if err != nil {
		log.Fatalf("failed to read filenames: %v", err)
	}
	log.Printf("found %d image files", len(imgFilenames))

	for _, filepath := range imgFilenames {
		select {
		case <-closeCh:
			log.Println("input processor: shutting down")
			return
		default:
		}
		img, err := LoadImage(filepath)
		if err != nil {
			log.Printf("error loading image '%s': %v", filepath, err)
			continue
		}
		wrGlob <- JPEGEncode(FilterBlur(img))
	}
}

func pipeFilter(in <-chan image.Image, out chan<- image.Image, closeCh chan struct{}) {
	defer func() {
		close(out)
		select {
		case <-closeCh:
			return
		default:
		}
	}()

	for {
		select {
		case <-closeCh:
			return
		case im, ok := <-in:
			if !ok {
				return
			}
			filtered := FilterBlur(im)
			out <- filtered
		}
	}
}

func pipeResize(in <-chan image.Image, out chan<- []byte, targetSize int, closeCh chan struct{}) {
	defer func() {
		close(out)
		select {
		case <-closeCh:
			return
		default:
		}
	}()

	for {
		select {
		case <-closeCh:
			return
		case im, ok := <-in:
			if !ok {
				return
			}
			resized := Resize(im, targetSize, targetSize)
			pngData := PNGEncode(resized)
			out <- pngData
		}
	}
}

func FilterBlur(im image.Image) image.Image {
	bounds := im.Bounds()
	xn := bounds.Max.X / 3
	yn := bounds.Max.Y / 3
	factors := [][]int{{0, 0}, {1, 0}, {2, 0}, {0, 1}, {1, 1}, {2, 1}, {0, 2}, {1, 2}, {2, 2}}
	blurred := image.NewNRGBA(bounds)

	for y := 0; y < yn; y++ {
		for x := 0; x < xn; x++ {
			cr := 0.0
			cg := 0.0
			cb := 0.0
			for _, factor := range factors {
				px, py := x*3+factor[0], y*3+factor[1]
				r, g, b, a := im.At(px, py).RGBA()
				cr += float64(r)
				cg += float64(g)
				cb += float64(b)
			}
			n := 1.0 / 9.0
			cr = math.Min(255.0, math.Max(0.0, cr*n))
			cg = math.Min(255.0, math.Max(0.0, cg*n))
			cb = math.Min(255.0, math.Max(0.0, cb*n))
			blurred.Set(x*3, y*3, color.NRGBA{uint8(cr), uint8(cg), uint8(cb), 255})
			blurred.Set(x*3+1, y*3, color.NRGBA{uint8(cr), uint8(cg), uint8(cb), 255})
			blurred.Set(x*3+2, y*3, color.NRGBA{uint8(cr), uint8(cg), uint8(cb), 255})
			blurred.Set(x*3, y*3+1, color.NRGBA{uint8(cr), uint8(cg), uint8(cb), 255})
			blurred.Set(x*3+1, y*3+1, color.NRGBA{uint8(cr), uint8(cg), uint8(cb), 255})
			blurred.Set(x*3+2, y*3+1, color.NRGBA{uint8(cr), uint8(cg), uint8(cb), 255})
			blurred.Set(x*3, y*3+2, color.NRGBA{uint8(cr), uint8(cg), uint8(cb), 255})
			blurred.Set(x*3+1, y*3+2, color.NRGBA{uint8(cr), uint8(cg), uint8(cb), 255})
			blurred.Set(x*3+2, y*3+2, color.NRGBA{uint8(cr), uint8(cg), uint8(cb), 255})
		}
	}
	return blurred
}

// Resize resizes an image to the specified target width and height using Bilinear interpolation.
func Resize(im image.Image, targetWidth, targetHeight int) image.Image {
	b := im.Bounds()
	w := b.Max.X
	h := b.Max.Y

	if w == targetWidth && h == targetHeight {
		return im
	}

	rm := float64(targetWidth) / float64(w)
	rh := float64(targetHeight) / float64(h)

	out := image.NewNRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.Bilinear(out, out.Bounds(), im, b, draw.Over)

	return out
}

func PNGEncode(im image.Image) []byte {
	var buf bytes.Buffer
	err := png.Encode(&buf, im)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func JPEGEncode(im image.Image) []byte {
	var buf bytes.Buffer
	img := FixedImage{im}
	err := jpeg.Encode(&buf, &img, &jpeg.Options{Quality: 50})
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
type FixedImage struct {
	image.Image
}

func (img *FixedImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, 1000, 1000)
}

func main() {
	// Input and output directories
	inputDir := "images"
	outputDir := "processed_images"

	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		log.Fatalf("input directory '%s' does not exist", inputDir)
	}

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, os.ModePerm)
		if err != nil {
			log.Fatalf("error creating output directory '%s': %v", outputDir, err)
		}
	}
	numImageProcessors := 4
	// Channel coefficients i/o buffering size here adjusts per usecase experiment and likely tune
	imgChanSize := 4
	pngChanSize := 4

	closeCh := make(chan struct{})
	
	// Read images and feed into first stage.
	imageCh := make(chan image.Image, imgChanSize)
	// pipline channel downstream stages yet.
	blurredCh := make(chan image.Image, imgChanSize)
	resizedCh := make(chan []byte, pngChanSize)
	wrGlob := make(chan []byte, pngChanSize)
	defer close(imageCh)
	defer close(blurredCh)
	defer close(resizedCh)

	// Fanout imageLoader to many parallel processor
	for i := 0; i < numImageProcessors; i++ {
		go processRawImages(inputDir, outputDir, closeCh, wrGlob)
	}

	// Stage 1: Image Filtering.
	go pipeFilter(imageCh, blurredCh, closeCh)

	// Stage 2: Image Resizing.
	go pipeResize(blurredCh, resizedCh, 256, closeCh)
	
	
	go func() {
		signal.Notify(closeCh, os.Interrupt, os.Kill)
		select {}
	}()

	//writeCh := make(chan [][]byte,NumProducts)
	fileName := fmt.Sprintf("%s/processbenchlog.txt", outputDir)
	file, _ := os.Create(fileName)
	defer file.Close()
	scanner := bufio.NewScanner(os.Stdin)
	for. or flag.Arg
	
		write, ok := <-wrGlob
		if !ok {
			return
		}
		on := time.Now()
		err = ioutil.WriteFile(fmt.Sprintf("%s/output-%d.jpg", outputDir, count), write, 0666)
		
		if err != nil {
			fmt.Printf("Can't write to file")
		}
		fmt.Println(time.Since(on))
	}
 }