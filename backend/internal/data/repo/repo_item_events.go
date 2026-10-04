package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sdmad737/catalog/backend/internal/data/ent"
	"github.com/sdmad737/catalog/backend/internal/data/ent/checkout"
	"github.com/sdmad737/catalog/backend/internal/data/ent/group"
	"github.com/sdmad737/catalog/backend/internal/data/ent/item"
	"github.com/sdmad737/catalog/backend/internal/data/ent/itemevent"
	"github.com/sdmad737/catalog/backend/internal/data/ent/location"
)

type ItemEventOut struct {
	ID        uuid.UUID `json:"id"`
	Kind      string    `json:"kind"`
	Summary   string    `json:"summary"`
	Note      string    `json:"note,omitempty"`
	ActorName string    `json:"actorName,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type MissingItemInput struct {
	LastKnownLocation string `json:"lastKnownLocation" validate:"max=255"`
	Note              string `json:"note" validate:"max=1000"`
}

type FoundItemInput struct {
	LocationID uuid.UUID `json:"locationId,omitempty"`
	Condition  string    `json:"condition" validate:"required,oneof=unknown excellent good fair poor damaged"`
	Note       string    `json:"note" validate:"max=1000"`
}

type ItemEventRepository struct{ db *ent.Client }

func mapItemEvent(v *ent.ItemEvent) ItemEventOut {
	return ItemEventOut{ID: v.ID, Kind: v.Kind, Summary: v.Summary, Note: v.Note, ActorName: v.ActorName, CreatedAt: v.CreatedAt}
}

func (r *ItemEventRepository) List(ctx context.Context, gid, itemID uuid.UUID) ([]ItemEventOut, error) {
	rows, err := r.db.ItemEvent.Query().Where(
		itemevent.HasItemWith(item.ID(itemID), item.HasGroupWith(group.ID(gid))),
	).Order(ent.Desc(itemevent.FieldCreatedAt)).Limit(100).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]ItemEventOut, len(rows))
	for i, row := range rows {
		out[i] = mapItemEvent(row)
	}
	return out, nil
}

func (r *ItemEventRepository) add(ctx context.Context, itemID uuid.UUID, kind, summary, note, actor string) error {
	_, err := r.db.ItemEvent.Create().SetItemID(itemID).SetKind(kind).SetSummary(summary).SetNote(note).SetActorName(actor).Save(ctx)
	return err
}

func (r *ItemEventRepository) MarkMissing(ctx context.Context, gid, itemID uuid.UUID, input MissingItemInput, actor string) (ItemOut, error) {
	tx, err := r.db.Tx(ctx)
	if err != nil {
		return ItemOut{}, err
	}
	defer func() { _ = tx.Rollback() }()
	row, err := tx.Item.Query().Where(item.ID(itemID), item.HasGroupWith(group.ID(gid))).WithLocation().Only(ctx)
	if err != nil {
		return ItemOut{}, err
	}
	lastLocation := input.LastKnownLocation
	if lastLocation == "" && row.Edges.Location != nil {
		lastLocation = row.Edges.Location.Name
	}
	now := time.Now()
	if _, err = row.Update().SetStatus(item.StatusMissing).SetMissingSince(now).SetMissingNote(input.Note).SetLastKnownLocation(lastLocation).Save(ctx); err != nil {
		return ItemOut{}, err
	}
	summary := "Marked missing"
	if lastLocation != "" {
		summary = fmt.Sprintf("Marked missing — last known at %s", lastLocation)
	}
	if _, err = tx.ItemEvent.Create().SetItemID(itemID).SetKind("missing").SetSummary(summary).
		SetNote(input.Note).SetActorName(actor).Save(ctx); err != nil {
		return ItemOut{}, err
	}
	if err = tx.Commit(); err != nil {
		return ItemOut{}, err
	}
	return (&ItemsRepository{db: r.db}).GetOneByGroup(ctx, gid, itemID)
}

func (r *ItemEventRepository) MarkFound(ctx context.Context, gid, itemID uuid.UUID, input FoundItemInput, actor string) (ItemOut, error) {
	tx, err := r.db.Tx(ctx)
	if err != nil {
		return ItemOut{}, err
	}
	defer func() { _ = tx.Rollback() }()
	row, err := tx.Item.Query().Where(item.ID(itemID), item.HasGroupWith(group.ID(gid)), item.StatusEQ(item.StatusMissing)).Only(ctx)
	if err != nil {
		return ItemOut{}, err
	}
	restoredStatus := item.StatusAvailable
	hasActiveCheckout, err := tx.Checkout.Query().Where(checkout.ReturnedAtIsNil(), checkout.HasItemWith(item.ID(itemID))).Exist(ctx)
	if err != nil {
		return ItemOut{}, err
	}
	if hasActiveCheckout {
		restoredStatus = item.StatusCheckedOut
	}
	update := row.Update().SetStatus(restoredStatus).SetCondition(item.Condition(input.Condition)).ClearMissingSince().SetMissingNote("").SetLastKnownLocation("")
	if input.LocationID != uuid.Nil {
		valid, err := tx.Location.Query().Where(location.ID(input.LocationID), location.HasGroupWith(group.ID(gid))).Exist(ctx)
		if err != nil {
			return ItemOut{}, err
		}
		if !valid {
			return ItemOut{}, fmt.Errorf("found location does not belong to this organization")
		}
		update.SetLocationID(input.LocationID)
	}
	if _, err = update.Save(ctx); err != nil {
		return ItemOut{}, err
	}
	if _, err = tx.ItemEvent.Create().SetItemID(itemID).SetKind("found").SetSummary("Marked found").
		SetNote(input.Note).SetActorName(actor).Save(ctx); err != nil {
		return ItemOut{}, err
	}
	if err = tx.Commit(); err != nil {
		return ItemOut{}, err
	}
	return (&ItemsRepository{db: r.db}).GetOneByGroup(ctx, gid, itemID)
}
