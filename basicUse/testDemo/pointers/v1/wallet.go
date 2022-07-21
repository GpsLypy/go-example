package main

import (
	"fmt"
)

//iht不具有描述性，因此换个名字
type Bitcoin int

type Wallet struct {
	balance Bitcoin //余额，包级私有只在main包可见
}

func (w *Wallet) Deposit(amount Bitcoin) {
	fmt.Println("address of balance in Deposit is", &w.balance)
	w.balance += amount
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}
