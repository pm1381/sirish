package test_samples

import (
	"github.com/pm1381/sirish/internal/visitors/test_samples/rand"
	rand2 "math/rand"
)

type RandConflict interface {
	Method1(t1 rand.Something, t2 rand2.Rand) rand2.Rand
}
