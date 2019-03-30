# bigdur

A small Go package for parsing larger duration units.

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/go.coder.com/bigdur)

## Install

```bash
go get -u go.coder.com/bigdur
```

## Overview

A duration token consists of a series of coefficient and unit pairs.

For example:

- `4d`
- `4m4s`
- `4mo2.2d5s`

are valid.

The following units are available: 

| Abbreviation | Description |
|--------------|-------------|
| s            | 1 second    |
| m            | 60 seconds  |
| h            | 60 minutes  |
| d            | 24 hours    |
| w            | 7 days      |
| mo           | 30 days     |
| y            | 12 months   |


