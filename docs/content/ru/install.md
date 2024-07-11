---
title: Установка
description: Как установить Yandex Pandora
categories: [Get started]
tags: [install]
weight: 1
---

[Скачать](https://github.com/yandex/pandora/releases) релиз или собрать из сходников.

Pandora использует **go modules**

```bash
git clone https://github.com/yandex/pandora.git
cd pandora
go mod download
```

Также возможна кросс-компиляция под другие arch/os:

```bash
GOOS=linux GOARCH=amd64 go build
```

