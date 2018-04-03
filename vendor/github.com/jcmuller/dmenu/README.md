# dmenu : Simple dmenu wrapper in go

## Overview
[![GoDoc](https://godoc.org/github.com/jcmuller/dmenu?status.svg)](https://godoc.org/github.com/jcmuller/dmenu)
[![Build Status](https://travis-ci.org/jcmuller/dmenu.svg?branch=master)](https://travis-ci.org/jcmuller/dmenu)
[![Code Climate](https://codeclimate.com/github/jcmuller/dmenu/badges/gpa.svg)](https://codeclimate.com/github/jcmuller/dmenu)
[![Go Report Card](https://goreportcard.com/badge/github.com/jcmuller/dmenu)](https://goreportcard.com/report/github.com/jcmuller/dmenu)
[![Sourcegraph](https://sourcegraph.com/github.com/jcmuller/dmenu/-/badge.svg)](https://sourcegraph.com/github.com/jcmuller/dmenu?badge)

dmenu lets you interact with dmenu from go.

## Install

```
go get github.com/jcmuller/dmenu
```

## Example

```go
output, err := dmenu.Popup("Choose:", "One", "Two", "Three")

if err != nil {
    if err, ok := err.(*dmenu.EmptySelectionError); ok {
        // It's ok. It's fine to not get a selection
    } else {
        fmt.Println(fmt.Errorf("Error getting output: %s", err))
    }
}

fmt.Println(output)

```

## Author

@jcmuller

## License

[MIT](License).
