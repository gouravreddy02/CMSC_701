# CMSC701 Assignment 1 — Suffix Array

A Go implementation of suffix array construction and querying

## Overview

This project implements a full suffix array pipeline over genomic reference sequences (FASTA format), including index construction, inspection, and pattern querying with three different search strategies.

## Tools

| Executable | Description |
|------------|-------------|
| `buildsa`  | Reads a FASTA reference, builds the suffix array and a prefix lookup table, and serializes everything to a binary file |
| `inspectsa` | Loads the binary index and outputs LCP1 statistics (mean, median, max) plus a sampled spot-check of the suffix array |
| `querysa`  | Loads the binary index and searches for patterns from a query FASTA file using one of three modes |

## Usage

### Build

```bash
bash build.sh
```

This compiles all three executables into the project root.

### buildsa

```bash
./buildsa <reference.fa> <k> <output>
```

- `reference.fa` — FASTA file containing the reference genome
- `k` — prefix length for the prefix lookup table (must be ≤ 16)
- `output` — path to write the serialized binary index

### inspectsa

```bash
./inspectsa <index> <sample_rate> <output>
```

Outputs a 4-line text file: mean LCP1, median LCP1, max LCP1, and tab-separated suffix array entries sampled at `sample_rate` intervals.

### querysa

```bash
./querysa <index> <queries.fa> <mode> <output>
```

- `mode` — one of `naive`, `simpaccel`, or `prefaccel`

Each output line is tab-separated: `query_name  lb  ub  k  hit_1 ... hit_k`

**Query modes:**
- `naive` — standard binary search; reports character comparison counts for lower and upper bound searches
- `simpaccel` — binary search with two-LCP accelerant (avoids redundant character comparisons)
- `prefaccel` — uses the precomputed prefix table to jump directly to the matching interval, then binary searches within it; reports the prefix interval bounds instead of comparison counts

## Project Structure

```
CMSC701_A1/
├── build.sh
├── go.mod
├── cmd/
│   ├── buildsa/main.go
│   ├── inspectsa/main.go
│   └── querysa/main.go
└── pkg/
    ├── fasta/fasta.go
    └── suffixarray/
        ├── build.go         # O(m² log m) suffix array construction
        ├── prefix_table.go  # Prefix lookup table construction
        ├── query.go         # naive, simpaccel, prefaccel search
        └── serialize.go     # gob-based binary serialization
```

## Language

Go

## Hardest Part

Getting the character comparison counts to match the reference implementation exactly was the most challenging aspect. The binary search logic itself was straightforward, but subtle differences in how comparisons are counted (when the suffix is shorter than the pattern, or when a full match is found) led to small discrepancies that were difficult to diagnose.

## Resources Consulted

- Course lectures and slides on suffix arrays and binary search accelerants
- Go standard library documentation (`encoding/gob`, `sort`, `bufio`)
