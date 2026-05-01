package engine

import (
	lua "github.com/yuin/gopher-lua"
)

// LuaSandbox 管理安全的 LUA 执行环境
type LuaSandbox struct {
	L *lua.LState
}

// NewLuaSandbox 创建新的 LUA 沙箱，移除危险包
func NewLuaSandbox() *LuaSandbox {
	s := &LuaSandbox{L: lua.NewState()}
	s.applySandbox()
	return s
}

// applySandbox 移除可能被恶意使用的标准库
func (s *LuaSandbox) applySandbox() {
	for _, name := range []string{"os", "io", "debug", "coroutine"} {
		s.L.SetGlobal(name, lua.LNil)
	}
}

// Close 关闭 LUA 状态，释放资源
func (s *LuaSandbox) Close() {
	s.L.Close()
}
