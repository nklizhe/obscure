[![GoDoc](https://godoc.org/github.com/nklizhe/obscure?status.svg)](https://godoc.org/github.com/nklizhe/obscure) [![Build Status](https://travis-ci.org/nklizhe/obscure.svg?branch=master)](https://travis-ci.org/nklizhe/obscure)

obscure
=======
[Security through obscurity](https://en.wikipedia.org/wiki/Security_through_obscurity).

Generates Go code for embedding secret text in your code.

# Install
```
go get github.com/nklizhe/obscure/cmd/obscure
```

# Usage
```
cat my_private_key.txt | obscure
```

# Author

* LI ZHE, Tom <nklizhe@gmail.com>

# License
The MIT License (MIT)

Copyright (c) 2015 LI ZHE

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.


