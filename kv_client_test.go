package consul_kv_client

import (
	"testing"
	"time"
	"log"
	"os"
	"github.com/lattecake/consul-kv-client/api"
)

func TestKvClient_Get(t *testing.T) {

	config := &api.Config{
		Address:  os.Getenv("CONSUL_HTTP_ADDR"),
		Scheme:   "http",
		Token:    os.Getenv("CONSUL_HTTP_TOKEN"),
		WaitTime: time.Second * 30,
	}

	kv, err := NewKvClient(config)
	if err != nil {
		t.Fatal("new kv client error: ", err)
	}

	var prefix = "hello"
	var key = "world"

	if err = kv.Start(prefix); err != nil {
		t.Fatalf(err.Error())
	}

	val, err := kv.Get(key)
	if err != nil {
		t.Fatal("kv get val error: ", err)
	}

	log.Println("fast get val ", val)

	go func() {
		if err = kv.Put(key, "6666666"); err != nil {
			t.Fatalf("kv put val error: %#v", err)
		}
	}()

	time.Sleep(time.Second * 50)
	val, err = kv.Get(key)
	if err != nil {
		t.Fatal("kv get val error: ", err)
	}

	log.Println("val: ", val)
}
