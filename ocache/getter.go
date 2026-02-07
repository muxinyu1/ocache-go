package ocache

type Getter interface {
	Get(group string, key string) (Value, error)
}
