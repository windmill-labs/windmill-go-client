# windmill-go-client

Go client for the windmill platform.

## Import

Import user-friendly client:

```go
import wmill "github.com/windmill-labs/windmill-go-client"

```

Import full api from the autogenerated openapi client:

```go
import api "github.com/windmill-labs/windmill-go-client/api"
```

## Usage

```go
a, _ := wmill.GetResource("u/ruben-user/test")
a, _ := wmill.GetVariable("u/ruben-user/test")
```
