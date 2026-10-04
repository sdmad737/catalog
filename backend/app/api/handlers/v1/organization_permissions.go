package v1

import (
	"errors"
	"net/http"

	"github.com/sdmad737/catalog/backend/internal/core/services"
	"github.com/sdmad737/catalog/backend/internal/sys/validate"
)

func requireOrganizationRole(r *http.Request, allowed ...string) error {
	actor := services.UseUserCtx(r.Context())
	if actor == nil || actor.Status == "disabled" {
		return validate.NewRequestError(errors.New("account is disabled"), http.StatusForbidden)
	}
	for _, role := range allowed {
		if actor.Role == role {
			return nil
		}
	}
	return validate.NewRequestError(errors.New("you do not have permission to perform this action"), http.StatusForbidden)
}
