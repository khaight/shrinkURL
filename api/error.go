package api

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse object for error responses
type errResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func errRender(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response",
		ErrorText:      err.Error(),
	}
}

func errNotFound(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Resource Not Found",
		ErrorText:      err.Error(),
	}
}

func errBadRequest(err error) render.Renderer {
	return &errResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid Request",
		ErrorText:      err.Error(),
	}
}

// Render response
func (e *errResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
