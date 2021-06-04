# gogrep

Command line tool to search patterns inside a directory

## Features

- Multiple goroutine support
- Regex pattern matching with multiline matching
- Respects gitignore

## Usage

```bash
# Pattern
gogrep PATTERN [PATH]

# Examples
gogrep helloworld
gogrep helloworld .
gogrep "hello world" ./data
```