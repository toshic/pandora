---
title: JSON provider
description: JSON provider converts to Payload json files
categories: [Provider]
tags: [provider, json]
weight: 20
---

The provider reads the payload files and passes a structure of the desired type to the generators.

Since you cannot use your own structures in the default build, the default provider returns
`map[string]interface{}{}`

But this provider is convenient for creating custom generators.
You only need to specify your type into which you want to marshal the json file

```go
import (
	"github.com/yandex/pandora/core"
	coreimport "github.com/yandex/pandora/core/import"
)

type MyCustomPayload struct {
	URL        string
	QueryParam string
}

//...
coreimport.RegisterCustomJSONProvider("my-custom-provider-name", func() core.Ammo { return &MyCustomPayload{} })
//...
```

And specify your ISP in the config

```yaml
provider:
  type: my-custom-provider-name
  ammo-queue-size: 1
  limit: 0
  passes: 0
  source:
    type: file
    path: my-costom-payload.json
```

You can also use `stdin`, `inline` for the source. See [payload sources](data-sources.md) for more details