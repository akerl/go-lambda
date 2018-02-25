package apigw

// TextHandler does basic text responses
type TextHandler struct {
	Func HandlerFunc
}

// Check is always true for text handlers
func (h *TextHandler) Check(req Request) bool {
	return true
}

// Run calls the func with the provided request
func (h *TextHandler) Run(req Request, params Params) (Response, error) {
	resp, err := h.Func(req, params)
	if err != nil {
		return Fail(err.Error())
	}
	text, ok := resp.(string)
	if !ok {
		return Fail("handler returned unexpected type")
	}
	return Succeed(text)
}
