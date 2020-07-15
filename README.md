# websqlview

A combination of things:

* A webview with access to sqlite
* A mini framework for creating desktop-like webapps

This is based upon [webview](https://github.com/zserge/webview/) and should (eventually) be cross-platform (Windows/MacOS/Linux) although it has only been developed and tested on Ubuntu/libgtk so far.

It exposes GO's [sqlite](https://pkg.go.dev/github.com/mattn/go-sqlite3)/[database/sql](https://pkg.go.dev/database/sql) API (and some other assorted APIs) to javascript so webapps running in websqlview can have native desktop-like capabilities.

## Building

Install GO on your system. Then:

```sh
go build
```

to compile `websqlview.go` (and associated components).

If you are using an old-ish version of libgtk (ie. 3.18), you may need to supply `-tags gtk_majorversion_minorversion` ie.

```sh
go build -tags gtk_3_18
```

## Usage

For development and demonstration purposes, you can probably use `websqlview` directly:

```sh
websqlview file:///<something>
```

For release, you probably want to fork `websqlview.go` and customise to suit.

## Examples

* [`examples/dialog-example.html`](examples/dialog-example.html)
* [`sqlite/test/20-basic.test.js`](sqlite/test/20-basic.test.js)
* [GigoBooks](https://github.com/gigobooks/gigobooks) - Clean and simple accounting software for solopreneurs, consultants, freelancers and other micro-businesses.
