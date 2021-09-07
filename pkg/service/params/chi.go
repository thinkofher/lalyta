/*
Package params contains implementation for QueryParameters api
interface.
*/
package params

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Chi implements QueryParameters api interface for
// chi router.
//
// You can safely use new bulit-in function to allocate
// new Chi instace.
type Chi struct{}

func (c *Chi) ID(r *http.Request) string {
	return chi.URLParam(r, "id")
}
