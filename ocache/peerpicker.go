package ocache

type PeerPicker interface {
	PickPeer(key string) Client
}
