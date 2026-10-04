package v1

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/hay-kot/httpkit/errchain"
	"github.com/sdmad737/catalog/backend/internal/core/services"
	"github.com/sdmad737/catalog/backend/internal/data/repo"
	"github.com/sdmad737/catalog/backend/internal/sys/validate"
	"github.com/sdmad737/catalog/backend/internal/web/adapters"
)

type PeopleQuery struct {
	Search string `json:"search" schema:"q"`
}

func (ctrl *V1Controller) HandlePeopleList() errchain.HandlerFunc {
	fn := func(r *http.Request, query PeopleQuery) ([]repo.PersonOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.People.List(r.Context(), auth.GID, query.Search)
	}
	return adapters.Query(fn, http.StatusOK)
}

func (ctrl *V1Controller) HandlePersonCreate() errchain.HandlerFunc {
	fn := func(r *http.Request, input repo.PersonCreate) (repo.PersonOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.People.Create(r.Context(), auth.GID, input)
	}
	return adapters.Action(fn, http.StatusCreated)
}

func (ctrl *V1Controller) HandlePersonUpdate() errchain.HandlerFunc {
	fn := func(r *http.Request, id uuid.UUID, input repo.PersonUpdate) (repo.PersonOut, error) {
		auth := services.NewContext(r.Context())
		return ctrl.repo.People.Update(r.Context(), auth.GID, id, input)
	}
	return adapters.ActionID("id", fn, http.StatusOK)
}

func (ctrl *V1Controller) HandleItemCheckout() errchain.HandlerFunc {
	fn := func(r *http.Request, itemID uuid.UUID, input repo.CheckoutInput) (repo.ItemOut, error) {
		auth := services.NewContext(r.Context())
		out, err := ctrl.repo.Checkouts.Checkout(r.Context(), auth.GID, itemID, input, auth.User.Name)
		if errors.Is(err, repo.ErrItemUnavailable) {
			return repo.ItemOut{}, validate.NewRequestError(err, http.StatusConflict)
		}
		return out, err
	}
	return adapters.ActionID("id", fn, http.StatusOK)
}

func (ctrl *V1Controller) HandleItemReturn() errchain.HandlerFunc {
	fn := func(r *http.Request, itemID uuid.UUID, input repo.ReturnInput) (repo.ItemOut, error) {
		auth := services.NewContext(r.Context())
		out, err := ctrl.repo.Checkouts.Return(r.Context(), auth.GID, itemID, input, auth.User.Name)
		if errors.Is(err, repo.ErrNoActiveCheckout) {
			return repo.ItemOut{}, validate.NewRequestError(err, http.StatusConflict)
		}
		return out, err
	}
	return adapters.ActionID("id", fn, http.StatusOK)
}
