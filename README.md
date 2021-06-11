# gogrep

[![Go Report Card](https://goreportcard.com/badge/github.com/scorpionknifes/gogrep)](https://goreportcard.com/report/github.com/scorpionknifes/gogrep) ![gopherbadger-tag-do-not-edit](https://img.shields.io/badge/Go%20Coverage-96%25-brightgreen.svg?longCache=true&style=flat)

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
