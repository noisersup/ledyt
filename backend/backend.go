package backend

import "github.com/noisersup/ledyt/backend/common"

type Backend interface {
	Search(query string) ([]common.Video, error)
}
