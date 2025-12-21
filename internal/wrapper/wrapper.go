package wrapper

type WrapperInterface interface {
	Generate(opts Options) error
}

type Options interface {
	ValidateOpts() error
}

type GeneralOptions struct {
	Version  string
	Imports  bool
	CreateTx bool
}

type APMTypeWrapperOptions struct {
	GeneralOptions
}

func (o APMTypeWrapperOptions) ValidateOpts() error {
	return nil
}
