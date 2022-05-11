package internal

type Transformer interface {
	Transform(args []string) error
}
