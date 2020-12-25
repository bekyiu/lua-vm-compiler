package state

func (this *luaState) PushNil() {
	this.stack.push(nil)
}

func (this *luaState) PushBoolean(b bool) {
	this.stack.push(b)
}
func (this *luaState) PushInteger(n int64) {
	this.stack.push(n)
}
func (this *luaState) PushNumber(n float64) {
	this.stack.push(n)
}
func (this *luaState) PushString(s string) {
	this.stack.push(s)
}