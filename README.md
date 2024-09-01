# iterutils

[![Go Reference](https://pkg.go.dev/badge/github.com/mackee/iterutils.svg)](https://pkg.go.dev/github.com/mackee/iterutils)

`iterutils` is a Go library that provides utilities for working with iterators (`iter.Seq` and `iter.Seq2`). It offers various functions to simplify iteration, transform sequences, and handle concurrent operations. This library is designed to enhance the flexibility and performance of your Go applications by providing easy-to-use iterator utilities.

## Installation

To install `iterutils`, use `go get`:

```sh
go get github.com/mackee/iterutils
```

## Features

- **FromNexter**: Convert traditional iteration objects into `iter.Seq`.
- **FromNexter2**: Convert traditional iteration objects returning pairs of values into `iter.Seq2`.
- **FromTryNexter2**: Handle iteration with error handling, converting objects that return both a value and an error into `iter.Seq2`.
- **NexterWithT**: Wrap existing iteration types to allow their use with functions like `FromNexter`.
- **TryNexterWithT**: Wrap iteration types that return values and errors, allowing their use with `FromTryNexter2`.
- **async.Map**: Apply a function to each element of an iterator asynchronously.
- **async.Map2**: Apply a function to each pair of elements from a sequence asynchronously.

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/mackee/iterutils"
    "github.com/mackee/iterutils/async"
)

func main() {
    // MyNexter is an object that implements the Next() method
    myNexter := &MyNexter{}
    // MyTransform is a function that transforms the value returned by MyNexter.Next()
    itr := iterutils.FromNexter(myNexter, myTransform)
    // myAsyncTransform is a function that takes lots of time to process but can be executed concurrently
    asyncIter := async.Map(itr, myAsyncTransform)

    for elem := range asyncSeq {
        fmt.Println(elem)
    }
}
```

### `FromNexter`

Converts a traditional iterator object into an `iter.Seq`. This is useful for wrapping objects like `database/sql.Rows` or `bufio.Scanner` to be used in an `iter.Seq`.

### `FromNexter2`

Similar to `FromNexter`, but for iterators that return pairs of values, converting them into an `iter.Seq2`.

### `FromTryNexter2`

Handles iterators that return both a value and an error, converting them into an `iter.Seq2`. This is useful for handling operations like reading from a file or a network stream, where each iteration may result in an error.

### `async.Map`

Applies a function asynchronously to each element of an iterator. This is especially useful when dealing with time-consuming operations that can be executed concurrently. The order of the results in the returned `iter.Seq` is maintained.

### `async.Map2`

Similar to `async.Map`, but applies the function to each pair of elements from an `iter.Seq2`. The order of the pairs in the returned `iter.Seq2` is maintained.

## License

This library is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests to improve the library.

## Author

This library is developed and maintained by [mackee](https://github.com/mackee).
