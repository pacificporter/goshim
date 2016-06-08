goshim
======

## Description

Better `go run`. Build go codes transparently and exec

## Installation

    % go get github.com/pacificporter/goshim/cmd/goshim

## Synopsis

    % goshim ./path/to/pkg [args...]

### Help

    % goshim -h

### Force rebuilding

    % goshim -f ./path/to/pkg [args...]

### Verbosely output

    % goshim -v ./path/to/pkg [args...]

### Combination

    % goshim -f -v ./path/to/pkg [args...]

## License

[MIT][license]

## Author

* [Songmu](https://github.com/Songmu)
    * Original
* [Pacificporter Inc.](https://github.com/pacificporter)
    * Use flag library
    * Add -f, -v flags
