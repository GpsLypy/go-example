package main

import (
	"encoding/json"
	"io"
)

// NewLeague creates a league from JSON.
func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	//json->[],解码
	err := json.NewDecoder(rdr).Decode(&league)
	return league, err
}
