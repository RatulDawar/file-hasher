package workerpool

func (wp *WorkerPool) AddJob(job Job) {
	wp.jobs <- job
}

func (wp *WorkerPool) Start() {
	for range wp.size {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for job := range wp.jobs {
				wp.results <- job.Execute(job.input)
			}
		}()
	}
}

func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

// GetResults returns the results channel for reading job results
func (wp *WorkerPool) GetResults() <-chan JobResult {
	return wp.results
}
