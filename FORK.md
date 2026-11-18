# Flashcat fork

Based on upstream v1.1.0, retaining its bundled prism.wasm and generated AST.
The Go module path is github.com/flashcatcloud/go-ruby-prism.

The bundled Prism submodule is 1463122480f078562809d8f0c9802419c2fb3c31.
Its options reader requires a freeze byte before scopes_count and a forwarding
byte after each locals_count. Both were missing from the Go serializer, which
could misread allocator contents as huge counts after repeated parsing.

Each Prism instance is limited to 64 MiB of linear memory. Failed allocations
are errors; a failed Parse closes its runtime, so callers must discard and
recreate that parser. This budget does not include Go/compiler overhead and is
not the Ruby execution runtime's memory setting.
