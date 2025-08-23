package codec

type StringCodec struct{}

func (c StringCodec) Encode(v string) (string, error) {
	return v, nil
}

func (c StringCodec) Decode(s string) (string, error) {
	return s, nil
}
