package platforms

import (
	stargen "github.com/stardustapp/dustgo/lib/utils/stargen/common"
)

func CountPlatforms() int {
  return len(stargen.Platforms)
}
