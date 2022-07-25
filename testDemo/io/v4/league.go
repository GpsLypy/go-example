package main

import (
	"encoding/json"
	"io"
)

type League []Player

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}

// NewLeague creates a league from JSON.
func NewLeague(rdr io.Reader) (League, error) {
	var league League
	//json->[],解码
	err := json.NewDecoder(rdr).Decode(&league)
	return league, err
}
