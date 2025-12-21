package test_samples

import "context"

type NamedParamsAndResults interface {
	Method1(a string, b *int, c []byte) (s string, err error)
}

type NoParams interface {
	Method1() error
	Method2() (s string, err error)
	Method3() (string, error)
}

type NoResult interface {
	Method1(s []string)
	Method2(a, b, c int, s context.Context)
	Method3(a, _, c int, _ *context.Context)
}
