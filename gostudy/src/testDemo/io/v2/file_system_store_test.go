package main

import (
	"io"
	"io/ioutil"
	"os"
	_ "strings"
	"testing"
)

func TestFileSystemStore(t *testing.T) {

	t.Run("/league from a reader", func(t *testing.T) {
		//NewReader创建一个从s读取数据的Reader,语法：不转议
		// database := strings.NewReader(`[
		// 	{"Name": "Cleo", "Wins": 10},
		// 	{"Name": "Chris", "Wins": 33}]`)
		//为每个测试创建一个临时文件。*os.File 实现 ReadWriteSeeker
		database, cleanDatabase := createTempFile(t, `[
			{"Name":"Cleo","Wins":10},
			{"Name":"Chris","Wins":33}]`)

		defer cleanDatabase()
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

	t.Run("get player score", func(t *testing.T) {
		// database := strings.NewReader(`[
		// 	{"Name": "Cleo", "Wins": 10},
		// 	{"Name": "Chris", "Wins": 33}]`)
		database, cleanDatabase := createTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)

		defer cleanDatabase()

		store := FileSystemPlayerStore{database}

		got := store.GetPlayerScore("Chris")

		want := 33

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})
	t.Run("store wins for existing players", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
		defer cleanDatabase()
		store := FileSystemPlayerStore{database}
		store.RecordWin("Chris")
		got := store.GetPlayerScore("Chris")
		want := 34
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})
}

func createTempFile(t *testing.T, initialData string) (io.ReadWriteSeeker, func()) {
	t.Helper()
	//​ 创建一个临时文件供我们使用
	tmpfile, err := ioutil.TempFile("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}
	tmpfile.Write([]byte(initialData))
	removeFile := func() {
		os.Remove(tmpfile.Name())
	}
	return tmpfile, removeFile
}
