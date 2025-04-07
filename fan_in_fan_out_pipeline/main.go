package main

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

type File struct {
	Name string
	Data string
}

type Step interface {
	Execute(ctx context.Context, inFile <-chan File) <-chan File
}

type Pipeline []Step

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) AddStep(step Step) *Pipeline {
	*p = append(*p, step)
	return p
}

func (p *Pipeline) Execute(ctx context.Context, inFile <-chan File) <-chan File {
	nextChan := inFile
	for _, step := range *p {
		// FAN-OUT: Distribute work among multiple workers
		nextChan = fanOut(ctx, step, nextChan, runtime.NumCPU())
	}
	return nextChan
}

// FAN-OUT: Distributes workload among multiple workers
func fanOut(ctx context.Context, step Step, inFile <-chan File, workerCount int) <-chan File {
	out := make(chan File)
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range step.Execute(ctx, inFile) {
				out <- file
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// FAN-IN: Merges multiple input channels into a single channel
func fanIn(ctx context.Context, channels ...<-chan File) <-chan File {
	out := make(chan File)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan File) {
			defer wg.Done()
			for file := range c {
				out <- file
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

type Step1 struct{}

func (s Step1) Execute(ctx context.Context, inFile <-chan File) <-chan File {
	out := make(chan File)
	go func() {
		defer close(out)
		for file := range inFile {
			file.Data += fmt.Sprintf("-> Step1[%s]", file.Name)
			out <- file
		}
	}()
	return out
}

type Step2 struct{}

func (s Step2) Execute(ctx context.Context, inFile <-chan File) <-chan File {
	out := make(chan File)
	go func() {
		defer close(out)
		for file := range inFile {
			file.Data += fmt.Sprintf("-> Step2[%s]", file.Name)
			out <- file
		}
	}()
	return out
}

type Step3 struct{}

func (s Step3) Execute(ctx context.Context, inFile <-chan File) <-chan File {
	out := make(chan File)
	go func() {
		defer close(out)
		for file := range inFile {
			file.Data += fmt.Sprintf("-> Step3[%s]", file.Name)
			time.Sleep(time.Millisecond)
			out <- file
		}
	}()
	return out
}

type Step4 struct{}

func (s Step4) Execute(ctx context.Context, inFile <-chan File) <-chan File {
	out := make(chan File)
	go func() {
		defer close(out)
		for file := range inFile {
			file.Data += fmt.Sprintf("-> Step4[%s]", file.Name)
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
		fileName := fmt.Sprintf("file%d", i+1)
		fileData := fmt.Sprintf("Data for %s", fileName)
		files = append(files, File{Name: fileName, Data: fileData})
	}

	return files
}

func main() {
	start := time.Now()
	ctx := context.Background()

	// FAN-IN: Combining multiple sources
	files := generateRandomFiles(1000)
	source1 := fileChanGenerator(files[:500])
	source2 := fileChanGenerator(files[500:])
	input := fanIn(ctx, source1, source2)

	// Create and execute a pipeline
	pipeline := NewPipeline()
	pipeline.
		AddStep(Step1{}).
		AddStep(Step2{}).
		AddStep(Step3{}).
		AddStep(Step4{})

	processedFiles := pipeline.Execute(ctx, input)
	for file := range processedFiles {
		fmt.Println("Processed File:", file.Name, "Data:", file.Data)
	}

	end := time.Since(start)
	fmt.Println("it took: ", end)
}
