# Random String

[![Build Status](https://travis-ci.org/neptulon/randstr.svg?branch=master)](https://travis-ci.org/neptulon/randstr)
[![GoDoc](https://godoc.org/github.com/neptulon/randstr?status.svg)](https://godoc.org/github.com/neptulon/randstr)

Random string generator of arbitrary size.

## Example

```go
import (
	"github.com/neptulon/randstr"
)

func main() {
	str := randstr.Get(96)
	// str: "VkNT!pQXdtHgyffWMIqNZcnOECWhVYYafBGTDjJvE  PlyaWs!UKiKxGQkquNafewfcU ECXgQfYtyZkFIXEJmIYVPRYaIzh"
}
```

## Alphabet

The default alphabet contains 'space', 'dot', 'exclamation mark', and upper and lowercase English alphabetic characters:

```go
var Alphabet = []rune(". !abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
```

You can replace the alphabet used in creating the random string via reassigning the `Alphabet` variable:

```go
randstr.Alphabet = []rune("abcdefg1234567")
```

## License

[MIT](LICENSE)
