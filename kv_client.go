package consul_kv_client

import (
	"errors"
	"fmt"
	"github.com/lattecake/consul-kv-client/api"
	"log"
	"strings"
	"sync"
)

type KvClient interface {
	Get(key string) (v string, err error)
	Start(key string) error
	Put(key string, val string) error
	Delete(key string) error
	EnableLog(log bool) KvClient
}

type kvClient struct {
	kv          *api.KV
	modifyIndex uint64
	//running     bool
	data    sync.Map
	watch   map[string]bool
	prefix  string
	logging bool
}

func NewKvClient(config *api.Config) (kv KvClient, err error) {
	client, err := api.NewClient(config)

	if err != nil {
		return
	}

	return &kvClient{kv: client.KV(), watch: map[string]bool{}}, nil

}

func (c *kvClient) EnableLog(log bool) KvClient {
	c.logging = log
	return c
}

func (c *kvClient) Get(key string) (v string, err error) {
	if !strings.HasPrefix(key, c.prefix) {
		key = c.prefix + key
	}
	var i = 1
Load:
	val, ok := c.data.Load(key)
	if !ok {
		if i < 1 {
			err = errors.New("key not exists! " + key)
			return
		}
		if err = c.get(key); err != nil {
			err = errors.New("key not exists! " + key)
			return
		}
		i--
		goto Load
	}

	if v, ok := val.([]byte); ok == true {
		return string(v), nil
	}

	return
}

func (c *kvClient) Put(key string, val string) error {
	if !strings.HasPrefix(key, c.prefix) {
		key = c.prefix + key
	}
	var p = &api.KVPair{
		Key:   key,
		Value: []byte(val),
	}
	if _, err := c.kv.Put(p, nil); err != nil {
		return err
	}
	c.loadData(key, []byte(val))
	return nil
}

func (c *kvClient) Delete(key string) error {
	if _, err := c.kv.Delete(key, nil); err != nil {
		return err
	}

	c.data.Delete(key)
	return nil
}

func (c *kvClient) Start(key string) error {
	if !strings.HasSuffix(key, "/") {
		key = key + "/"
	}
	c.prefix = key
	if c.watch[key] == true {
		return errors.New(fmt.Sprintf("Key %s already watch", key))
	}

	pairs, _, err := c.kv.List(key, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("kv list err %s ", err.Error()))
	}

	for _, val := range pairs {
		go func(k string) {
			for {
				if err := c.watchKey(k); err != nil {
					log.Fatalf("watch key %s err %s", k, err.Error())
					break
				}
			}
		}(val.Key)
		if err = c.get(val.Key); err != nil {
			log.Fatalf("watch key %s err %s", key, err.Error())
		}
	}

	c.watch[key] = true
	return nil
}

func (c *kvClient) get(key string) error {
	pair, _, err := c.kv.Get(key, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("kv get %s ", err.Error()))
	}

	if pair == nil {
		return errors.New(fmt.Sprintf("%s value is not found", key))
	}

	c.loadData(key, pair.Value)
	return nil
}

func (c *kvClient) watchKey(key string) (err error) {
	var prevIndex = c.modifyIndex
	pair, meta, err := c.kv.Get(key, &api.QueryOptions{
		WaitIndex: prevIndex,
	})
	if err != nil {
		err = errors.New(fmt.Sprintf("kv get err: %s", err.Error()))
		return
	}
	if meta.LastIndex == 0 {
		err = errors.New("meta last index is zero.")
		return
	}
	c.getModifyIndex(pair)
	if c.modifyIndex > prevIndex {
		c.printLog("modifyIndex: ", c.modifyIndex, " prevIndex: ", prevIndex, " val: ", string(pair.Value))
		c.loadData(key, pair.Value)
	}
	return
}

func (c *kvClient) printLog(keyVales ...interface{}) {
	if c.logging {
		log.Println(keyVales)
	}
}

func (c *kvClient) getModifyIndex(pair *api.KVPair) {
	var maxIndex = c.modifyIndex
	if pair.ModifyIndex > maxIndex {
		c.modifyIndex = pair.ModifyIndex
	}
}

func (c *kvClient) loadData(key string, value []byte) {
	c.data.Store(key, value)
}
