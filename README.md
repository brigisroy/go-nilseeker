# NilSeeker - Go Nil Pointer Analyzer

NilSeeker is a static analysis tool for Go that detects potential nil pointer dereferences in your code. It helps you find and fix nil pointer issues before they cause panics in production.

## Features

- Detects potential nil pointer dereferences in field access and method calls
- Finds explicit nil pointer dereferences
- Detects nil slice and map indexing
- Tracks variables that have been nil-checked to reduce false positives
- Provides clear and actionable error messages

## Installation

```bash
go install github.com/brigisroy/go-nilseeker/cmd/nilseeker@latest
```

## Usage

To analyze your Go code, run:

```bash
nilseeker [flags] <packages>
```

For example, to scan all packages in your current module:

```bash
nilseeker ./...
```

Or to scan a specific package:

```bash
nilseeker ./pkg/mypackage
```

### Command Line Flags

- `-r`: Recursively scan subdirectories
- `-v`: Enable verbose output
- `--exclude-dirs`: Comma-separated list of directories to exclude
- `--exclude-files`: Comma-separated list of file patterns to exclude
- `--version`: Show version information

## Examples

### Basic Usage

```bash
# Scan current package
nilseeker .

# Scan all packages in current module
nilseeker ./...

# Scan with verbose output
nilseeker -v ./pkg/mypackage
```

### CI Integration

Add NilSeeker to your CI pipeline:

```yaml
# Example GitHub Actions workflow
steps:
  - name: Check out code
    uses: actions/checkout@v3

  - name: Set up Go
    uses: actions/setup-go@v4
    with:
      go-version: "1.20"

  - name: Install nilseeker
    run: go install github.com/brigisroy/go-nilseeker/cmd/nilseeker@latest

  - name: Run nilseeker
    run: nilseeker ./...
```

## How It Works

NilSeeker uses Go's `go/analysis` framework to analyze your code. It examines:

1. **Selector expressions** (`a.b`, `a.Method()`) - Checks if `a` could be nil
2. **Star expressions** (`*p`) - Checks if `p` could be nil
3. **Index expressions** (`slice[i]`, `map[key]`) - Checks if the slice or map could be nil
4. **Nil checks** (`if x != nil { ... }`) - Tracks which variables have been checked

## Limitations

- False positives: NilSeeker may report issues that aren't actual problems in your code
- Limited interprocedural analysis: NilSeeker primarily operates within function boundaries
- No flow-sensitive analysis: NilSeeker doesn't fully track program flow to determine variable states
