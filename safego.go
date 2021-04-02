package safego

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync/atomic"
)

type PanicHandler interface {
	// 异常处理方法
	Catch(Recover interface{}, Stack []byte)
	// Recover recover信息
	// Stack   抛出的异常栈信息
}

// 默认Panic处理器
type DefaultPanicHandler struct {
}

// 默认Catch方法
func (*DefaultPanicHandler) Catch(Recover interface{}, Stack []byte) {
	fmt.Println(Recover, string(Stack))
}

// Panic结构
type Panic struct {
	Recover interface{}
	Stack   []byte
}

// PanicGroup
type PanicGroup struct {
	panics  chan Panic    // 协程panic通知信道
	dones   chan struct{} // 协程完成通知信道，struct{}不占用存储仅作为信号
	conNum  int32         // 协程并发数量
	handler PanicHandler  // 协程Panic处理器
}

// 工厂方法
func NewPanicGroup() *PanicGroup {
	return &PanicGroup{
		panics:  make(chan Panic, 8),
		dones:   make(chan struct{}, 8),
		handler: &DefaultPanicHandler{},
	}
}

// 工厂方法（带协程数量控制）
func NewPanicGroupWithLimit(limit int32) *PanicGroup {
	return &PanicGroup{
		panics:  make(chan Panic, limit),
		dones:   make(chan struct{}, limit),
		handler: &DefaultPanicHandler{},
	}
}

// 注册Panic处理器
func (g *PanicGroup) RegPanicHandler(handler PanicHandler) *PanicGroup {
	g.handler = handler
	return g // 方便链式调用
}

// 封装Go方法
func (g *PanicGroup) Go(f func()) *PanicGroup {
	// 并发数量加1
	atomic.AddInt32(&g.conNum, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				g.panics <- Panic{Recover: r, Stack: debug.Stack()}
				return
			}
			g.dones <- struct{}{} // 发送结束信号
		}()
		f()
	}()
	return g // 方便链式调用
}

// 封装Wait方法
func (g *PanicGroup) Wait(ctx context.Context) error {
	for {
		select {
		case <-g.dones: // 正常结束
			if atomic.AddInt32(&g.conNum, -1) == 0 {
				return nil
			}
		case p := <-g.panics: // 异常处理
			g.handler.Catch(p.Recover, p.Stack)
		case <-ctx.Done(): // 上下文信号
			return ctx.Err()
		}
	}
}
