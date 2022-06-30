package poker

import (
	"encoding/json"
	"fmt"
	_ "io"
	"os"
	"sort"

	_ "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
)

// FileSystemPlayerStore stores players in the filesystem.
type FileSystemPlayerStore struct {
	//database io.Writer
	database *json.Encoder
	league   League
}

//func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
// func NewFileSystemPlayerStore(file *os.File) *FileSystemPlayerStore {
// 	file.Seek(0, 0)
// 	league, _ := NewLeague(file)
// 	return &FileSystemPlayerStore{
// 		//database: &tape{database},
// 		database: json.NewEncoder(&tape{file}),
// 		league:   league,
// 	}
// }

func FileSystemPlayerStoreFromFile(path string) (*FileSystemPlayerStore, error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, fmt.Errorf("problem opening %s %v", path, err)
	}

	store, err := NewFileSystemPlayerStore(db)

	if err != nil {
		return nil, fmt.Errorf("problem creating file system player store, %v ", err)
	}

	return store, nil
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	err := initialisePlayerDBFile(file)
	league, err := NewLeague(file)
	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		//database: &tape{file},
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}

func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, 0)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return nil
}

// GetLeague returns the scores of all the players.
func (f *FileSystemPlayerStore) GetLeague() League {
	// offset是偏移的位置，whence是偏移起始位置
	// f.database.Seek(0, 0)
	// league, _ := NewLeague(f.database)
	// return league
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Wins > f.league[j].Wins
	})
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	// var wins int
	// for _, player := range f.GetLeague() {
	// 	if player.Name == name {
	// 		wins = player.Wins
	// 		break
	// 	}
	// }
	// return wins
	player := f.league.Find(name)
	if player != nil {
		return player.Wins
	}
	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
	//league := f.GetLeague()
	// for i, player := range league {
	// 	if player.Name == name {
	// 		league[i].Wins++
	// 	}
	// }
	player := f.league.Find(name)
	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}
	//f.database.Seek(0, 0)
	//编码结果暂存到f.database,创建编码器
	f.database.Encode(f.league)
}
