package ocache

type Client interface {
	Get(groupName string, key string) (Value, error)
}
