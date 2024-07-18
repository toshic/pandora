---
title: Dummy Provider
description: Dummy provider is used for generators without external Payload.
categories: [Provider]
tags: [provider, dummy]
weight: 30
---

Pandora requires a Provider configuration. But there are cases where you don't need any Payload.
In such cases, you can use an empty `dummy` Provider

```yaml
ammo:
type: dummy
```

This provider is useful in 2 cases:

1. for checks on your payload profile configuration
2. for custom generators, where you prepare your Payloads inside the generator