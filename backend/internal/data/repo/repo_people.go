package repo

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sdmad737/catalog/backend/internal/data/ent"
	"github.com/sdmad737/catalog/backend/internal/data/ent/group"
	"github.com/sdmad737/catalog/backend/internal/data/ent/person"
)

type PersonOut struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email,omitempty"`
	Identifier string    `json:"identifier,omitempty"`
	Notes      string    `json:"notes,omitempty"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"createdAt"`
}

type PersonCreate struct {
	Name       string `json:"name" validate:"required,min=1,max=255"`
	Email      string `json:"email" validate:"omitempty,email,max=255"`
	Identifier string `json:"identifier" validate:"max=100"`
	Notes      string `json:"notes" validate:"max=1000"`
}

type PersonUpdate struct {
	Name       string `json:"name" validate:"required,min=1,max=255"`
	Email      string `json:"email" validate:"omitempty,email,max=255"`
	Identifier string `json:"identifier" validate:"max=100"`
	Notes      string `json:"notes" validate:"max=1000"`
	Active     bool   `json:"active"`
}

type PersonRepository struct{ db *ent.Client }

func mapPerson(v *ent.Person) PersonOut {
	return PersonOut{
		ID: v.ID, Name: v.Name, Email: v.Email, Identifier: v.Identifier,
		Notes: v.Notes, Active: v.Active, CreatedAt: v.CreatedAt,
	}
}

func (r *PersonRepository) List(ctx context.Context, gid uuid.UUID, search string) ([]PersonOut, error) {
	q := r.db.Person.Query().Where(person.HasGroupWith(group.ID(gid)))
	if search = strings.TrimSpace(search); search != "" {
		q.Where(person.Or(person.NameContainsFold(search), person.EmailContainsFold(search), person.IdentifierContainsFold(search)))
	}
	rows, err := q.Order(ent.Asc(person.FieldName)).All(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]PersonOut, len(rows))
	for i, row := range rows {
		out[i] = mapPerson(row)
	}
	return out, nil
}

func (r *PersonRepository) Create(ctx context.Context, gid uuid.UUID, input PersonCreate) (PersonOut, error) {
	row, err := r.db.Person.Create().SetGroupID(gid).SetName(strings.TrimSpace(input.Name)).
		SetEmail(strings.TrimSpace(input.Email)).SetIdentifier(strings.TrimSpace(input.Identifier)).SetNotes(input.Notes).Save(ctx)
	if err != nil {
		return PersonOut{}, err
	}
	return mapPerson(row), nil
}

func (r *PersonRepository) Update(ctx context.Context, gid, id uuid.UUID, input PersonUpdate) (PersonOut, error) {
	_, err := r.db.Person.Update().Where(person.ID(id), person.HasGroupWith(group.ID(gid))).
		SetName(strings.TrimSpace(input.Name)).SetEmail(strings.TrimSpace(input.Email)).
		SetIdentifier(strings.TrimSpace(input.Identifier)).SetNotes(input.Notes).SetActive(input.Active).Save(ctx)
	if err != nil {
		return PersonOut{}, err
	}
	row, err := r.db.Person.Query().Where(person.ID(id), person.HasGroupWith(group.ID(gid))).Only(ctx)
	if err != nil {
		return PersonOut{}, err
	}
	return mapPerson(row), nil
}
