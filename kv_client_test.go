package consul_kv_client

import (
	"testing"
	"time"
	"log"
	"os"
	"github.com/lattecake/consul-kv-client/api"
)

var (
	consul_url = "http://127.0.0.1:8500"
)

func TestKvClient_Get(t *testing.T) {

	consulKey := os.Getenv("POD_CONSUL_KEY")
	if consulKey == "" {
		consulKey = "daea5ec6-ebde-3a8e-33f2-d223c2fa0bb9"
	} else {
		consul_url = "http://consul:8500"
	}

	config := &api.Config{
		Address:  consul_url,
		Scheme:   "http",
		Token:    consulKey,
		WaitTime: time.Second * 5,
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
