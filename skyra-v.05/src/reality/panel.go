package reality

type Panel interface {
	Reality
	Port(name string) Reality
}
