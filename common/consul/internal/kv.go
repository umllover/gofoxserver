package internal

import (
	"strings"
	"time"

	"strconv"

	"github.com/lovelly/leaf/log"
	"github.com/hashicorp/consul/api"
)

type kvValue map[string]int

//定时获取最新负载配置
func watchKV(client *api.Client, path string) {
	var lastIndex uint64

	for {
		value, index, err := getKVAll(client, path, lastIndex)
		if err != nil {
			log.Debug("consul: Error fetching config from %s. %v", path, err)
			time.Sleep(time.Second * UpdateConfigTicke)
			continue
		}

		if index != lastIndex && len(value) > 0{
			log.Debug("consul: Manual kvpath changed to #%d", index)
			ChanRPC.Go("KvUpdate",  value)
			lastIndex = index
		} else {
			time.Sleep(time.Second * UpdateConfigTicke)
		}
	}
}

//获取扣个key下所有值
func getKVAll(client *api.Client, key string, waitIndex uint64) (kvValue, uint64, error) {
	q := &api.QueryOptions{RequireConsistent: true, WaitIndex: waitIndex}
	kvpairs, meta, err := client.KV().List(key, q)
	values := make(kvValue, 0)
	if err != nil {
		return nil, 0, err
	}
	if kvpairs == nil {
		return nil, meta.LastIndex, nil
	}

	for _, v := range kvpairs {
		keys := strings.Split(v.Key, "/")
		klen := len(keys)
		if klen <= 0 {
			log.Debug("at getKVAll Split key error: key len <= 0")
			continue
		}

		strvalue := strings.TrimSpace(string(v.Value))
		IntVal, err1 := strconv.Atoi(strvalue)
		if err1 != nil {
			log.Error("at getKVAll Split key error: v.Value.(int) faild, ERROR:%s", err1.Error())
			continue
		}
		values[keys[klen-1]] = IntVal
	}
	return values, meta.LastIndex, nil
}

//取值
func getKV(client *api.Client, key string, waitIndex uint64) (kvValue, uint64, error) {
	q := &api.QueryOptions{RequireConsistent: true, WaitIndex: waitIndex}
	kvpair, meta, err := client.KV().Get(key, q)
	values := make(kvValue, 0)

	if err != nil {
		return nil, 0, err
	}
	if kvpair == nil {
		return nil, meta.LastIndex, nil
	}

	keys := strings.Split(kvpair.Key, "/")
	klen := len(keys)
	if klen <= 0 {
		log.Debug("at getKVAll Split key error: key len <= 0")
		return nil, meta.LastIndex, nil
	}

	strvalue := strings.TrimSpace(string(kvpair.Value))
	IntVal, err1 := strconv.Atoi(strvalue)
	if err1 != nil {
		log.Error("at getKVAll Split key error: v.Value.(int) faild, error:%s",  err1.Error())
		return nil, meta.LastIndex, nil
	}
	values[keys[klen-1]] = IntVal
	return values, meta.LastIndex, nil
}

//存值
func putKV(client *api.Client, key, value string) (bool, error) {
	p := &api.KVPair{Key: key, Value: []byte(value)}
	_, err := client.KV().Put(p, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

//CAS存值
func CASputKV(client *api.Client, key, value string, index uint64) (bool, error) {
	p := &api.KVPair{Key: key, Value: []byte(value), ModifyIndex: index}
	ok, _, err := client.KV().CAS(p, nil)
	if err != nil {
		return false, err
	}
	return ok, nil
}

//删除kv
func DelKV(client *api.Client, key string) {
	log.Debug("delete kv id is ", key)
	client.KV().Delete(key, nil)
}
