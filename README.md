# GAE/Go data deserialization benchmark

## Introduction

Inspired by <https://github.com/alecthomas/go_serialization_benchmarks>,
this benchmark measures time and memory usage in
loading/deserializing a constant asset on
[GAE/Go standard runtime](https://cloud.google.com/appengine/docs/standard/go/).

The asset mentioned here is like
a ~100k-element list of struct declared in `data.go`,
which mimics our actual use case.

```
type Data struct {
	Bool   bool      `json:"bool" msg:"bool"`
	Int    int       `json:"int" msg:"int"`
	Float  float64   `json:"float" msg:"float"`
	String string    `json:"string" msg:"string"`
	Time   time.Time `json:"time" msg:"time"`
}
```

An asset for benchmark is generated with a specific random seed,
then stored and serialized in various ways, referred as the following:

- `embed` ... embedded in Go source code as `[]Data` variable
- `json` ... stored in a file using [encoding/json](https://golang.org/pkg/encoding/json/)
- `gob` ... stored in a file using [encoding/gob](https://golang.org/pkg/encoding/gob/)
- `msgp` ... stored in a file using [github.com/tinylib/msgp](https://github.com/tinylib/msgp)
- `msgpack` ... stored in a file using [github.com/vmihailenco/msgpack](https://github.com/vmihailenco/msgpack)

In benchmark results, a notation like `json100000` is used,
which means a dataset of `json` with 100000-element list of Data.

In each result there's an SHA1 hash to check integrity of deserialized data.
They should be identical for the same list element size (100000).

```
Name            Count           Time                    Memory Bytes    Memory Allocs           Hash
embed100000     2000000000               0.36 ns/op            0 B/op          0 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
```

## Instructions

### Building

First you need to install `msgp` tool to generate msgpack encoder/decoder.

```
$ go get -u github.com/tinylib/msgp
```

Clone this repository under `$GOPATH/src` and run `go generate .`
to generate GAE/Go deployables in `app` directory.

```
$ cd $GOPATH/src
$ git clone https://github.com/t-yaegashi/goddbench
$ cd goddbench
$ go generate .
======== MessagePack Code Generator =======
>>> Input: "data.go"
>>> Wrote and formatted "data_gen.go"
2018/01/06 21:49:52 Writing to app/embed100000/list.go
2018/01/06 21:50:04 Writing to app/embed100000/app.yaml
2018/01/06 21:50:04 Writing to app/embed100000/app.go
2018/01/06 21:50:05 Writing to app/json100000/list.go
2018/01/06 21:50:05 Writing to app/json100000/list.json
2018/01/06 21:50:05 Writing to app/json100000/app.yaml
2018/01/06 21:50:05 Writing to app/json100000/app.go
2018/01/06 21:50:05 Writing to app/msgpack100000/list.go
2018/01/06 21:50:05 Writing to app/msgpack100000/list.msgpack
2018/01/06 21:50:07 Writing to app/msgpack100000/app.yaml
2018/01/06 21:50:07 Writing to app/msgpack100000/app.go
2018/01/06 21:50:08 Writing to app/gob100000/list.go
2018/01/06 21:50:08 Writing to app/gob100000/list.gob
2018/01/06 21:50:08 Writing to app/gob100000/app.yaml
2018/01/06 21:50:08 Writing to app/gob100000/app.go
2018/01/06 21:50:08 Writing to app/msgp100000/list.go
2018/01/06 21:50:08 Writing to app/msgp100000/list.msgp
2018/01/06 21:50:08 Writing to app/msgp100000/app.yaml
2018/01/06 21:50:08 Writing to app/msgp100000/app.go
```

### Testing locally

Launch `dev_appserver.py` for all deployables to test locally:

```
$ dev_appserver.py app/*
```

Use `curl` to retrieve the results in another console:

```
$ curl http://localhost:808{0..5}
embed100000     2000000000               0.36 ns/op            0 B/op          0 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
gob100000             20          59716073 ns/op        20654066 B/op     100237 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
json100000             2         524417147 ns/op        85308624 B/op     400024 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
msgp100000            20          74252328 ns/op        17230405 B/op     200005 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
msgpack100000         10         175795553 ns/op        14219875 B/op     700012 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
```

### Running on GAE/Go

Use `gcloud` to deploy on the real GAE/Go environment:

```
$ glcoud app deploy --project PROJECT app/*
```

You can retrieve the results with the following URLs:

```
https://goddbench-embed100000-dot-PROJECT.appspot.com/
https://goddbench-json100000-dot-PROJECT.appspot.com/
https://goddbench-gob100000-dot-PROJECT.appspot.com/
https://goddbench-msgpack100000-dot-PROJECT.appspot.com/
https://goddbench-msgp100000-dot-PROJECT.appspot.com/
```

Run `curl.sh` to get all of them in one action:

```
./curl.sh PROJECT
```

## Results

Result using `dev_appserver.py` on n1-standard-1:
```
Name            Count           Time                    Memory Bytes    Memory Allocs           Hash
embed100000     2000000000               0.36 ns/op            0 B/op          0 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
gob100000             20          59716073 ns/op        20654066 B/op     100237 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
json100000             2         524417147 ns/op        85308624 B/op     400024 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
msgp100000            20          74252328 ns/op        17230405 B/op     200005 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
msgpack100000         10         175795553 ns/op        14219875 B/op     700012 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
```

Result on GAE in us-central:
```
Name            Count           Time                    Memory Bytes    Memory Allocs           Hash
gob100000             10         176936904 ns/op        45133789 B/op     109994 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
json100000             1        1467064292 ns/op        152637112 B/op    418667 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
msgp100000             5         261154645 ns/op        46860164 B/op     239570 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
msgpack100000          3         423261986 ns/op        42541709 B/op     723842 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
```

> Note: `embed100000` failed in deployment due to precompilation time out.

Result on GAE in asia-northeast1:
```
Name            Count           Time                    Memory Bytes    Memory Allocs           Hash
embed100000     2000000000               0.45 ns/op            0 B/op          0 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
gob100000             10         129381925 ns/op        45133789 B/op     109994 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
json100000             2         707772450 ns/op        125792668 B/op    416212 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
msgp100000            10         145570381 ns/op        46859890 B/op     239569 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
msgpack100000          5         271111025 ns/op        42541668 B/op     723841 allocs/op      1ba1706bea3a667d765e505133e4261d1652371c
```

## Thoughts

- `embed` is doubtlessly the fastest and the most efficient solution, but some GAE regions including us-central are unable to deploy huge golang sources such as `embed100000` due to compilation time out.  So we have to seek for alternative runtime solutions.
- `json`, `gob`, `msgpack` are easy-to-integrate solutions because they are fully reflection-based (no need for writing IDLs or generating code) and provide the same API for `io.Reader` / `io.Writer` streaming.
- `msgp` needs to generate dedicated code for each type to encode/decode it, but the code generator (msgp) is reflection aware, so you don't have to write any IDLs by yourself.
- As expected, `json` was the worst in both speed and memory consumption.
- I expected `msgp` would be the best in performance, but actually `gob` was faster and less memory consuming overall.
- [GAE standard refuses a file larger than 32MB](https://cloud.google.com/appengine/quotas#Code).
