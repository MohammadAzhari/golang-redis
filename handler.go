package main

import "sync"

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

var handlers = map[string]func(v []Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

var SETs = map[string]string{}
var setMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "error: SET command takes 2 args"}
	}

	key := args[0].bulk
	val := args[1].bulk

	setMu.Lock()
	SETs[key] = val
	setMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "error: GET command takes 1 args"}
	}

	key := args[0].bulk
	setMu.RLock()
	val, ok := SETs[key]
	setMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "string", str: val}
}

var HSETs = map[string]map[string]string{}
var hsetMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "error: GET command takes 3 args"}
	}

	collection := args[0].bulk
	key := args[1].bulk
	val := args[2].bulk

	hsetMu.Lock()
	if _, ok := HSETs[collection]; !ok {
		HSETs[collection] = map[string]string{}
	}
	HSETs[collection][key] = val
	hsetMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "error: GET command takes 2 args"}
	}

	collection := args[0].bulk
	key := args[1].bulk

	hsetMu.RLock()
	val, ok := HSETs[collection][key]
	hsetMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "string", str: val}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "error: GET command takes 1 arg"}
	}

	collection := args[0].bulk

	hsetMu.RLock()
	col, ok := HSETs[collection]
	hsetMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	array := make([]Value, 0)
	for _, val := range col {
		array = append(array, Value{typ: "string", str: val})
	}

	return Value{typ: "array", array: array}
}
