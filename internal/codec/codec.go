package codec

type Codec[T any] interface {
	Encode(T) (string, error)
	Decode(string) (T, error)
}
