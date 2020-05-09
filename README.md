# geojson

This package provides types for working with [GeoJSON](https://geojson.org/)
data as defined in [RFC7946](https://tools.ietf.org/html/rfc7946).

There are two layers available:
* High-level types that provide a friendlier interface and attempt to do some
  validation. These live in the top-level package.
* Low-level types that closely match GeoJSON wire format. These live under
  `github.com/bsidhom/geojson/wire`.

Both levels can be deserialized from raw JSON, but serialization is
currently only available for the low-level types.

Parsing JSON using the low layer:

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/bsidhom/geojson/wire"
)

func main() {
    s := `{"type":"Point","coordinates":[0,1]}`
    var w wire.Wrapper
    json.Unmarshal([]byte(s), &w)
    point := w.Value.(*wire.Point)
    fmt.Printf("%#v\n", point)
}
```

Parsing JSON using the higher layer:

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/bsidhom/geojson"
)

func main() {
    s := `{"type":"Point","coordinates":[0,1]}`
    var w geojson.Wrapper
    json.Unmarshal([]byte(s), &w)
    point := w.Value.(*geojson.Point)
    fmt.Printf("%#v\n", point)
}
```
