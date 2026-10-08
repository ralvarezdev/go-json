# go-json

JSON helpers for Go projects. It defines `Encoder` and `Decoder` interfaces and provides implementations on top of `encoding/json` (including streaming variants) and protobuf `protojson`. Requires Go 1.25.1 (per `go.mod`).

## Installation

```bash
go get github.com/ralvarezdev/go-json
```

Direct dependencies: `go-reflect`, `go-strings` and `google.golang.org/protobuf`.

## API

**Decoders (`decoder`)**

- **`Decoder`** — `Decode(body, dest)` and `DecodeReader(reader, dest)`; `ToReader(reader any)` utility.
- **`decoder/json`** — `NewDecoder()`, `NewStreamDecoder()`.
- **`decoder/protojson`** — `NewOptions(...)`, `NewDecoder(options)`, `NewMapper(destinationInstance)`, which maps destination struct fields including proto messages and nested structs.

**Encoders (`encoder`)**

- **`Encoder`** — `Encode(body)` and `EncodeAndWrite(writer, beforeWriteFn, body)`.
- **`ProtoJSONEncoder`** — adds `PrecomputeMarshal(body)`.
- **`encoder/json`** — `NewEncoder()`, `NewStreamEncoder()`.
- **`encoder/protojson`** — `NewOptions(...)`, `NewEncoder(options)`, `NewMapper(structInstance)`.

## Usage

```go
enc := gojsonencoder.NewEncoder()
data, err := enc.Encode(map[string]string{"hello": "world"})

dec := gojsondecoder.NewDecoder()
var out map[string]string
err = dec.DecodeReader(bytes.NewReader(data), &out)
```

## Development

```bash
go build ./...
go vet ./...
```

There are no tests.

## License

GNU General Public License v3.0 (see [LICENSE](LICENSE)).
