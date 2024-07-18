---
title: Data Sources
description: To simplify the implementation of your custom generators
categories: [Provider]
tags: [provider, dummy]
weight: 40
---

Used for json providers

There are 3 data sources

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
