package repo

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sdmad737/catalog/backend/internal/data/ent"
	"github.com/sdmad737/catalog/backend/internal/data/ent/group"
	"github.com/sdmad737/catalog/backend/internal/data/ent/item"
	"github.com/sdmad737/catalog/backend/internal/data/ent/itemtag"
)

var ErrActiveTagExists = errors.New("item already has an active physical tag")

type ItemTagRepository struct{ db *ent.Client }

type ItemTagOut struct {
	ID            uuid.UUID  `json:"id"`
	Token         string     `json:"token,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	RevokedAt     *time.Time `json:"revokedAt,omitempty"`
	VerifiedAt    *time.Time `json:"verifiedAt,omitempty"`
	LastScannedAt *time.Time `json:"lastScannedAt,omitempty"`
}

type PublicTagOut struct {
	Status       string  `json:"status"`
	Name         string  `json:"name,omitempty"`
	AssetID      AssetID `json:"assetId,omitempty" swaggertype:"string"`
	Organization string  `json:"organization,omitempty"`
	ItemStatus   string  `json:"itemStatus,omitempty"`
}

func mapItemTag(v *ent.ItemTag) ItemTagOut {
	return ItemTagOut{
		ID: v.ID, Token: v.PublicToken, CreatedAt: v.CreatedAt, RevokedAt: v.RevokedAt,
		VerifiedAt: v.VerifiedAt, LastScannedAt: v.LastScannedAt,
	}
}

func newPublicToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (r *ItemTagRepository) List(ctx context.Context, gid, itemID uuid.UUID) ([]ItemTagOut, error) {
	rows, err := r.db.ItemTag.Query().Where(
		itemtag.HasItemWith(item.ID(itemID), item.HasGroupWith(group.ID(gid))),
	).Order(ent.Desc(itemtag.FieldCreatedAt)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]ItemTagOut, len(rows))
	for i, row := range rows {
		out[i] = mapItemTag(row)
	}
	return out, nil
}

func (r *ItemTagRepository) Assign(ctx context.Context, gid, itemID uuid.UUID) (ItemTagOut, error) {
	exists, err := r.db.Item.Query().Where(item.ID(itemID), item.HasGroupWith(group.ID(gid))).Exist(ctx)
	if err != nil {
		return ItemTagOut{}, err
	}
	if !exists {
		return ItemTagOut{}, &ent.NotFoundError{}
	}
	active, err := r.db.ItemTag.Query().Where(itemtag.HasItemWith(item.ID(itemID)), itemtag.RevokedAtIsNil()).Exist(ctx)
	if err != nil {
		return ItemTagOut{}, err
	}
	if active {
		return ItemTagOut{}, ErrActiveTagExists
	}
	token, err := newPublicToken()
	if err != nil {
		return ItemTagOut{}, err
	}
	row, err := r.db.ItemTag.Create().SetItemID(itemID).SetPublicToken(token).Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			return ItemTagOut{}, ErrActiveTagExists
		}
		return ItemTagOut{}, err
	}
	return mapItemTag(row), nil
}

func (r *ItemTagRepository) Replace(ctx context.Context, gid, itemID uuid.UUID) (ItemTagOut, error) {
	if _, err := r.db.Item.Query().Where(item.ID(itemID), item.HasGroupWith(group.ID(gid))).Only(ctx); err != nil {
		return ItemTagOut{}, err
	}
	tx, err := r.db.Tx(ctx)
	if err != nil {
		return ItemTagOut{}, err
	}
	defer func() { _ = tx.Rollback() }()
	now := time.Now()
	if _, err = tx.ItemTag.Update().Where(itemtag.HasItemWith(item.ID(itemID)), itemtag.RevokedAtIsNil()).SetRevokedAt(now).Save(ctx); err != nil {
		return ItemTagOut{}, err
	}
	token, err := newPublicToken()
	if err != nil {
		return ItemTagOut{}, err
	}
	row, err := tx.ItemTag.Create().SetItemID(itemID).SetPublicToken(token).Save(ctx)
	if err != nil {
		return ItemTagOut{}, err
	}
	if err = tx.Commit(); err != nil {
		return ItemTagOut{}, err
	}
	return mapItemTag(row), nil
}

func (r *ItemTagRepository) Revoke(ctx context.Context, gid, itemID, tagID uuid.UUID) error {
	updated, err := r.db.ItemTag.Update().Where(
		itemtag.ID(tagID), itemtag.RevokedAtIsNil(),
		itemtag.HasItemWith(item.ID(itemID), item.HasGroupWith(group.ID(gid))),
	).SetRevokedAt(time.Now()).Save(ctx)
	if err != nil {
		return err
	}
	if updated == 0 {
		return &ent.NotFoundError{}
	}
	return nil
}

func (r *ItemTagRepository) Public(ctx context.Context, token string, recordOpen bool) (PublicTagOut, error) {
	row, err := r.db.ItemTag.Query().Where(itemtag.PublicToken(token)).WithItem(func(q *ent.ItemQuery) {
		q.WithGroup()
	}).Only(ctx)
	if err != nil {
		return PublicTagOut{}, err
	}
	if row.RevokedAt != nil {
		return PublicTagOut{Status: "revoked"}, nil
	}

	if recordOpen {
		// Creating or rendering a QR code is not verification. A tag becomes
		// verified only after its actual public page is opened.
		now := time.Now()
		update := row.Update().SetLastScannedAt(now)
		if row.VerifiedAt == nil {
			update.SetVerifiedAt(now)
		}
		if _, err = update.Save(ctx); err != nil {
			return PublicTagOut{}, err
		}
	}

	organization := ""
	if row.Edges.Item.Edges.Group != nil {
		organization = row.Edges.Item.Edges.Group.Name
	}
	return PublicTagOut{
		Status: "active", Name: row.Edges.Item.Name, AssetID: AssetID(row.Edges.Item.AssetID),
		Organization: organization, ItemStatus: string(row.Edges.Item.Status),
	}, nil
}

func (r *ItemTagRepository) Resolve(ctx context.Context, gid uuid.UUID, token string) (ItemTagOut, ItemOut, error) {
	row, err := r.db.ItemTag.Query().Where(
		itemtag.PublicToken(token), itemtag.RevokedAtIsNil(),
		itemtag.HasItemWith(item.HasGroupWith(group.ID(gid))),
	).WithItem().Only(ctx)
	if err != nil {
		return ItemTagOut{}, ItemOut{}, err
	}
	itemID := row.Edges.Item.ID
	now := time.Now()
	update := row.Update().SetLastScannedAt(now)
	if row.VerifiedAt == nil {
		update.SetVerifiedAt(now)
	}
	row, err = update.Save(ctx)
	if err != nil {
		return ItemTagOut{}, ItemOut{}, err
	}
	full, err := (&ItemsRepository{db: r.db}).GetOneByGroup(ctx, gid, itemID)
	return mapItemTag(row), full, err
}
