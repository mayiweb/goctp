/**
 * ctp 队列
 */

package main

// ctp 队列
type CtpQueue struct{
    element []string
}

// 向队列中添加元素
func (q *CtpQueue) Push(str string) {
    q.element = append(q.element, str)
}

// 向队列中添加元素（防止重复）
func (q *CtpQueue) PushUnique(str string) {

    isPush := true
    for _, val := range q.element {
        if str == val {
            isPush = false
            break
        }
    }

    if isPush {
        q.element = append(q.element, str)
    }
}

// 移除队列第一个元素并返回
func (q *CtpQueue) Poll() interface{} {
    if q.IsEmpty() {
        return nil
    }

    // 取队列中第一个元素并返回
    FirstElement := q.element[0]

    // 移除队列中第一个元素
    q.element = q.element[1:]

    return FirstElement
}

// 清空队列
func (q *CtpQueue) Clear() bool {
    if q.IsEmpty() {
        return true
    }

    q.element = nil

    return true
}

// 队列是否为空
func (q *CtpQueue) IsEmpty() bool {
    if (q.Size() == 0) {
        return true
    }

    return false
}

// 返回队列中的元素数量
func (q *CtpQueue) Size() int {
    return len(q.element)
}