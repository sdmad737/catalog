package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hay-kot/httpkit/errchain"
	"github.com/hay-kot/httpkit/server"
	"github.com/rs/zerolog/log"
	"github.com/sdmad737/catalog/backend/internal/core/services"
	"github.com/sdmad737/catalog/backend/internal/data/ent"
	"github.com/sdmad737/catalog/backend/internal/data/ent/user"
	"github.com/sdmad737/catalog/backend/internal/data/repo"
	"github.com/sdmad737/catalog/backend/internal/sys/validate"
)

type StaffRoleUpdate struct {
	Role string `json:"role" validate:"required,oneof=admin staff viewer"`
}

type StaffStatusUpdate struct {
	Disabled bool `json:"disabled"`
}

type OwnershipTransfer struct {
	UserID uuid.UUID `json:"userId" validate:"required"`
}

func (ctrl *V1Controller) HandleUsersList() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := requireOrganizationRole(r, "owner", "admin"); err != nil {
			return err
		}
		auth := services.NewContext(r.Context())
		users, err := ctrl.repo.Users.GetByGroup(r.Context(), auth.GID)
		if err != nil {
			return err
		}
		return server.JSON(w, http.StatusOK, WrapResults(users))
	}
}

func (ctrl *V1Controller) HandleUserRoleUpdate() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := requireOrganizationRole(r, "owner"); err != nil {
			return err
		}
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			return validate.NewRouteKeyError("id")
		}
		var body StaffRoleUpdate
		if err = server.Decode(r, &body); err != nil {
			return validate.NewRequestError(err, http.StatusBadRequest)
		}
		if body.Role != "admin" && body.Role != "staff" && body.Role != "viewer" {
			return validate.NewRequestError(errors.New("invalid staff role"), http.StatusBadRequest)
		}
		auth := services.NewContext(r.Context())
		target, err := ctrl.repo.Users.GetOneID(r.Context(), id)
		if err != nil || target.GroupID != auth.GID {
			return &ent.NotFoundError{}
		}
		if target.Role == "owner" {
			return validate.NewRequestError(errors.New("transfer ownership before changing the owner role"), http.StatusConflict)
		}
		updated, err := ctrl.repo.Users.SetRole(r.Context(), auth.GID, id, user.Role(body.Role))
		if err != nil {
			return err
		}
		return server.JSON(w, http.StatusOK, Wrap(updated))
	}
}

func (ctrl *V1Controller) HandleUserStatusUpdate() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := requireOrganizationRole(r, "owner", "admin"); err != nil {
			return err
		}
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			return validate.NewRouteKeyError("id")
		}
		var body StaffStatusUpdate
		if err = server.Decode(r, &body); err != nil {
			return validate.NewRequestError(err, http.StatusBadRequest)
		}
		auth := services.NewContext(r.Context())
		if id == auth.UID {
			return validate.NewRequestError(errors.New("you cannot disable your own account"), http.StatusBadRequest)
		}
		target, err := ctrl.repo.Users.GetOneID(r.Context(), id)
		if err != nil || target.GroupID != auth.GID {
			return &ent.NotFoundError{}
		}
		if target.Role == "owner" {
			return validate.NewRequestError(errors.New("transfer ownership before disabling the owner"), http.StatusConflict)
		}
		if auth.User.Role == "admin" && target.Role == "admin" {
			return validate.NewRequestError(errors.New("only the owner can disable an administrator"), http.StatusForbidden)
		}
		updated, err := ctrl.repo.Users.SetDisabled(r.Context(), auth.GID, id, body.Disabled)
		if err != nil {
			return err
		}
		return server.JSON(w, http.StatusOK, Wrap(updated))
	}
}

func (ctrl *V1Controller) HandleOwnershipTransfer() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := requireOrganizationRole(r, "owner"); err != nil {
			return err
		}
		var body OwnershipTransfer
		if err := server.Decode(r, &body); err != nil {
			return validate.NewRequestError(err, http.StatusBadRequest)
		}
		auth := services.NewContext(r.Context())
		if err := ctrl.repo.Users.TransferOwnership(r.Context(), auth.GID, auth.UID, body.UserID); err != nil {
			return err
		}
		return server.JSON(w, http.StatusNoContent, nil)
	}
}

// HandleUserRegistration godoc
//
//	@Summary Register New User
//	@Tags    User
//	@Produce json
//	@Param   payload body services.UserRegistration true "User Data"
//	@Success 204
//	@Router  /v1/users/register [Post]
func (ctrl *V1Controller) HandleUserRegistration() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		regData := services.UserRegistration{}

		if err := server.Decode(r, &regData); err != nil {
			log.Err(err).Msg("failed to decode user registration data")
			return validate.NewRequestError(err, http.StatusInternalServerError)
		}

		if !ctrl.allowRegistration && regData.GroupToken == "" {
			return validate.NewRequestError(fmt.Errorf("user registration disabled"), http.StatusForbidden)
		}

		_, err := ctrl.svc.User.RegisterUser(r.Context(), regData)
		if err != nil {
			log.Err(err).Msg("failed to register user")
			return validate.NewRequestError(err, http.StatusInternalServerError)
		}

		return server.JSON(w, http.StatusNoContent, nil)
	}
}

// HandleUserSelf godoc
//
//	@Summary  Get User Self
//	@Tags     User
//	@Produce  json
//	@Success  200 {object} Wrapped{item=repo.UserOut}
//	@Router   /v1/users/self [GET]
//	@Security Bearer
func (ctrl *V1Controller) HandleUserSelf() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		token := services.UseTokenCtx(r.Context())
		usr, err := ctrl.svc.User.GetSelf(r.Context(), token)
		if usr.ID == uuid.Nil || err != nil {
			log.Err(err).Msg("failed to get user")
			return validate.NewRequestError(err, http.StatusInternalServerError)
		}

		return server.JSON(w, http.StatusOK, Wrap(usr))
	}
}

// HandleUserSelfUpdate godoc
//
//	@Summary  Update Account
//	@Tags     User
//	@Produce  json
//	@Param    payload body     repo.UserUpdate true "User Data"
//	@Success  200     {object} Wrapped{item=repo.UserUpdate}
//	@Router   /v1/users/self [PUT]
//	@Security Bearer
func (ctrl *V1Controller) HandleUserSelfUpdate() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		updateData := repo.UserUpdate{}
		if err := server.Decode(r, &updateData); err != nil {
			log.Err(err).Msg("failed to decode user update data")
			return validate.NewRequestError(err, http.StatusBadRequest)
		}

		actor := services.UseUserCtx(r.Context())
		newData, err := ctrl.svc.User.UpdateSelf(r.Context(), actor.ID, updateData)
		if err != nil {
			return validate.NewRequestError(err, http.StatusInternalServerError)
		}

		return server.JSON(w, http.StatusOK, Wrap(newData))
	}
}

// HandleUserSelfDelete godoc
//
//	@Summary  Delete Account
//	@Tags     User
//	@Produce  json
//	@Success  204
//	@Router   /v1/users/self [DELETE]
//	@Security Bearer
func (ctrl *V1Controller) HandleUserSelfDelete() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if ctrl.isDemo {
			return validate.NewRequestError(nil, http.StatusForbidden)
		}

		actor := services.UseUserCtx(r.Context())
		if err := ctrl.svc.User.DeleteSelf(r.Context(), actor.ID); err != nil {
			return validate.NewRequestError(err, http.StatusInternalServerError)
		}

		return server.JSON(w, http.StatusNoContent, nil)
	}
}

type (
	ChangePassword struct {
		Current string `json:"current,omitempty"`
		New     string `json:"new,omitempty"`
	}
)

// HandleUserSelfChangePassword godoc
//
//	@Summary  Change Password
//	@Tags     User
//	@Success  204
//	@Param    payload body ChangePassword true "Password Payload"
//	@Router   /v1/users/change-password [PUT]
//	@Security Bearer
func (ctrl *V1Controller) HandleUserSelfChangePassword() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if ctrl.isDemo {
			return validate.NewRequestError(nil, http.StatusForbidden)
		}

		var cp ChangePassword
		err := server.Decode(r, &cp)
		if err != nil {
			log.Err(err).Msg("user failed to change password")
		}

		ctx := services.NewContext(r.Context())

		ok := ctrl.svc.User.ChangePassword(ctx, cp.Current, cp.New)
		if !ok {
			return validate.NewRequestError(err, http.StatusInternalServerError)
		}

		return server.JSON(w, http.StatusNoContent, nil)
	}
}
