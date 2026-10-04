package v1

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hay-kot/httpkit/errchain"
	"github.com/hay-kot/httpkit/server"
	"github.com/sdmad737/catalog/backend/internal/core/services"
	"github.com/sdmad737/catalog/backend/internal/data/repo"
	"github.com/sdmad737/catalog/backend/internal/sys/validate"
)

type ResolvedItemTag struct {
	Tag        repo.ItemTagOut `json:"tag"`
	Item       repo.ItemOut    `json:"item"`
	CanOperate bool            `json:"canOperate"`
}

func tagToken(r *http.Request) (string, error) {
	token := strings.TrimSpace(chi.URLParam(r, "token"))
	if len(token) < 40 || len(token) > 64 {
		return "", validate.NewRequestError(errors.New("invalid tag"), http.StatusNotFound)
	}
	return token, nil
}

func itemTagIDs(r *http.Request) (uuid.UUID, uuid.UUID, error) {
	itemID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		return uuid.Nil, uuid.Nil, validate.NewRouteKeyError("id")
	}
	tagID := uuid.Nil
	if raw := chi.URLParam(r, "tag_id"); raw != "" {
		tagID, err = uuid.Parse(raw)
		if err != nil {
			return uuid.Nil, uuid.Nil, validate.NewRouteKeyError("tag_id")
		}
	}
	return itemID, tagID, nil
}

func (ctrl *V1Controller) HandlePublicItemTag() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		token, err := tagToken(r)
		if err != nil {
			return err
		}
		out, err := ctrl.repo.ItemTags.Public(r.Context(), token, true)
		if err != nil {
			return err
		}
		w.Header().Set("Cache-Control", "no-store")
		return server.JSON(w, http.StatusOK, out)
	}
}

func (ctrl *V1Controller) HandlePublicItemTagQRCode() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		token, err := tagToken(r)
		if err != nil {
			return err
		}
		tag, err := ctrl.repo.ItemTags.Public(r.Context(), token, false)
		if err != nil {
			return err
		}
		if tag.Status != "active" {
			return validate.NewRequestError(errors.New("inactive tag"), http.StatusGone)
		}

		origin := strings.TrimSuffix(strings.TrimSpace(r.URL.Query().Get("origin")), "/")
		if origin == "" {
			scheme := "https"
			if strings.HasPrefix(r.Host, "localhost") || strings.HasPrefix(r.Host, "127.0.0.1") {
				scheme = "http"
			}
			origin = scheme + "://" + r.Host
		}
		parsed, err := url.Parse(origin)
		if err != nil || parsed.Host == "" || (parsed.Scheme != "https" && parsed.Scheme != "http") || parsed.Path != "" {
			return validate.NewRequestError(errors.New("invalid origin"), http.StatusBadRequest)
		}
		if parsed.Scheme == "http" && parsed.Hostname() != "localhost" && parsed.Hostname() != "127.0.0.1" {
			return validate.NewRequestError(errors.New("tag links require HTTPS"), http.StatusBadRequest)
		}
		if !strings.EqualFold(parsed.Host, r.Host) {
			return validate.NewRequestError(errors.New("tag origin must match request host"), http.StatusBadRequest)
		}
		w.Header().Set("Cache-Control", "private, max-age=300")
		return writeQRCode(w, origin+"/t/"+token)
	}
}

func (ctrl *V1Controller) HandleResolveItemTag() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		token, err := tagToken(r)
		if err != nil {
			return err
		}
		auth := services.NewContext(r.Context())
		tag, item, err := ctrl.repo.ItemTags.Resolve(r.Context(), auth.GID, token)
		if err != nil {
			return err
		}
		actor := services.UseUserCtx(r.Context())
		canOperate := actor != nil && actor.Status != "disabled" &&
			(actor.Role == "owner" || actor.Role == "admin" || actor.Role == "staff")
		return server.JSON(w, http.StatusOK, ResolvedItemTag{Tag: tag, Item: item, CanOperate: canOperate})
	}
}

func (ctrl *V1Controller) HandleItemTagsList() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		itemID, _, err := itemTagIDs(r)
		if err != nil {
			return err
		}
		auth := services.NewContext(r.Context())
		items, err := ctrl.repo.ItemTags.List(r.Context(), auth.GID, itemID)
		if err != nil {
			return err
		}
		return server.JSON(w, http.StatusOK, WrapResults(items))
	}
}

func (ctrl *V1Controller) HandleItemTagAssign(replace bool) errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		itemID, _, err := itemTagIDs(r)
		if err != nil {
			return err
		}
		auth := services.NewContext(r.Context())
		var tag repo.ItemTagOut
		if replace {
			tag, err = ctrl.repo.ItemTags.Replace(r.Context(), auth.GID, itemID)
		} else {
			tag, err = ctrl.repo.ItemTags.Assign(r.Context(), auth.GID, itemID)
		}
		if errors.Is(err, repo.ErrActiveTagExists) {
			return validate.NewRequestError(err, http.StatusConflict)
		}
		if err != nil {
			return err
		}
		return server.JSON(w, http.StatusCreated, tag)
	}
}

func (ctrl *V1Controller) HandleItemTagRevoke() errchain.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		itemID, tagID, err := itemTagIDs(r)
		if err != nil {
			return err
		}
		auth := services.NewContext(r.Context())
		if err = ctrl.repo.ItemTags.Revoke(r.Context(), auth.GID, itemID, tagID); err != nil {
			return err
		}
		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}
