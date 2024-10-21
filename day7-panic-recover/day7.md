## day7
- 主动触发 panic
```go
func main() {
	fmt.Println("before panic")
	panic("crash")
	fmt.Println("after panic")
}
```
```bash
$ go run hello.go

before panic
panic: crash

goroutine 1 [running]:
main.main()
        ~/go_demo/hello/hello.go:7 +0x95
exit status 2
```
- 数组越界触发的 panic
```go
func main() {
	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
}
```
```bash
$ go run hello.go
panic: runtime error: index out of range [4] with length 3
```
- defer
  - panic 会导致程序被中止，但是在退出前，会先处理完当前协程上已经defer 的任务，执行完成后再退出。
```go
func main() {
	defer func() {
		fmt.Println("defer func")
	}()

	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
}
```
```bash
$ go run hello.go 
defer func
panic: runtime error: index out of range [4] with length 3
```
- recover
  - 避免因为 panic 发生而导致整个程序终止，recover 函数只在 defer 中生效
```go
// hello.go
func test_recover() {
	defer func() {
		fmt.Println("defer func")
		if err := recover(); err != nil {
			fmt.Println("recover success")
		}
	}()

	arr := []int{1, 2, 3}
	fmt.Println(arr[4])
	fmt.Println("after panic")
}

func main() {
	test_recover()
	fmt.Println("after recover")
}

```
```bash
$ go run hello.go 
defer func
recover success
after recover
```