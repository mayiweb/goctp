/**
 * 安全 Map 防止并发读写
 */

package safe

import (
    "sync"
)

// 一级 Map
type Map struct {
    Lock sync.RWMutex
    Map map[interface{}]interface{}
}

// 二级 Map
type Map2 struct {
    Lock sync.RWMutex
    Map map[interface{}]map[interface{}]interface{}
}

// 设置 map
func (p *Map) Set(name interface{}, value interface{}) {

    p.Lock.Lock()
    defer p.Lock.Unlock()

    if len(p.Map) == 0 {
        p.Map = make(map[interface{}]interface{}, 0)
    }

    p.Map[name] = value
}

// 获取 map
func (p *Map) Get(name interface{}) (interface{}, bool) {

    p.Lock.RLock()
    defer p.Lock.RUnlock()

    mVal, mOk := p.Map[name]
    if !mOk {
        return nil, false
    }

    return mVal, true
}

// 删除 map
func (p *Map) Del(name interface{}) {

    p.Lock.Lock()
    defer p.Lock.Unlock()

    delete(p.Map, name)
}

// 清空 map
func (p *Map) Clear() {

    p.Lock.Lock()
    defer p.Lock.Unlock()

    p.Map = make(map[interface{}]interface{}, 0)
}

// 数据条数
func (p *Map) Size() int {
    p.Lock.RLock()
    defer p.Lock.RUnlock()

    return len(p.Map)
}

// 获取全部 map
func (p *Map) GetAll() map[interface{}]interface{} {
    p.Lock.RLock()
    defer p.Lock.RUnlock()

    // 申明一个新的 map 将 p.map 数据写入新 map 中，再访问数据就不是同一份数据了
    Result := make(map[interface{}]interface{}, 0)

    for k, v := range p.Map {
        Result[k] = v
    }

    return Result
}

// ------------------- 二级 map -------------------

// 设置 map
func (p *Map2) Set(key1 interface{}, key2 interface{}, value interface{}) {

    p.Lock.Lock()
    defer p.Lock.Unlock()

    if len(p.Map) == 0 {
        p.Map = make(map[interface{}]map[interface{}]interface{}, 0)
    }

    if len(p.Map[key1]) == 0 {
        p.Map[key1] = make(map[interface{}]interface{}, 0)
    }

    p.Map[key1][key2] = value
}

// 获取 map
func (p *Map2) Get(key1 interface{}, key2 interface{}) (interface{}, bool) {

    p.Lock.RLock()
    defer p.Lock.RUnlock()

    mVal, mOk := p.Map[key1][key2]
    if !mOk {
        return nil, false
    }

    return mVal, true
}

// 获取二级 key 列表
func (p *Map2) GetList(key1 interface{}) (map[interface{}]interface{}, bool) {
    p.Lock.RLock()
    defer p.Lock.RUnlock()

    mLists, mOk := p.Map[key1]
    if !mOk {
        return nil, false
    }

    // 申明一个新的 map 将 p.map 数据写入新 map 中，再访问数据就不是同一份数据了
    Result := make(map[interface{}]interface{}, 0)

    for k, v := range mLists {
        Result[k] = v
    }

    return Result, true
}

// 获取二级 key 列数据条数
func (p *Map2) GetListSize(key1 interface{}) int {
    p.Lock.RLock()
    defer p.Lock.RUnlock()

    _, mOk := p.Map[key1]
    if !mOk {
        return 0
    }

    return len(p.Map[key1])
}

// 删除一级 map
func (p *Map2) Del(key1 interface{}) {

    p.Lock.Lock()
    defer p.Lock.Unlock()

    delete(p.Map, key1)
}

// 删除二级 map
func (p *Map2) DelList(key1 interface{}, key2 interface{}) {

    p.Lock.Lock()
    defer p.Lock.Unlock()

    delete(p.Map[key1], key2)
}

// 清空 map
func (p *Map2) Clear() {

    p.Lock.Lock()
    defer p.Lock.Unlock()

    p.Map = make(map[interface{}]map[interface{}]interface{}, 0)
}

// 数据条数
func (p *Map2) Size() int {
    p.Lock.RLock()
    defer p.Lock.RUnlock()

    return len(p.Map)
}

// 获取全部 map
func (p *Map2) GetAll() map[interface{}]map[interface{}]interface{} {
    p.Lock.RLock()
    defer p.Lock.RUnlock()

    // 申明一个新的 map 将 p.map 数据写入新 map 中，再访问数据就不是同一份数据了
    Result := make(map[interface{}]map[interface{}]interface{}, 0)

    for k, v := range p.Map {

        ResultList := make(map[interface{}]interface{}, 0)

        for sk, sv := range v {
            ResultList[sk] = sv
        }

        Result[k] = ResultList
    }

    return Result
}