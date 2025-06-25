Projects Focused on Inverted Index Concepts
1. dwayhs/go-search-engine

A beginner-friendly, pure-Go inverted-index search engine (no mmap but great structure):

    Implements indexing, query parsing, and TF-IDF ranking

    Clear code split across index/, core/, and search/ modules
    github.com+15github.com+15github.com+15
    github.com+1github.com+1
    github.com+1github.com+1

Why it helps: Though it uses in-memory maps, it clearly shows how to structure indexing, querying, ranking — a great reference to layer mmap + binary into.
2. el10savio/Inverted-Index-Generator

A simple Go example for generating inverted indexes from document collections.

Why it helps: Focuses on core inverted-index building logic — tokenization, storing word→doc mapping — which you can adapt to a mmap-backed system with fixed struct entries.
3. CyberKatze/simple-binary-search-engine

Small Go binary that builds a search engine with a postings list binary structure.
github.com+4github.com+4github.com+4

Why it helps: Operates on byte-level binary layout, similar to what your mmap + binary.PutUintXX code will use.
📚 Low-Level mmap in Go
edsrzf/mmap-go

A popular, cross-platform Go mmap package with example usage and tests.
github.com+11github.com+11github.com+11

Why it helps: While it doesn’t build search features, it's an excellent guide to correct mmap usage, unmapping, remapping, and syncing — the foundational I/O layer your engine needs.
💡 How to Make the Most of These

    Learn data structure flow from go-search-engine — how it tokenizes, indexes, and scores.

    Study binary layout in CyberKatze’s project to understand how postings lists are written as bytes.

    Implement mmap I/O using patterns from mmap-go.

    Integrate both:

        Replace in-memory maps with memory-mapped files

        Replace slices with structured byte offsets and lengths

        Use binary.PutUint32/64 and Uint32/64 for fixed-length records

🏗️ Step-by-Step Suggested Path

Step 1: Clone go-search-engine, run tests, and get a feel for its patterns.

Step 2: Build a small experiment with mmap-go to write/read a struct into a memory-mapped file:

buf := mmap[0:12]
// write offset + length values using binary.PutUintXX

Step 3: Modify go-search-engine to persist its inverted index into mmap-backed files using your buffer logic.

Step 4: Swap in your own search structures, and re-run queries — verifying correctness and speed gains.
📌 Summary Table
Project	Focus Area	How It Helps You
go-search-engine	Search logic, postings list	Understand structure & ranking
simple-binary-search-engine	Byte-level postings	Learn binary layout for postings
mmap-go	Memory-mapped file I/O	Correct usage of mmap, sync, remap basics

These real-world examples provide clear, modular foundations upon which to build your own high-performance, disk-backed search engine using mmap + binary encoding. 🚀 Let me know if you'd like the relevant code snippets extracted or mapped to your project structure!

https://dev.to/ravikishan/building-a-high-performance-full-text-search-engine-in-go-m4b#:~:text=This%20blog%20will%20walk%20you,performance%20applications

https://artem.krylysov.com/blog/2020/07/28/lets-build-a-full-text-search-engine/