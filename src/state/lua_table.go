package state

import (
	"math"
	"write_lua/src/number"
)

// lua表
type luaTable struct {
	// 当表被作为数组使用时, 存放在arr中
	arr []luaValue
	// 当表被作为map使用时, 存放在_map中
	_map map[luaValue]luaValue
	// 与当前表关联的元表
	metatable *luaTable
	// 用于键值遍历, 记录了当前键对应的下一个键
	keys map[luaValue]luaValue
	// table中最后一个key
	lastKey luaValue
	// table是否更新过
	changed bool
}

func newLuaTable(nArr, nRec int) *luaTable {
	t := &luaTable{}
	if nArr > 0 {
		t.arr = make([]luaValue, 0, nArr)
	}
	if nRec > 0 {
		t._map = make(map[luaValue]luaValue, nRec)
	}
	return t
}

// 尝试把key转为整数
func _floatToInteger(key luaValue) luaValue {
	if f, ok := key.(float64); ok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}

// 从arr, 或者map里根据key取值
func (this *luaTable) get(key luaValue) luaValue {
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(this.arr)) {
			return this.arr[idx-1]
		}
	}
	return this._map[key]
}

func (this *luaTable) put(key, val luaValue) {
	if key == nil {
		panic("table index is nil!")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN!")
	}

	this.changed = true
	key = _floatToInteger(key)
	// 存数组
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(this.arr))
		if idx <= arrLen {
			this.arr[idx-1] = val
			if idx == arrLen && val == nil {
				this._shrinkArray()
			}
			return
		}
		if idx == arrLen+1 {
			// 如果arrIdx很大, 就会存在map中
			// 扩容时, 先在map中删一次, 以免arr和map中存在相同的key
			delete(this._map, key)
			if val != nil {
				this.arr = append(this.arr, val)
				this._expandArray()
			}
			return
		}
	}
	// 不是数字key, 或者key超过数组的大小, 存入map
	if val != nil {
		if this._map == nil {
			this._map = make(map[luaValue]luaValue, 8)
		}
		this._map[key] = val
	} else {
		delete(this._map, key)
	}
}

// 去除arr末尾的洞
func (this *luaTable) _shrinkArray() {
	for i := len(this.arr) - 1; i >= 0; i-- {
		if this.arr[i] == nil {
			this.arr = this.arr[0:i]
		}
	}
}

// 把map中的数字key挪到数组中
func (this *luaTable) _expandArray() {
	for idx := int64(len(this.arr)) + 1; true; idx++ {
		if val, found := this._map[idx]; found {
			delete(this._map, idx)
			this.arr = append(this.arr, val)
		} else {
			break
		}
	}
}

func (this *luaTable) len() int {
	return len(this.arr)
}

// 是否有指定元方法
func (this *luaTable) hasMetafield(fieldName string) bool {
	return this.metatable != nil &&
		this.metatable.get(fieldName) != nil
}

func (this *luaTable) initKeys() {
	this.keys = make(map[luaValue]luaValue)
	// 当前key
	var key luaValue = nil
	for i, v := range this.arr {
		if v != nil {
			// 当前key对应的下一个key
			this.keys[key] = int64(i + 1)
			key = int64(i + 1)
		}
	}
	for k, v := range this._map{
		if v != nil {
			this.keys[key] = k
			key = k
		}
	}
	this.lastKey = key
}

// 给出当前key, 返回下一个key
func (this *luaTable) nextKey(key luaValue) luaValue {
	if this.keys == nil || (key == nil && this.changed) {
		this.initKeys()
		this.changed = false
	}
	nextKey := this.keys[key]
	if nextKey == nil && key != nil && key != this.lastKey {
		panic("invalid key to 'next'")
	}

	return nextKey
}