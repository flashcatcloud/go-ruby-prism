package parser

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/flashcatcloud/go-ruby-prism/wasm"
)

func TestOptionsWireLayout(t *testing.T) {
	for _, tc := range []struct {
		name, hex string
		scopes    [][][]byte
	}{
		{name: "default", hex: "0000000001000000000000000000000000000000000000"},
		{name: "scope", hex: "000000000100000000000000000000000000000100000001000000000100000078", scopes: [][][]byte{{[]byte("x")}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := serializeParserOptions(nil, 0, nil, false, nil, SyntaxVersionLatest, false, false, false, tc.scopes)
			if err != nil {
				t.Fatal(err)
			}
			want, err := hex.DecodeString(tc.hex)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("wire=%x, want %x", got, want)
			}
		})
	}
}

func TestScopeVariableAcrossRepeatedParses(t *testing.T) {
	ctx := context.Background()
	p, err := NewParser(ctx, WithScopes([][][]byte{{[]byte("x")}}))
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close(ctx)
	for i := 0; i < 200; i++ {
		res, err := p.Parse(ctx, []byte("x"))
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		if len(res.Errors) != 0 {
			t.Fatalf("iteration %d: %v", i, res.Errors)
		}
		if _, ok := res.Value.Statements.Body[0].(*LocalVariableReadNode); !ok {
			t.Fatalf("iteration %d: %T", i, res.Value.Statements.Body[0])
		}
	}
}

func TestParseMemoryLimitClosesFailedRuntime(t *testing.T) {
	ctx := context.Background()
	p, err := NewParser(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close(ctx)
	_, err = p.Parse(ctx, []byte(strings.Repeat("x", 65<<20)))
	if !errors.Is(err, wasm.ErrMemoryLimit) {
		t.Fatalf("got %v, want memory limit error", err)
	}
	if _, err = p.Parse(ctx, []byte("1")); err == nil {
		t.Fatal("failed runtime remained usable")
	}
	fresh, err := NewParser(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer fresh.Close(ctx)
	if r, err := fresh.Parse(ctx, []byte("1")); err != nil || len(r.Errors) != 0 {
		t.Fatalf("fresh parser: %v %v", r, err)
	}
}
