---
title: Источники Payload
description: Для упрощения реализации ваших кастомных генераторов
categories: [Provider]
tags: [provider, dummy]
weight: 40
---

Используется для json провайдеров

Есть 3 источника 

1. `file`

```yaml
source:
  type: file
  path: you_path
```

2. `stdin`

```yaml
source:
  type: stdin
```

3. `inline`

```yaml
source:
  type: inline
  data: |
    {"you": "json"}
```
