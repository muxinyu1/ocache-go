package impl

type ByteView struct {
	data []byte
}

func (v *ByteView) Clone() *ByteView {
	newData := make([]byte, len(v.data))
	copy(newData, v.data)
	return &ByteView{newData}
}

func (v *ByteView) AsBytes() []byte {
	return v.data
}

func FromString(s string) *ByteView {
	return &ByteView{
		[]byte(s),
	}
}

func (v *ByteView) Len() int {
	return len(v.data)
}
