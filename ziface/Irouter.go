package ziface

type IRouter interface {
	// pre
	BeforeHandler(request IRequest)
	// handle
	Handler(request IRequest)
	// after
	AfterHandler(request IRequest)
}
