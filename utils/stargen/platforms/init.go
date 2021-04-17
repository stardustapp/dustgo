package platforms

import (
	stargen "github.com/stardustapp/dustgo/utils/stargen/common"
)

func CountPlatforms() int {
	return len(stargen.Platforms)
}
