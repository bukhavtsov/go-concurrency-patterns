/*
In Go, a pipeline is just a way to organize concurrent programs. It works like an assembly line, where data moves step by step through different stages.

Each stage is a group of goroutines doing the same task.
Data flows between stages using channels.

How it works:
- A stage receives data from the previous stage
through a channel.
- It processes the data (modifies it or creates new data).
- It sends the output to the next stage via another channel.
*/
package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type File struct {
	Name string
	Data string
}

type Stage interface {
	Execute(ctx context.Context, inFile <-chan File) <-chan File
}

type Pipeline []Stage

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) AddStage(stage Stage) *Pipeline {
	*p = append(*p, stage)
	return p
}

// Execute starts pipeline stages.
// Each stage processes one file at a time, but as soon as a stage finishes with a file,
// it immediately starts on the next one, ensuring all stages are continuously working on different files without delay.
func (p *Pipeline) Execute(ctx context.Context, inFile <-chan File) <-chan File {
	nextChan := inFile
	for _, stage := range *p {
		nextChan = stage.Execute(ctx, nextChan)
	}
	return nextChan
}

type Stage1 struct{}

func (s Stage1) Execute(ctx context.Context, inFile <-chan File) <-chan File {
	out := make(chan File)
	go func() {
		defer close(out)
		for file := range inFile {
			file.Data += " -> Stage1"
			out <- file
		}
	}()
	return out
}

type Stage2 struct{}

func (s Stage2) Execute(ctx context.Context, inFile <-chan File) <-chan File {
	out := make(chan File)
	go func() {
		defer close(out)
		for file := range inFile {
			file.Data += " -> Stage2"
			out <- file
		}
	}()
	return out
}

type Stage3 struct{}

func (s Stage3) Execute(ctx context.Context, inFile <-chan File) <-chan File {
	out := make(chan File)
	go func() {
		defer close(out)
		for file := range inFile {
			file.Data += " -> Stage3"
			time.Sleep(time.Millisecond)
			out <- file
		}
	}()
	return out
}

type Stage4 struct{}

func (s Stage4) Execute(ctx context.Context, inFile <-chan File) <-chan File {
	out := make(chan File)
	go func() {
		defer close(out)
		for file := range inFile {
			file.Data += " -> Stage4"
			out <- file
		}
	}()
	return out
}

func fileChanGenerator(files []File) <-chan File {
	out := make(chan File)

	go func() {
		defer close(out)
		for _, f := range files {
			out <- f
		}
	}()

	return out
}

func generateRandomFiles(numFiles int) []File {
	var files []File
	rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < numFiles; i++ {
		// Create random file names and data
		fileName := fmt.Sprintf("file%d", i+1)
		fileData := fmt.Sprintf("Data for %s", fileName)

		// Append the file to the slice
		files = append(files, File{Name: fileName, Data: fileData})
	}

	return files
}

func main() {
	start := time.Now()

	ctx := context.Background()

	// Create input channel
	files := generateRandomFiles(1000)
	input := fileChanGenerator(files)

	// Create and execute pipeline
	pipeline := NewPipeline()
	pipeline.
		AddStage(Stage1{}).
		AddStage(Stage2{}).
		AddStage(Stage3{}).
		AddStage(Stage4{})

	processedFiles := pipeline.Execute(ctx, input)
	for file := range processedFiles {
		fmt.Println("Processed File:", file.Name, "Data:", file.Data)
	}

	end := time.Since(start)
	fmt.Println("it took: ", end)
}
