// File: internal/service/chat_service_test.go

package service

import (
	"testing"
)

// TestEchoService_GetResponse 是测试函数的标准命名格式：Test<结构体名>_<方法名>
func TestEchoService_GetResponse(t *testing.T) {
	// 1. 定义测试用例表格 (Table-Driven Tests)
	// 这是一个匿名结构体切片，用来存放多组测试数据
	tests := []struct {
		name     string // 测试用例的名称
		input    string // 输入的消息
		expected string // 期望得到的回复
	}{
		{
			name:     "Normal message",
			input:    "Hello World",
			expected: "You said: Hello World",
		},
		{
			name:     "Empty message",
			input:    "",
			expected: "You said: ",
		},
		{
			name:     "Special characters",
			input:    "!@#$%",
			expected: "You said: !@#$%",
		},
	}

	// 2. 初始化被测试的服务
	svc := NewChatService()

	// 3. 遍历表格执行测试
	for _, tt := range tests {
		// t.Run 启动一个子测试，这样如果某个用例挂了，我们可以清楚地知道是哪一个
		t.Run(tt.name, func(t *testing.T) {
			// 执行实际逻辑
			got := svc.GetResponse(tt.input)

			// 断言：对比实际结果 (got) 和期望结果 (expected)
			if got != tt.expected {
				t.Errorf("GetResponse() = %v, want %v", got, tt.expected)
			}
		})
	}
}