package poker

import (
	"strings"
	"testing"
)

// func TestCLI(t *testing.T) {

// 	playerStore := &StubPlayerStore{}
// 	cli := &CLI{playerStore}
// 	cli.PlayPoker()
// 	if len(playerStore.winCalls) != 1 {
// 		t.Fatal("expected a win call but didn't get any")
// 	}
// }

// func TestCLI(t *testing.T) {
// 	in := strings.NewReader("Chris wins\n")
// 	playerStore := &StubPlayerStore{}
// 	cli := &CLI{playerStore, in}
// 	cli.PlayPoker()
// 	if len(playerStore.winCalls) < 1 {
// 		t.Fatal("expected a win call but didn't get any")
// 	}
// 	got := playerStore.winCalls[0]
// 	want := "Chris"
// 	if got != want {
// 		t.Errorf("didn't record correct winner, got '%s', want '%s'", got, want)
// 	}
// }

// func TestCLI(t *testing.T) {
// 	in := strings.NewReader("Chris wins\n")
// 	playerStore := &StubPlayerStore{}

// 	cli := &CLI{playerStore, in}
// 	cli.PlayPoker()

// 	assertPlayerWin(t, playerStore, "Chris")
// }

func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &StubPlayerStore{}

		//cli := &CLI{playerStore, in}
		cli := NewCLI(playerStore, in)
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &StubPlayerStore{}

		cli := NewCLI(playerStore, in)
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Cleo")
	})

}

func assertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got '%s' want '%s'", store.winCalls[0], winner)
	}
}
