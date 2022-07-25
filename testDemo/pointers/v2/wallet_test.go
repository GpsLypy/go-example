package main

import (
	_ "fmt"
	"testing"
)

func TestWallet(t *testing.T) {

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		wallet.Deposit(Bitcoin(10))
		assertBalance(t, wallet, Bitcoin(10))
	})

	t.Run("Withdraw", func(t *testing.T) {
		wallet := Wallet{balance: Bitcoin(20)}
		err := wallet.Withdraw(10)
		assertBalance(t, wallet, Bitcoin(10))
		assertNoError(t, err)
	})

	t.Run("Withdraw insufficient funds", func(t *testing.T) {
		startingBlance := Bitcoin(20)
		wallet := Wallet{startingBlance}
		err := wallet.Withdraw(100)
		assertBalance(t, wallet, startingBlance)
		assertError(t, err, InsufficientFundsError)
	})
}

func assertBalance(t *testing.T, wallet Wallet, want Bitcoin) {
	t.Helper()
	got := wallet.Balance()
	if got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
func assertError(t *testing.T, got error, want error) {
	if got == nil {
		//如果它被调用，它将停止测试。这是因为我们不希望对返回的错误进行更多断言。如果没有这个，
		//测试将继续进行下一步，并且因为一个空指针而引起 panic。
		t.Fatal("wanted an error but didnt get one")
	}
	if got != want {
		t.Errorf("got '%s',want '%s'", got, want)
	}
}

func assertNoError(t *testing.T, got error) {
	if got != nil {
		t.Fatal("got an error but didnt want one")
	}
}
