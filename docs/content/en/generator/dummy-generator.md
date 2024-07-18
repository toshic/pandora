---
title: Dummy generator
description: Dummy generator is used to check load profiles
categories: [Provider]
tags: [provider, dummy]
weight: 1
---

When you are creating a complex load profile that consists of multiple phases, or just want to test
what your load profile will look like, you can use the Dummy Generator to keep your test server unloaded in the first stage.

```yaml
gun:
  type: dummy
  sleep: 10ms # optional
```