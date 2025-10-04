package appcontext

import (
	"context"
	"go-hasher/pkg/memorycache"
	"go-hasher/pkg/workerpool"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// WorkerPoolKey is the context key for the worker pool
	WorkerPoolKey ContextKey = "workerPool"
	// MemoryCacheKey is the context key for the memory cache
	MemoryCacheKey ContextKey = "memoryCache"
)

// AppContext holds application-wide dependencies
type AppContext struct {
	WorkerPool  *workerpool.WorkerPool
	MemoryCache *memorycache.MemoryCache
}

// NewAppContext creates a new application context with the given dependencies
func NewAppContext(workerPool *workerpool.WorkerPool, memoryCache *memorycache.MemoryCache) *AppContext {
	return &AppContext{
		WorkerPool:  workerPool,
		MemoryCache: memoryCache,
	}
}

// WithAppContext adds the application context to the given context
func WithAppContext(ctx context.Context, appCtx *AppContext) context.Context {
	ctx = context.WithValue(ctx, WorkerPoolKey, appCtx.WorkerPool)
	ctx = context.WithValue(ctx, MemoryCacheKey, appCtx.MemoryCache)
	return ctx
}

// GetWorkerPool retrieves the worker pool from the context
func GetWorkerPool(ctx context.Context) *workerpool.WorkerPool {
	if wp, ok := ctx.Value(WorkerPoolKey).(*workerpool.WorkerPool); ok {
		return wp
	}
	return nil
}

// GetMemoryCache retrieves the memory cache from the context
func GetMemoryCache(ctx context.Context) *memorycache.MemoryCache {
	if mc, ok := ctx.Value(MemoryCacheKey).(*memorycache.MemoryCache); ok {
		return mc
	}
	return nil
}

// MustGetWorkerPool retrieves the worker pool from the context and panics if not found
func MustGetWorkerPool(ctx context.Context) *workerpool.WorkerPool {
	wp := GetWorkerPool(ctx)
	if wp == nil {
		panic("worker pool not found in context")
	}
	return wp
}

// MustGetMemoryCache retrieves the memory cache from the context and panics if not found
func MustGetMemoryCache(ctx context.Context) *memorycache.MemoryCache {
	mc := GetMemoryCache(ctx)
	if mc == nil {
		panic("memory cache not found in context")
	}
	return mc
}
