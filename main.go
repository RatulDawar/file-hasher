package main

import (
	"bufio"
	"context"
	"fmt"
	"go-hasher/internal/controllers/filehasher"
	"go-hasher/pkg/appcontext"
	"go-hasher/pkg/filehandler"
	"go-hasher/pkg/memorycache"
	"go-hasher/pkg/workerpool"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {

	workerPool := workerpool.NewWorkerPool(10, 100)
	workerPool.Start()

	memoryCache := memorycache.NewMemoryCache()

	appCtx := appcontext.NewAppContext(workerPool, memoryCache)
	ctx := appcontext.WithAppContext(context.Background(), appCtx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Server starting on :8081...")
	fmt.Println("Worker pool initialized with 10 workers")

	go func() {
		<-sigChan
		fmt.Println("\nShutting down server...")
		listener.Close()
		workerPool.Close()
		fmt.Println("Worker pool shut down")
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {

			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				break
			}
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleConnection(conn, ctx)
	}
}

func processMessage(ctx context.Context, message string) {

	message = strings.TrimSpace(message)

	wp := appcontext.MustGetWorkerPool(ctx)
	memoryCache := appcontext.MustGetMemoryCache(ctx)

	if hash, exists := memoryCache.Get(message); exists {
		fmt.Printf("Cache hit for %s: %s\n", message, hash)
		return
	}

	filehasher := filehasher.NewFileHasherController()

	wp.AddJob(workerpool.Job{
		Execute: func(jobInput workerpool.JobInput) workerpool.JobResult {
			hash, err := filehasher.HashFile(ctx, filehandler.Path(message))
			return workerpool.JobResult{Output: hash, Error: err}
		},
	})

	fmt.Printf("Processed and cached %s: \n", message)
}

func handleConnection(conn net.Conn, ctx context.Context) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		processMessage(ctx, message)

		conn.Write([]byte("File processed and cached!\n"))
	}
}
