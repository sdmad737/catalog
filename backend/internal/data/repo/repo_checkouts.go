package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sdmad737/catalog/backend/internal/data/ent"
	"github.com/sdmad737/catalog/backend/internal/data/ent/checkout"
	"github.com/sdmad737/catalog/backend/internal/data/ent/group"
	"github.com/sdmad737/catalog/backend/internal/data/ent/item"
	entlocation "github.com/sdmad737/catalog/backend/internal/data/ent/location"
	"github.com/sdmad737/catalog/backend/internal/data/ent/person"
)

var (
	ErrItemUnavailable  = errors.New("item is not available for checkout")
	ErrNoActiveCheckout = errors.New("item is not currently checked out")
)

type CheckoutOut struct {
	ID           uuid.UUID  `json:"id"`
	Person       PersonOut  `json:"person"`
	CheckedOutAt time.Time  `json:"checkedOutAt"`
	DueAt        *time.Time `json:"dueAt,omitempty"`
	ReturnedAt   *time.Time `json:"returnedAt,omitempty"`
	ConditionOut string     `json:"conditionOut"`
	ConditionIn  string     `json:"conditionIn,omitempty"`
	CheckoutNote string     `json:"checkoutNote,omitempty"`
	ReturnNote   string     `json:"returnNote,omitempty"`
	CheckedOutBy string     `json:"checkedOutBy,omitempty"`
	ReturnedBy   string     `json:"returnedBy,omitempty"`
}

type CheckoutInput struct {
	PersonID  uuid.UUID  `json:"personId" validate:"required"`
	DueAt     *time.Time `json:"dueAt,omitempty"`
	Condition string     `json:"condition" validate:"required,oneof=unknown excellent good fair poor damaged"`
	Note      string     `json:"note" validate:"max=1000"`
}

type ReturnInput struct {
	LocationID uuid.UUID `json:"locationId,omitempty"`
	Condition  string    `json:"condition" validate:"required,oneof=unknown excellent good fair poor damaged"`
	Note       string    `json:"note" validate:"max=1000"`
}

type CheckoutRepository struct{ db *ent.Client }

func mapCheckout(v *ent.Checkout) CheckoutOut {
	var borrower PersonOut
	if v.Edges.Person != nil {
		borrower = mapPerson(v.Edges.Person)
	}
	return CheckoutOut{
		ID: v.ID, Person: borrower, CheckedOutAt: v.CheckedOutAt, DueAt: v.DueAt,
		ReturnedAt: v.ReturnedAt, ConditionOut: v.ConditionOut.String(), ConditionIn: v.ConditionIn.String(),
		CheckoutNote: v.CheckoutNote, ReturnNote: v.ReturnNote, CheckedOutBy: v.CheckedOutBy, ReturnedBy: v.ReturnedBy,
	}
}

func (r *CheckoutRepository) ActiveForItem(ctx context.Context, gid, itemID uuid.UUID) (*CheckoutOut, error) {
	row, err := r.db.Checkout.Query().Where(
		checkout.ReturnedAtIsNil(), checkout.HasGroupWith(group.ID(gid)),
		checkout.HasItemWith(item.ID(itemID), item.HasGroupWith(group.ID(gid))),
	).WithPerson().Only(ctx)
	if ent.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	out := mapCheckout(row)
	return &out, nil
}

func (r *CheckoutRepository) Checkout(ctx context.Context, gid, itemID uuid.UUID, input CheckoutInput, actor string) (ItemOut, error) {
	tx, err := r.db.Tx(ctx)
	if err != nil {
		return ItemOut{}, err
	}
	defer func() { _ = tx.Rollback() }()

	asset, err := tx.Item.Query().Where(item.ID(itemID), item.HasGroupWith(group.ID(gid))).Only(ctx)
	if err != nil {
		return ItemOut{}, err
	}
	if asset.Status != item.StatusAvailable {
		return ItemOut{}, ErrItemUnavailable
	}
	borrower, err := tx.Person.Query().Where(person.ID(input.PersonID), person.Active(true), person.HasGroupWith(group.ID(gid))).Only(ctx)
	if err != nil {
		return ItemOut{}, err
	}
	hasActive, err := tx.Checkout.Query().Where(checkout.ReturnedAtIsNil(), checkout.HasItemWith(item.ID(itemID))).Exist(ctx)
	if err != nil {
		return ItemOut{}, err
	}
	if hasActive {
		return ItemOut{}, ErrItemUnavailable
	}

	_, err = tx.Checkout.Create().SetGroupID(gid).SetItemID(itemID).SetPersonID(borrower.ID).
		SetNillableDueAt(input.DueAt).SetConditionOut(checkout.ConditionOut(input.Condition)).
		SetCheckoutNote(input.Note).SetCheckedOutBy(actor).Save(ctx)
	if err != nil {
		return ItemOut{}, err
	}
	if _, err = tx.Item.UpdateOneID(itemID).SetStatus(item.StatusCheckedOut).SetCondition(item.Condition(input.Condition)).Save(ctx); err != nil {
		return ItemOut{}, err
	}
	if _, err = tx.ItemEvent.Create().SetItemID(itemID).SetKind("checkout").
		SetSummary(fmt.Sprintf("Checked out to %s", borrower.Name)).SetNote(input.Note).SetActorName(actor).Save(ctx); err != nil {
		return ItemOut{}, err
	}
	if err = tx.Commit(); err != nil {
		return ItemOut{}, err
	}
	return (&ItemsRepository{db: r.db}).GetOneByGroup(ctx, gid, itemID)
}

func (r *CheckoutRepository) Return(ctx context.Context, gid, itemID uuid.UUID, input ReturnInput, actor string) (ItemOut, error) {
	tx, err := r.db.Tx(ctx)
	if err != nil {
		return ItemOut{}, err
	}
	defer func() { _ = tx.Rollback() }()

	row, err := tx.Checkout.Query().Where(
		checkout.ReturnedAtIsNil(), checkout.HasGroupWith(group.ID(gid)),
		checkout.HasItemWith(item.ID(itemID), item.HasGroupWith(group.ID(gid))),
	).WithPerson().Only(ctx)
	if ent.IsNotFound(err) {
		return ItemOut{}, ErrNoActiveCheckout
	}
	if err != nil {
		return ItemOut{}, err
	}
	now := time.Now()
	if _, err = row.Update().SetReturnedAt(now).SetConditionIn(checkout.ConditionIn(input.Condition)).
		SetReturnNote(input.Note).SetReturnedBy(actor).Save(ctx); err != nil {
		return ItemOut{}, err
	}
	itemUpdate := tx.Item.UpdateOneID(itemID).SetStatus(item.StatusAvailable).SetCondition(item.Condition(input.Condition))
	if input.LocationID != uuid.Nil {
		valid, err := tx.Location.Query().Where(entlocation.ID(input.LocationID), entlocation.HasGroupWith(group.ID(gid))).Exist(ctx)
		if err != nil {
			return ItemOut{}, err
		}
		if !valid {
			return ItemOut{}, errors.New("return location does not belong to this organization")
		}
		itemUpdate.SetLocationID(input.LocationID)
	}
	if _, err = itemUpdate.Save(ctx); err != nil {
		return ItemOut{}, err
	}
	summary := fmt.Sprintf("Returned from %s", row.Edges.Person.Name)
	if _, err = tx.ItemEvent.Create().SetItemID(itemID).SetKind("return").SetSummary(summary).
		SetNote(input.Note).SetActorName(actor).Save(ctx); err != nil {
		return ItemOut{}, err
	}
	if err = tx.Commit(); err != nil {
		return ItemOut{}, err
	}
	return (&ItemsRepository{db: r.db}).GetOneByGroup(ctx, gid, itemID)
}
