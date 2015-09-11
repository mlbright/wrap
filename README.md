# wrap

Wrap serves up 1 file over HTTP on the port of your choosing. Then it exits.
Download the binary for your platform, and put it in your system's PATH.

## Usage:

#### Linux: `wrap -port 8081 /tmp/some-file`

#### Windows: `wrap.exe -port 8081 c:\src\wrap\wrap.go`

## Contributing:

* Please show me an easier way
* Fork, and submit a pull request

## References:

* [Library|https://github.com/tylerb/graceful] that wraps the Golang std library net/http objects
for graceful shutdown
* Another [approach|http://www.hydrogen18.com/blog/stop-listening-http-server-go.html]
