[Home](../index.md)

---

# Scenario generator / HTTP

- [Configuration](#configuration)
    - [Generator](#generator)
    - [Provider](#provider)
- [Description of the scenario format](#description-of-the-scenario-format)
    - [General principle](#general-principle)
    - [HCL example](#hcl-example)
    - [YAML example](#yaml-example)
- [Features](#features)
    - [Requests](#requests)
        - [Templater](#templater)
            - [Variable names in templates](#variable-names-in-templates)
        - [Preprocessors](#preprocessors)
        - [Postprocessors](#postprocessors)
            - [var/jsonpath](#varjsonpath)
            - [var/xpath](#varxpath)
            - [var/header](#varheader)
            - [assert/response](#assertresponse)
    - [Scenarios](#scenarios)
    - [Sources](#sources)
        - [csv file](#csv-file)
        - [json file](#json-file)
        - [variables](#variables)

## Configuration

You need to use a generator and a provider of type `http/scenario`

```yaml
pools:
  - id: Pool name
    gun:
      type: http/scenario
      target: localhost:80
    ammo:
      type: http/scenario
      file: payload.hcl
```

### Generator

The minimum generator configuration is as follows

```yaml
gun:
  type: http/scenario
  target: localhost:80
```

There is also a `type: http2/scenario` generator

```yaml
gun:
  type: http2/scenario
  target: localhost:80
```

All the settings of the regular [HTTP generator](http-generator.md) are supported for the scenario generator

### Provider

The provider accepts only one parameter - the path to the file with the scenario description

```yaml
ammo:
  type: http/scenario
  file: payload.hcl
```

Supports file extensions

- hcl
- yaml
- json

## Description of the scenario format

Supports formats

- hcl
- yaml
- json

### General principle

Several scenarios can be described in one file. A scenario has a name that distinguishes one scenario from another.

A scenario is a sequence of requests. That is, you will need to describe in the script which requests in what order
should be executed.

Request - HTTP request. Has the standard HTTP request fields plus additional fields. See [Requests](#requests).

### HCL example

```terraform
variable_source "source_name" "file/csv" {
  file              = "file.csv"
  fields            = ["id", "name"]
  ignore_first_line = true
  delimiter         = ","
}

request "request_name" {
  method  = "POST"
  uri     = "/uri"
  headers = {
    HeaderName = "header value"
  }
  tag       = "tag"
  body      = <<EOF
<body/>
EOF
  templater = "text"

  preprocessor {
    mapping = {
      new_var = "source.var_name[next].0"
    }
  }
  postprocessor "var/jsonpath" {
    mapping = {
      new_var = "$.auth_key"
    }
  }
}


scenario "scenario_name" {
  weight           = 1
  min_waiting_time = 1000
  requests         = [
    "request_name",
  ]
}
```

### YAML example

```yaml
variable_sources:
  - type: "file/csv"
    name: "source_name"
    ignore_first_line: true
    delimiter: ","
    file: "file.csv"
    fields: [ "id", "name" ]

requests:
  - name: "request_name"
    uri: '/uri'
    method: POST
    headers:
      Header-Name: "header value"
    tag: tag
    body: '<body/>'
    preprocessor:
      mapping:
        new_var: source.var_name[next].0
    templater:
      type: text
    postprocessors:
      - type: var/jsonpath
        mapping:
          token: "$.auth_key"

scenarios:
  - name: scenario_name
    weight: 1
    min_waiting_time: 1000
    requests: [
      request_name
    ]
```

## Features

### Requests

Поля

- method
- uri
- headers
- body
- **name**
- tag
- templater
- preprocessors
- postprocessors

### Templater

The `uri`, `headers`, `body` fields are templatized.

The standard go template is used.

#### Variable names in templates

Variable names have the full path of their definition.

For example

Variable `users` from source `user_file` - `{% raw %}{{.source.user_file.users}}{% endraw %}`

Variable `token` from the `list_req` query postprocessor - `{% raw %}{{.request.list_req.postprocessor.token}}{% endraw %}`

Variable `item` from the `list_req` query preprocessor - `{% raw %}{{.request.list_req.preprocessor.item}}{% endraw %}`

#### Preprocessors

Preprocessor - actions are performed before templating

It is used for creating new variable mapping

The preprocessor has the ability to work with arrays using modifiers

- next
- last
- rand

##### yaml

```yaml
requests:
  - name: req_name
    ...
    preprocessor:
      mapping:
        user_id: source.users[next].id
```

##### hcl

```terraform
request "req_name" {
  preprocessor {
    mapping = {
      user_id = "source.users[next].id"
    }
  }
}
```

#### Postprocessors

##### var/jsonpath

HCL example

```terraform
postprocessor "var/jsonpath" {
  mapping = {
    token = "$.auth_key"
  }
}
```

##### var/xpath

```terraform
postprocessor "var/xpath" {
  mapping = {
    data = "//div[@class='data']"
  }
}
```

##### var/header

Creates a new variable from response headers

It is possible to specify simple string manipulations via pipe

- lower
- upper
- substr(from, length)
- replace(search, replace)

```terraform
postprocessor "var/header" {
  mapping = {
    ContentType      = "Content-Type|upper"
    httpAuthorization = "Http-Authorization"
  }
}
```

##### assert/response

Checks header and body content

Upon assertion, further scenario execution is dropped

```terraform
postprocessor "assert/response" {
  headers = {
    "Content-Type" = "application/json"
  }
  body        = ["token"]
  status_code = 200

  size {
    val = 10000
    op  = ">"
  }
}
```

### Scenarios

The minimum fields for the script are name and list of requests

```terraform
scenario "scenario_name" {
  requests = [
    "list_req",
    "order_req",
    "order_req",
    "order_req"
  ]
}
```

You can specify a multiplicator for request repetition

```terraform
scenario "scenario_name" {
  requests = [
    "list_req",
    "order_req(3)"
  ]
}
```

You can specify the sleep() delay. Parameter in milliseconds

```terraform
scenario "scenario_name" {
  requests = [
    "list_req",
    "sleep(100)",
    "order_req(3)"
  ]
}
```

The second argument to request is **sleep** for requests with multipliers

```terraform
scenario "scenario_name" {
  requests = [
    "list_req",
    "sleep(100)",
    "order_req(3, 100)"
  ]
}
```

The `min_waiting_time` parameter describes the minimum scenario execution time. That is, a **sleep** will be added at the end of the entire
scenario if the scenario is executed faster than this parameter.

```terraform
scenario "scenario_name" {
  min_waiting_time = 1000
  requests         = [
    "list_req",
    "sleep(100)",
    "order_req(3, 100)"
  ]
}
```

Multiple scenarios can be described in one file.

The `weight` parameter is the distribution weight of each scenario. The greater the weight, the more often the scenario will be executed.


```terraform
scenario "scenario_first" {
  weight   = 1
  requests = [
    "auth_req(1, 100)",
    "list_req(1, 100)",
    "order_req(3, 100)"
  ]
}

scenario "scenario_second" {
  weight   = 50
  requests = [
    "mainpage",
  ]
}

```

### Sources

Variable sources

#### csv file

Example

```terraform
variable_source "users" "file/csv" {
  file              = "users.csv"                   # required
  fields            = ["user_id", "name", "pass"]   # optional
  ignore_first_line = true                          # optional
  delimiter         = ","                           # optional
}
```

Creating a source from csv. Adding the name `users` to it.

Using variables from this source

```gotempate
{% raw %}{{.source.users[0].user_id}}{% endraw %}
```

The `fields` parameter is optional.

If this parameter is not present, the names in the first line of the csv file will be used as field names,
if `ignore_first_line = false`.

If `ignore_first_line = true` and there are no fields, then ordinal numbers will be used as names

```gotempate
{% raw %}{{.source.users[0].0}}{% endraw %}
```

#### json file

Example

```terraform
variable_source "users" "file/json" {
  file = "users.json"     # required
}
```

Creating a source from a json file. Add the name `users` to it.

The file must contain any valid json. For example:

```json
{
    "data": [
        {
            "id": 1,
            "name": "user1"
        },
        {
            "id": 2,
            "name": "user2"
        }
    ]
}
```

Using variables from this source

```gotempate
{% raw %}{{.source.users.data[next].id}}{% endraw %}
```

Или пример с массивом

```json
 [
    {
        "id": 1,
        "name": "user1"
    },
    {
        "id": 2,
        "name": "user2"
    }
]
```

Using variables from this source

```gotempate
{% raw %}{{.source.users[next].id}}{% endraw %}
```

#### variables

Пример

```terraform
variable_source "variables" "variables" {
  variables = {
    host = localhost
    port = 8090
  }
}
```

Creating a source with variables. Add the name `variables` to it.

Using variables from this source

```gotempate
{% raw %}{{.source.variables.host}}:{{.source.variables.port}}{% endraw %}
```

---

[Home](../index.md)
