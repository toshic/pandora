---
title: Производительность Pandora
description: Результаты тестов производительности
categories: [Best practices]
tags: [config, test, performance]
weight: 13
---

[Александр Иванов](mailto:ival.net@yandex.ru) провел несколько тестов производительности. Вот результаты.

* Server: NGinx, 32 cores, 64G RAM.
* Tank: 32 cores, 128G RAM.
* Network: 1G.

## HTTP запросы в nginx


Статические страницы разных размеров. Задержки сервера реализованы в скрипте Lua, мы можем
установить время задержки с помощью параметра запроса `sleep`:

```lua
server {
    listen          12999      default;
    listen          [::]:12999 default         ipv6only=on;
    server_name     pandora.test.yandex.net;

    location ~* / {

        rewrite_by_lua_block {
            local args = ngx.req.get_uri_args()
            if args['sleep'] then
                            ngx.sleep(args['sleep']/1000)
            end;
        }

        root /etc/nginx/pandora;
        error_page 404 = 404;

        }

        access_log off;
        error_log off;
}
```

* **Connection: Close** 23k RPS

![Connection:Close, response times distribution](../../images/http_connection_close_td.png)


* **Connection: Keep-Alive** 95k RPS

![Keep-Alive, response times distribution](../../images/http_keep_alive_td.png)


* **Response size 10kB** maxed out network interface. OK.
* **Response size 100kb** maxed out network interface. OK.
* **POST requests 10kB** maxed out network interface. OK.
* **POST requests 100kB** maxed out network interface. OK.
* **POST requests 1MB** maxed out network interface. OK.

![100 kb responses, network load](../../images/http_100kb_net.png)



* **50ms server delay** 30k RPS. OK.
* **500ms server delay** 30k RPS, 30k instances. OK.
* **1s server delay** 50k RPS, 50k instances. OK.
* **10s server delay** 5k RPS, 5k instances. OK.

**All good.**

![10s server delay, response times distribution](../../images/http_delay_10s_td.png)


![10s server delay, instances count](../../images/http_delay_10s_instances.png)



* **Server fail during test** OK.

![server fail emulation, response times quantiles](../../images/http_srv_fail_q.png)



## Пользовательские сценарии

Производительность пользовательских сценариев во многом зависит от их реализации. В некоторых наших
тестах мы наблюдали скачки, вызванные GC. Их можно избежать, уменьшив размер распределения.
Хорошей идеей будет оптимизировать ваши сценарии.
В Go есть множество <https://github.com/golang/go/wiki/Performance>`_ инструментов, помогающих вам
сделать это.

> Мы использовали payload в формате JSON для задания параметров для каждого сценария.

* **Small requests** 35k RPS. OK.
* **Some scenario steps with big JSON bodies** 35k RPS. OK.

![scenario steps](../../images/scn_cases.png)

