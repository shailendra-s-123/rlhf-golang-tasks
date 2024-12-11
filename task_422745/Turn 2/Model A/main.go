package main

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"sync"
	"time"
)

const (
	// Number of workers for each stage in the pipeline
	numWorkers = 4
	// Buffer size for each channel in the pipeline
	bufferSize = 8
	// Max number of images to process
	numImages = 20
)

// Error type for pipeline errors
type PipelineError struct {
	err error
}

func (pe PipelineError) Error() string {
	return "Pipeline Error: " + pe.err.Error()
}

type Image struct {
	Filename string
	Data     []byte
}

func main() {
	// Initialize pipeline channels
	input := make(chan Image, bufferSize)
	processed := make(chan Image, bufferSize)
	output := make(chan Image, bufferSize)

	// Start worker Goroutines for each stage
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go readImages(&wg, input)
		wg.Add(1)
		go processImages(&wg, input, processed)
		wg.Add(1)
		go writeImages(&wg, processed, output)
	}

	// Generate and send sample image filenames
	for i := 1; i <= numImages; i++ {
		input <- Image{
			Filename: fmt.Sprintf("image_%03d.jpg", i),
			Data:     sampleImageData(i), // Generate sample image data
		}
	}

	// Close input channel to signal completion
	close(input)

	// Wait for all workers to finish
	wg.Wait()

	// Close remaining channels
	close(processed)
	close(output)

	fmt.Println("Pipeline processing completed.")
}

func readImages(wg *sync.WaitGroup, in chan<- Image) {
	defer wg.Done()
	for {
		select {
		case image, ok := <-in:
			if !ok {
				return
			}
			// Simulate reading image data from file
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			in <- image
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func processImages(wg *sync.WaitGroup, in <-chan Image, out chan<- Image) {
	defer wg.Done()
	for {
		select {
		case image, ok := <-in:
			if !ok {
				return
			}

			// Process the image (e.g., resize, apply filters)
			processedImage := image
			processedImage.Filename = fmt.Sprintf("%s_processed.png", image.Filename[:len(image.Filename)-4])
			out <- processedImage

		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func writeImages(wg *sync.WaitGroup, in <-chan Image, out chan<- Image) {
	defer wg.Done()
	for {
		select {
		case image, ok := <-in:
			if !ok {
				return
			}

			// Write the image to file
			err := writeImage(image.Filename, image.Data)
			if err != nil {
				out <- Image{} // Propagate error downstream
				return
			}

		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func sampleImageData(id int) []byte {
	// Generate sample image data (in this case, just return some dummy data)
	return []byte(fmt.Sprintf("Sample image data for %d", id))
}

func writeImage(filename string, data []byte) error {
	// Write image data to file
	// Replace this with actual image encoding logic using "image/jpeg" and "image/png"
	if strings.HasSuffix(filename, ".jpg") {
		return errors.New("Not implemented")
	} else {
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		_, err = f.Write(data)
		f.Close()
		return err
	}
}