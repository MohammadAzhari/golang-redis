package main

import (
	"sync"

	"github.com/MohammadAzhari/golang-redis/resp"
)

var sets = map[string]string{}
var setMu = sync.RWMutex{}

var hsets = map[string]map[string]string{}
var hsetMu = sync.RWMutex{}

var handlers = map[string]func(v []resp.Value) resp.Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func ping(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.Value{Type: "string", String: "PONG"}
	}

	return resp.Value{Type: "string", String: args[0].String}
}

func set(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Type: "error", String: "error: SET command takes 2 args"}
	}

	key := args[0].String
	val := args[1].String

	setMu.Lock()
	sets[key] = val
	setMu.Unlock()

	return resp.Value{Type: "string", String: "OK"}
}

func get(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Type: "error", String: "error: GET command takes 1 args"}
	}

	key := args[0].String
	setMu.RLock()
	val, ok := sets[key]
	setMu.RUnlock()

	if !ok {
		return resp.Value{Type: "null"}
	}

	return resp.Value{Type: "string", String: val}
}

func hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{Type: "error", String: "error: GET command takes 3 args"}
	}

	collection := args[0].String
	key := args[1].String
	val := args[2].String

	hsetMu.Lock()
	if _, ok := hsets[collection]; !ok {
		hsets[collection] = map[string]string{}
	}
	hsets[collection][key] = val
	hsetMu.Unlock()

	return resp.Value{Type: "string", String: "OK"}
}

func hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Type: "error", String: "error: GET command takes 2 args"}
	}

	collection := args[0].String
	key := args[1].String

	hsetMu.RLock()
	val, ok := hsets[collection][key]
	hsetMu.RUnlock()

	if !ok {
		return resp.Value{Type: "null"}
	}

	return resp.Value{Type: "string", String: val}
}

func hgetall(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Type: "error", String: "error: GET command takes 1 arg"}
	}

	collection := args[0].String

	hsetMu.RLock()
	col, ok := hsets[collection]
	hsetMu.RUnlock()

	if !ok {
		return resp.Value{Type: "null"}
	}

	array := make([]resp.Value, 0)
	for _, val := range col {
		array = append(array, resp.Value{Type: "string", String: val})
	}

	return resp.Value{Type: "array", Array: array}
}
