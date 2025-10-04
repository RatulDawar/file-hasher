package workerpool

import (
	"crypto/sha256"
	"testing"
)

func TestWorkerPool(t *testing.T) {
	wp := NewWorkerPool(3, 10)
	wp.Start()

	numJobs := 10
	for i := 0; i < numJobs; i++ {
		wp.AddJob(Job{
			Execute: func(jobInput JobInput) JobResult {
				value := jobInput.input.(int)
				hash := sha256.Sum256([]byte{byte(value)})
				return JobResult{Output: hash, Error: nil}
			},
			input: JobInput{input: i},
		})
	}
	wp.Close()

	results := make([][32]byte, 0)
	for jobResult := range wp.results {
		if jobResult.Error != nil {
			t.Error(jobResult.Error)
		}
		results = append(results, jobResult.Output.([32]byte))
	}

	expectedResults := make(map[[32]byte]bool)
	for i := 0; i < numJobs; i++ {
		expectedResults[sha256.Sum256([]byte{byte(i)})] = true
	}

	if len(results) != numJobs {
		t.Errorf("Expected %d results, got %d", numJobs, len(results))
	}

	for _, result := range results {
		if !expectedResults[result] {
			t.Error("Unexpected result found")
		}
	}
}
