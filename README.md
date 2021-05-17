# consul kv 客户端


参考了 [https://github.com/hashicorp/consul/tree/master/api](https://github.com/hashicorp/consul/tree/master/api) 实现

去除了一些没必要的文件。

## 安装

```
$ go get github.com/lattecake/consul-kv-client
```

## 使用说明

```go
package main

import (
    consulKv "github.com/lattecake/consul-kv-client"
    "github.com/lattecake/consul-kv-client/api"
    "os"
)

var (
    kv *consulKv.KvClient
    var prefix = "hello"
    var key = "world"

)

func init() {
    config := &api.Config{
        Address:  os.Getenv("CONSUL_HTTP_ADDR"),
        Scheme:   "http",
        Token:    os.Getenv("CONSUL_HTTP_TOKEN"),
        WaitTime: time.Second * 30,
    }

    kv, err := consulKv.NewKvClient(config)
    if err != nil {
        log.Fatalf("new kv client error: %#v", err)
    }
    if err = kv.Start(prefix); err != nil {
        log.Fatalf(err)
    }
}

func main() {
    val, err := kv.Get(key)
    if err != nil {
        log.Fatalf("kv get val error: %#v", err)
    }
    log.Println(val)
}

```

将kv设置为

- `kv.Get(key)` 这个方法将从data直接获取不调用api 启用start将自动load进data

- `kv.Start(key)` 这个方法其实就是启动实时监听，当版本号一旦有变化，将实时更到要data里

> 如果想要同时监听多个键，可以执行多个 `kv.start()`

*详情请看go test*

## Testing

`go test`

![](https://ofbudvg4c.qnssl.com//images/2018/06/c9/22/a2/20180603-c4b1802c39542936a2333a83220c9716.jpeg)

![](https://ofbudvg4c.qnssl.com//images/2018/06/b8/23/1f/20180603-c09d8f1a999cfd4311b5f08d238b8a5d.jpeg)


