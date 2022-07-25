package main

import (
	"io"
)

// FileSystemPlayerStore stores players in the filesystem.
type FileSystemPlayerStore struct {
	database io.ReadSeeker
}

// GetLeague returns the scores of all the players.
func (f *FileSystemPlayerStore) GetLeague() []Player {
	// offset是偏移的位置，whence是偏移起始位置
	f.database.Seek(0, 0)
	league, _ := NewLeague(f.database)
	return league
}
