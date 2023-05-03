package routes

type Router interface {
	Route() error
}

type Closer interface {
	Close() error
}