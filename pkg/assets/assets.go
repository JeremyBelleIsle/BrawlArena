package assets

import (
	_ "embed"
)

//go:embed pixel-driftcopie.mp3
var MusiqueBytes []byte

//go:embed 584131_Geometry-Dash-Menu-Theme.mp3
var MusiqueBytesm []byte

// GetMusiqueBytes returns the music file bytes
func GetMusiqueBytes() []byte {
	return MusiqueBytes
}
func GetMusiqueBytesm() []byte {
	return MusiqueBytesm
}
