package backend

import "github.com/noisersup/ledyt/backend/common"

// Backend is the representation of the Youtube API Client.
// It implements all required functions needed for UI to interact
// with the API.
type Backend interface {
	Search(query string) ([]common.Video, error)
}
