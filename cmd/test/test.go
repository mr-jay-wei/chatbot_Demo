package main

import (
	"fmt"
	"sync"
	"time"
)

// 1. 纯粹的业务函数
// 你看，这里不再需要 wg 参数了，它只关心下载逻辑
func download(file string) {
	fmt.Printf("开始下载 %s...\n", file)
	time.Sleep(2 * time.Second) // 模拟耗时操作
	fmt.Printf("%s 下载完成！\n", file)
}

func main() {
	start := time.Now()

	var wg sync.WaitGroup
	files := []string{"file1.zip", "file2.zip", "file3.zip"}

	for _, file := range files {
		wg.Add(1)

		go func(f string) {
			defer wg.Done()
			download(f)
		}(file)
	}

	wg.Wait()

	fmt.Printf("总耗时: %v\n", time.Since(start))
}
