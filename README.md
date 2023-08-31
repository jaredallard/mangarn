# mangarn

A tool to bulk rename files based on regex patterns. Primarily developed
for renaming manga exports/downloads from various sources.

**Note**: Currently designed to work best with images per page.

## What's supported?

See [./internal/parser/parser_test.go](./internal/parser/parser_test.go)
for tested file name patterns.

## Usage

```bash
# Run the CLI in the current directory, outputting files into the
# ./output directory.
go run ./cmd/mangarn
```

## License

GPL-3.0-or-later
