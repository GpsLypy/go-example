package main

import (
	"context"
	"time"
)

func slowOperationWithTimeout(ctx context.Context) (Result, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel() // 一旦慢操作完成就立马调用cancel
	return slowOperation(ctx)
}

func main() {

}
