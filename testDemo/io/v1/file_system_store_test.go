package main

import (
	"strings"
	"testing"
)

func TestFileSystemStore(t *testing.T) {

	t.Run("/league from a reader", func(t *testing.T) {
		//NewReader创建一个从s读取数据的Reader,语法：不转议
		database := strings.NewReader(`[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		store := FileSystemPlayerStore{database}

		got := store.GetLeague()

		want := []Player{
			{"Cleo", 10},
			{"Chris", 33},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})
}
