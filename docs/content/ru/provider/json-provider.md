---
title: JSON провайдер
description: JSON провайдер преобразует в Payload json файлы
categories: [Provider]
tags: [provider, json]
weight: 20
---

Данные провайдер читает payload файлы и передает в генераторы структуру нужного типа.

Так как в дефолтной сборке нельзя использовать свои собственные структуры дефолтный провайдер возвращает
`map[string]interface{}{}`

Но данный провайдер удобен для создания пользовательских генераторов. 
Вам достаточно только указать свой тип, в который необходимо маршалить json файл

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

И в конфиге указать свой провайдер

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

Для источника можно так же использовать `stdin`, `inline`. Подробнее смотрите [источники payload](data-sources.md)