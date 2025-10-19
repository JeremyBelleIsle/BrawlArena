package assets

import (
	_ "embed"
)

//go:embed pixel-driftcopie.mp3
var MusiqueBytes []byte

// GetMusiqueBytes returns the music file bytes
func GetMusiqueBytes() []byte {
	return MusiqueBytes
}
