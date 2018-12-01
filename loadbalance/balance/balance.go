package balance

type Balance interface {
	DoBalance([]*Instance, ...string) (*Instance, error)
}
