package main

import "testing"

//以下是一个表驱动测试的demo样子，我们可以使用-run 参数运行单个测试，例如，如果只想运行subtest1
//go test -run=TestFoo/subtest_1 -v

// func TestFoo(t *testing.T) {
// 	t.Run("subtest 1", func(t *testing.T) {
// 		if false {
// 			t.Error()
// 		}
// 	})

// 	t.Run("subtest 2", func(t *testing.T) {
// 		if 2 != 2 {
// 			t.Error()
// 		}
// 	})
// }

// 如何利用子测试来防止重复测试逻辑。实现思路是为每个案例点创建一个子测试，
// 定义一个map结构，map的键代表测试名称，map的值代表测试数据的输入值和预期值。实现代码如下：

func TestRemoveNewLineSuffixes(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		`empty`: {
			input: "",
			want:  "",
		},
		`ending with \r\n`: {
			input: "a\r\n",
			want:  "a",
		},
		`ending with \n`: {
			input: "a\n",
			want:  "a",
		},
		`ending with multiple \n`: {
			input: "a\n\n\n",
			want:  "a",
		},
		`ending without newline`: {
			input: "a",
			want:  "a",
		},
	}

	for name, tt := range tests {
		//避免闭包使用错误的tt变量值
		tt := tt
		t.Run(name, func(t *testing.T) {
			//标记测试必须并行运行
			t.Parallel()
			got := removeNewLineSuffixes(tt.input)
			if got != tt.want {
				t.Errorf("got: %s,want: %s", got, tt.want)
			}
		})
	}
}
