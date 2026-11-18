package wasm

import (
	"context"
	"errors"
	"testing"
)

func TestLinearMemoryCap(t *testing.T) {
	ctx := context.Background()
	r, err := NewRuntime(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close(ctx)
	before := r.mem.mem.Size()
	if _, ok := r.mem.mem.Grow(MemoryLimitPages); ok {
		t.Fatal("memory grew beyond cap")
	}
	if r.mem.mem.Size() != before {
		t.Fatal("rejected grow changed memory")
	}
	if _, err := r.Calloc(ctx, 64<<20, 1); !errors.Is(err, ErrMemoryLimit) {
		t.Fatalf("calloc: %v", err)
	}
	if r.mem.mem.Size() > 64<<20 {
		t.Fatalf("linear bytes=%d", r.mem.mem.Size())
	}
}
