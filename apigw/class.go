package apigw

// HandlerFunc describes a user-provided handler for API GW requests
type HandlerFunc func(Request, Params) (interface{}, error)

// Handler describes a wrapper for a user-provided function
type Handler interface {
	Check(Request) bool
	Run(Request, Params) (Response, error)
}

// HandlerSet is an ordered slice of HandlerShims
type HandlerSet []Handler

// DefaultSet is a set of parameter defaults
type DefaultSet map[string]string

// Router defines a dynamic handler
type Router struct {
	Handlers HandlerSet
	Defaults DefaultSet
}

// Handler returns the dynamic router
func (r *Router) Handler() func(req Request) (Response, error) {
	return func(req Request) (Response, error) {
		params := Params{
			Request:  &req,
			Defaults: r.Defaults,
		}
		for _, h := range r.Handlers {
			if h.Check(req) {
				return h.Run(req, params)
			}
		}
		return Fail("no handler found")
	}
}
