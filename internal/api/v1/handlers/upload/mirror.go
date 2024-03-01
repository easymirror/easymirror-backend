package upload

import "github.com/labstack/echo/v4"

// Mirror handles incoming POST requests for mirroring sites.
func (h *Handler) Mirror(c echo.Context) error {
	// TODO: Get user data form the JWT token
	// TODO Parse which sites to mirror to
	// TODO: Begin mirroring process
	return nil
}
