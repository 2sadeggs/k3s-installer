package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 4*1024)
	},
}

// Display progress with a given frequency
func displayProgress(iteration, totalIterations, frequency int, operation string) {
	if iteration%frequency == 0 {
		progress := float64(iteration) / float64(totalIterations) * 100
		fmt.Printf("\r%s Progress: %.2f%%", operation, progress)
	}
}

// Test sequential read/write speed
func testSequential(filename string, size int) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Failed to open file %s: %v", filename, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	buf := make([]byte, 1024*1024)
	startTime := time.Now()

	// Sequential Write
	fmt.Println("Sequential Write in progress...")
	for i := 0; i < size/len(buf); i++ {
		_, err := f.Write(buf)
		if err != nil {
			return fmt.Errorf("Error writing to file: %v", err)
		}
		displayProgress(i, size/len(buf), 1000, "Sequential Write")
	}

	writeDuration := time.Since(startTime).Seconds()
	writeSpeed := float64(size) / (1024 * 1024) / writeDuration
	fmt.Printf("\nSequential Write complete. Write Speed: %.2f MB/s\n", writeSpeed)

	f.Seek(0, 0)
	startTime = time.Now()

	// Sequential Read
	fmt.Println("Sequential Read in progress...")
	for i := 0; i < size/len(buf); i++ {
		_, err := f.Read(buf)
		if err != nil {
			return fmt.Errorf("Error reading from file: %v", err)
		}
		displayProgress(i, size/len(buf), 1000, "Sequential Read")
	}

	readDuration := time.Since(startTime).Seconds()
	readSpeed := float64(size) / (1024 * 1024) / readDuration
	fmt.Printf("\nSequential Read complete. Read Speed: %.2f MB/s\n", readSpeed)

	return nil
}

// Test random read/write IOPS
func testRandomIO(filename string, size int) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Failed to open file %s: %v", filename, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	var wg sync.WaitGroup

	// Random Write
	fmt.Println("Random Write in progress...")
	startTimeWrite := time.Now()
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := bufPool.Get().([]byte)
			defer bufPool.Put(buf)

			for i := 0; i < size/len(buf); i++ {
				offset := rand.Int63n(int64(size))
				_, err := f.WriteAt(buf, offset)
				if err != nil {
					fmt.Printf("Error writing to file: %v\n", err)
				}
				displayProgress(i, size/len(buf), 1000, "Random Write")
			}
		}()
	}

	wg.Wait()
	writeDuration := time.Since(startTimeWrite).Seconds()
	writeIOPS := float64(size) / (1024 * 1024) / writeDuration
	fmt.Printf("\nRandom Write complete. Write IOPS: %.0f\n", writeIOPS)

	// Random Read
	fmt.Println("Random Read in progress...")
	startTimeRead := time.Now()
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := bufPool.Get().([]byte)
			defer bufPool.Put(buf)

			for i := 0; i < size/len(buf); i++ {
				offset := rand.Int63n(int64(size))
				_, err := f.ReadAt(buf, offset)
				if err != nil {
					fmt.Printf("Error reading from file: %v\n", err)
				}
				displayProgress(i, size/len(buf), 1000, "Random Read")
			}
		}()
	}

	wg.Wait()
	readDuration := time.Since(startTimeRead).Seconds()
	readIOPS := float64(size) / (1024 * 1024) / readDuration
	fmt.Printf("\nRandom Read complete. Read IOPS: %.0f\n", readIOPS)

	return nil
}

func diskRWTest() error {
	testFileSize := 10 * 1024 * 1024 * 1024 // 10GB

	fmt.Println("Sequential read/write test start...(10GB file)")
	err := testSequential("disk_test_file.dat", testFileSize)
	if err != nil {
		return fmt.Errorf("Sequential read/write test failed: %v", err)
	}
	fmt.Println("Sequential read/write test complete!")

	fmt.Println("Random read/write IOPS test start...(10GB file)")
	err = testRandomIO("disk_test_file.dat", testFileSize)
	if err != nil {
		return fmt.Errorf("Random read/write IOPS test failed: %v", err)
	}
	fmt.Println("Random read/write IOPS test complete!")

	return nil
}
