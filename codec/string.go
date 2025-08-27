package codec

type StringCodec struct{}

func (StringCodec) Encode(v string) (string, error) {
	return v, nil
}

func (StringCodec) Decode(s string) (string, error) {
	return s, nil
}
