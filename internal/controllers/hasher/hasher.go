package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

type HasherController struct {
}

func NewHasherController() *HasherController {
	return &HasherController{}
}

// HashFile computes the SHA-256 hash of a file's content
func (hc *HasherController) HashFile(content []byte) string {
	return hc.hash(content)
}

func (hc *HasherController) hash(content []byte) string {
	hasher := sha256.New()
	hasher.Write(content)
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
