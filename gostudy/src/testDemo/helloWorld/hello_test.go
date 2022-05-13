package main

import "testing"

// t:hook（钩子）
// func TestHello(t *testing.T) {
// 	got := Hello()
// 	want := "Hello,world"
// 	if got != want {
// 		t.Errorf("got '%q' but want '%q'", got, want)
// 	}
// }

// func TestHello(t *testing.T) {
// 	got := Hello("Chris")
// 	want := "Hello,Chris"
// 	if got != want {
// 		t.Errorf("got '%q' but want '%q'", got, want)
// 	}
// }

// func TestHello(t *testing.T) {

// 	t.Run("saying hello to people", func(t *testing.T) {
// 		got := Hello("Chris")
// 		want := "Hello,Chris"

// 		if got != want {
// 			t.Errorf("got '%q' want '%q'", got, want)
// 		}
// 	})

// 	t.Run("say hello world when an empty string is supplied", func(t *testing.T) {
// 		got := Hello("")
// 		want := "Hello, World"

// 		if got != want {
// 			t.Errorf("got '%q' want '%q'", got, want)
// 		}
// 	})

// }

//子测试
func TestHello(t *testing.T) {

	assertCorrectMessage := func(t *testing.T, got, want string) {
		//t.Helper() 需要告诉测试套件这个方法是辅助函数（helper）通过这样做，当测试失败时所报告的行号将在函数调用中而不是在辅助函数内部
		t.Helper()
		if got != want {
			t.Errorf("got '%q' want '%q'", got, want)
		}
	}

	t.Run("saying hello to people", func(t *testing.T) {
		got := Hello("Chris", "")
		want := "Hello, Chris"
		assertCorrectMessage(t, got, want)
	})

	t.Run("empty string defaults to 'world'", func(t *testing.T) {
		got := Hello("", "")
		want := "Hello, World"
		assertCorrectMessage(t, got, want)
	})

	t.Run("in Spanish", func(t *testing.T) {
		got := Hello("Elodie", "Spanish")
		want := "Hola, Elodie"
		assertCorrectMessage(t, got, want)
	})

}
