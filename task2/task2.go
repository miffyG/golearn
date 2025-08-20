package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	num := 5
	increaseValue(&num)
	fmt.Println("修改后的值:", num)

	slice := &[]int{1, 2, 3}
	doubleSliceValues(slice)
	fmt.Println("修改后的切片:", *slice)

	printNumbers()

	taskScheduler(func() {
		time.Sleep(1 * time.Second)
		fmt.Println("任务1完成")
	}, func() {
		time.Sleep(2 * time.Second)
		fmt.Println("任务2完成")
	}, func() {
		time.Sleep(3 * time.Second)
		fmt.Println("任务3完成")
	})

	rect := Rectangle{Width: 5.2, Height: 3.8}
	circle := Circle{Radius: 4.5}
	fmt.Println("矩形的面积:", rect.Area())
	fmt.Println("矩形的周长:", rect.Perimeter())
	fmt.Println("圆形的面积:", circle.Area())
	fmt.Println("圆形的周长:", circle.Perimeter())

	e := Employee{
		Person: Person{
			Name: "Alice",
			Age:  40,
		},
		EmployeeID: "E9527",
	}
	e.PrintInfo()

	chanCommunication()

	bufferedChan()

	mutexCounter()

	atomicCounter()
}

// 编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
func increaseValue(ptr *int) {
	*ptr += 10
}

// 实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
func doubleSliceValues(slice *[]int) {
	for i := range *slice {
		(*slice)[i] *= 2
	}
}

// 编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数
func printNumbers() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i += 2 {
			fmt.Println("奇数:", i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 2; i <= 10; i += 2 {
			fmt.Println("偶数:", i)
		}
	}()

	wg.Wait()
}

// 设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间
func taskScheduler(tasks ...func()) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(t func()) {
			defer wg.Done()
			start := time.Now()
			t()
			fmt.Println("任务执行时间:", time.Since(start))
		}(task)
	}
	wg.Wait()
}

// 定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法
type Shape interface {
	Area() float32
	Perimeter() float32
}

type Rectangle struct {
	Width  float32
	Height float32
}

type Circle struct {
	Radius float32
}

func (rect Rectangle) Area() float32 {
	return rect.Width * rect.Height
}

func (rect Rectangle) Perimeter() float32 {
	return (rect.Width + rect.Height) * 2
}

func (circle Circle) Area() float32 {
	return math.Pi * circle.Radius * circle.Radius
}

func (circle Circle) Perimeter() float32 {
	return 2 * math.Pi * circle.Radius
}

// 使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息
type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID string
}

func (e Employee) PrintInfo() {
	fmt.Printf("员工ID: %s\n", e.EmployeeID)
	fmt.Printf("姓名: %s\n", e.Name)
	fmt.Printf("年龄: %d\n", e.Age)
}

// 编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
func chanCommunication() {
	var wg sync.WaitGroup
	ch := make(chan int)
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for num := range ch {
			fmt.Println("接收到的数字:", num)
		}
	}()

	wg.Wait()
}

// 实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印
func bufferedChan() {
	var wg sync.WaitGroup
	ch := make(chan int, 10)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 100; i++ {
			ch <- i
		}
		close(ch)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for num := range ch {
			fmt.Println("接收到的数字:", num)
		}
	}()
	wg.Wait()
}

// 编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值
func mutexCounter() {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	var counter int

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mutex.Lock()
				counter++
				mutex.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Println("最终计数器的值:", counter)
}

// 使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值
func atomicCounter() {
	var wg sync.WaitGroup
	var counter atomic.Int32

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter.Add(1)
			}
		}()
	}
	wg.Wait()
	fmt.Println("最终计数器的值:", counter.Load())
}
