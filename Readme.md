# safego 使用样例
## 无参数调用
 ```
func f1() {
	fmt.Println("f1")
	panic("f1")
}
ctx := context.Background()
g := NewPanicGroup()
g.Go(f1)
err := g.Wait(ctx)
if err != nil {
	fmt.Println(err)
}
 ```
## 带参数方法调用
 ```
// 闭包方式传入参数
func f1WithParams(a int, b string) func() {
	return func() {
		fmt.Println(a, b)
		fmt.Println("f1WithParams")
		panic("f1WithParams")
	}
}
ctx := context.Background()
g := NewPanicGroup()
g.Go(f1WithParams(123, "abc"))
err := g.Wait(ctx)
 ```
## 链式调用
 ```
ctx := context.Background()
err := NewPanicGroup().Go(f1).Go(f2).Wait(ctx)
 ```
## 注册自定义Catch处理器
 ```
// 默认Panic处理器
type PanicHandler struct {
}
// 默认Catch方法
func (*PanicHandler) Catch(Recover interface{}, Stack []byte) {
	fmt.Println(Recover, string(Stack))
}
g.RegPanicHandler(PanicHandler)
 ```
