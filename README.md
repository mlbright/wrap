wrap
========

Wrap serves up a file over HTTP on the port of your choosing. Then it exits.

## Installation

To install, you could do:

```
git clone https://github.com/mlbright/wrap.git
cd wrap
go build wrap.go
```

Then put the binary in your system's PATH.

## Usage

Linux and OS X

```
wrap -port 8081 /tmp/some-file
```

Windows

```
wrap.exe -port 8081 c:\src\wrap\wrap.go
```

## Contributing

If you would like to contribute, please:

1. Create a GitHub issue regarding the contribution. Features and bugs should be discussed beforehand.
2. Fork the repository.
3. Create a pull request with your solution. This pull request should reference and close the issues (Fix #2).

All pull requests should:

1. Be `go fmt` formatted.

## References

* [Graceful](https://github.com/tylerb/graceful) is a library that wraps the Golang std library net/http objects
for graceful shutdown.
* Another [approach](http://www.hydrogen18.com/blog/stop-listening-http-server-go.html).
* [Manners](https://github.com/braintree/manners) is what is used
