package state

// 返回栈顶位置
func (this *luaState) GetTop() int {
	return this.stack.top
}

// 索引转换为绝对索引
func (this *luaState) AbsIndex(idx int) int {
	return this.stack.absIndex(idx)
}

// check若不足n个则扩容
func (this *luaState) CheckStack(n int) bool {
	this.stack.check(n)
	return true
}

// pop n个值
func (this *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		this.stack.pop()
	}
}

// 根据索引copy值
func (this *luaState) Copy(fromIdx, toIdx int) {
	e := this.stack.get(fromIdx)
	this.stack.set(toIdx, e)
}

// 把索引处的值复制一份到栈顶
func (this *luaState) PushValue(idx int) {
	e := this.stack.get(idx)
	this.stack.push(e)
}

// pop栈顶, 写到指定位置
func (this *luaState) Replace(idx int) {
	e := this.stack.pop()
	this.stack.set(idx, e)
}

// 如果i为正数, 则把[idx, top]之间的元素向栈顶平移n个位置, 越界的元素按顺序移动到空出来的位置
// 如果i为负数, 则向栈底平移
func (this *luaState) Rotate(idx int, n int) {
	t := this.stack.top - 1
	p := this.stack.absIndex(idx) - 1
	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	this.stack.reverse(p, m)
	this.stack.reverse(m + 1, t)
	this.stack.reverse(p, t)
}

// pop栈顶, 插入到指定位置
func (this *luaState) Insert(idx int) {
	this.Rotate(idx, 1)
}

// 删除指定位置的值
func (this *luaState) Remove(idx int) {
	this.Rotate(idx, -1)
	this.Pop(1)
}

// 设置栈顶 索引
func (this *luaState) SetTop(idx int) {
	newTop := this.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}
	n := this.stack.top - newTop
	if n > 0 {
		this.Pop(n)
	} else if n < 0 {
		for i := 0; i > n; i-- {
			this.stack.push(nil)
		}
	}
}