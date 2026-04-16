# CMSC701 Assignment 2 — Bitvector Rank/Select & Sparse Array

A Go implementation of succinct bitvector rank and select operations, applied to build a sparse array data structure.

## Overview

This project implements Jacobson's rank using a two-level index (superblocks and blocks) with popcount-based within-block rank. Select is answered in O(log N) time via binary search over the rank structure. These primitives are then used to build a sparse array that efficiently maps string values to positions in a potentially large, mostly-empty index space.

## Tools

| Executable | Description |
|------------|-------------|
| `rsbuild` | Reads a binary bitvector file, builds a rank index with packed superblock/block counters, and serializes the bitvector + index to a binary output file |
| `rsquery_rank` | Loads the index built by `rsbuild` and answers rank queries from a text query file |
| `rsquery_select` | Loads the index built by `rsbuild` and answers select queries using binary search over rank |
| `sarray` | Builds and queries a sparse array backed by bitvector rank/select |

## Usage

### Build

```bash
bash build.sh
```

This compiles all four executables into the project root.

### rsbuild

```bash
./rsbuild <input_bitvector_file> <output_index_path>
```

- `input_bitvector_file` — binary file containing a bitvector (little-endian: 4-byte size in bits, then ceil(N/8) bytes of data)
- `output_index_path` — path to write the serialized bitvector and rank index

### rsquery_rank

```bash
./rsquery_rank <index> <query_file> <output_path>
```

- `index` — path to the index file produced by `rsbuild`
- `query_file` — text file with one `rank <idx>` query per line
- `output_path` — output file; each line is `<idx>:<rank>` (or `<idx>:-1` if out of bounds)

### rsquery_select

```bash
./rsquery_select <index> <query_file> <output_path>
```

- `index` — path to the index file produced by `rsbuild`
- `query_file` — text file with one `select <rank>` query per line
- `output_path` — output file; each line is `<rank>:<select>` (or `<rank>:-1` if rank exceeds maximum)

### sarray

```bash
./sarray <command_script> <output_path>
```

- `command_script` — text file containing init, insert, finish, and query commands
- `output_path` — output file with one result per query

**Supported commands:**

| Command | Output Format | Description |
|---------|---------------|-------------|
| `query_index <idx>` | `idx:<idx>:<key>` or `idx:<idx>:-1` | Return the key at a given index |
| `query_rank <rank>` | `qr:<rank>:<key>` or `qr:<rank>:-1` | Return the key at a given rank |
| `rank_at_index <idx>` | `rai:<idx>:<rank>` | Return the bitvector rank at a given index |
| `index_of_next_key <idx>` | `ink:<idx>:<pos>` or `ink:<idx>:-1` | Return the index of the next present key after idx |
| `density` | floating point number (5+ digits) | Return the density of the bitvector (ones / length) |

## Project Structure

```
CMSC701_A2/
├── build.sh
├── go.mod
├── bitvector/
│   ├── bitvector.go       # Core bitvector type and access
│   ├── intvector.go       # Packed integer vector (arbitrary bit-width)
│   ├── rank.go            # Jacobson's rank with superblock/block index
│   ├── select.go          # Select via binary search over rank
│   └── serialization.go   # Binary serialization of bitvector + index
├── cmd/
│   ├── rsbuild/main.go
│   ├── rsquery_rank/main.go
│   ├── rsquery_select/main.go
│   └── sarray/main.go
└── project_2_sample_data/
    ├── test_input/        # Sample bitvectors and query files
    └── test_output/       # Expected outputs for validation
```

## Language

Go

## Hardest Part

The hardest part was getting select to work correctly. Select uses binary search on top of the rank data structure to find answers in O(log N) time. The main challenge was understanding what select actually means, it does not just find where the r-th 1-bit is. Instead, it finds the highest position in the bitvector where the rank still equals r. Getting the binary search condition right for this took some time. There were also edge cases to handle, like select(0), which should return the position right before the very first 1-bit instead of returning -1. On top of that, the sparse array needed a different version of select that returns the actual position of the k-th 1-bit (not the highest position with that rank), so I had to implement two separate binary search functions with slightly different conditions.

The `IntVector` used for the rank index superblock and block counters is entirely self-implemented, no external library was used.

## Resources Consulted

- Course lectures on Jacobson's rank and succinct data structures
- Go standard library documentation (`math/bits`, `encoding/binary`)
