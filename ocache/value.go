package ocache

type Value interface {
	Len() int
	AsBytes() []byte
}
