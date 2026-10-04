package v1

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hay-kot/httpkit/errchain"
	"github.com/sdmad737/catalog/backend/internal/core/services"
	"github.com/sdmad737/catalog/backend/internal/data/repo"
	"github.com/sdmad737/catalog/backend/internal/web/adapters"
)

func (ctrl *V1Controller) HandleItemMarkMissing() errchain.HandlerFunc {
	fn := func(r *http.Request, itemID uuid.UUID, body repo.MissingItemInput) (repo.ItemOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.ItemEvents.MarkMissing(r.Context(), auth.GID, itemID, body, auth.User.Name)
	}
	return adapters.ActionID("id", fn, http.StatusOK)
}

func (ctrl *V1Controller) HandleItemMarkFound() errchain.HandlerFunc {
	fn := func(r *http.Request, itemID uuid.UUID, body repo.FoundItemInput) (repo.ItemOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.ItemEvents.MarkFound(r.Context(), auth.GID, itemID, body, auth.User.Name)
	}
	return adapters.ActionID("id", fn, http.StatusOK)
}

func (ctrl *V1Controller) HandleItemHistory() errchain.HandlerFunc {
	fn := func(r *http.Request, itemID uuid.UUID) ([]repo.ItemEventOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.ItemEvents.List(r.Context(), auth.GID, itemID)
	}
	return adapters.CommandID("id", fn, http.StatusOK)
}
