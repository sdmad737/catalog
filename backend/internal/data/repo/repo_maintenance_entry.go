package repo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sdmad737/catalog/backend/internal/data/ent"
	"github.com/sdmad737/catalog/backend/internal/data/ent/group"
	"github.com/sdmad737/catalog/backend/internal/data/ent/item"
	"github.com/sdmad737/catalog/backend/internal/data/ent/maintenanceentry"
	"github.com/sdmad737/catalog/backend/internal/data/types"
)

// MaintenanceEntryRepository is a repository for maintenance entries that are
// associated with an item in the database. An entry represents a maintenance event
// that has been performed on an item.
type MaintenanceEntryRepository struct {
	db *ent.Client
}

type MaintenanceEntryCreate struct {
	CompletedDate types.Date `json:"completedDate"`
	ScheduledDate types.Date `json:"scheduledDate"`
	Name          string     `json:"name"          validate:"required"`
	Description   string     `json:"description"`
	Cost          float64    `json:"cost,string"`
}

func (mc MaintenanceEntryCreate) Validate() error {
	if mc.CompletedDate.Time().IsZero() && mc.ScheduledDate.Time().IsZero() {
		return errors.New("either completedDate or scheduledDate must be set")
	}
	return nil
}

type MaintenanceEntryUpdate struct {
	CompletedDate types.Date `json:"completedDate"`
	ScheduledDate types.Date `json:"scheduledDate"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Cost          float64    `json:"cost,string"`
}

func (mu MaintenanceEntryUpdate) Validate() error {
	if mu.CompletedDate.Time().IsZero() && mu.ScheduledDate.Time().IsZero() {
		return errors.New("either completedDate or scheduledDate must be set")
	}
	return nil
}

type (
	MaintenanceEntry struct {
		ID            uuid.UUID  `json:"id"`
		CompletedDate types.Date `json:"completedDate"`
		ScheduledDate types.Date `json:"scheduledDate"`
		Name          string     `json:"name"`
		Description   string     `json:"description"`
		Cost          float64    `json:"cost,string"`
	}

	MaintenanceLog struct {
		ItemID      uuid.UUID          `json:"itemId"`
		CostAverage float64            `json:"costAverage"`
		CostTotal   float64            `json:"costTotal"`
		Entries     []MaintenanceEntry `json:"entries"`
	}
)

var (
	mapMaintenanceEntryErr  = mapTErrFunc(mapMaintenanceEntry)
	mapEachMaintenanceEntry = mapTEachFunc(mapMaintenanceEntry)
)

func mapMaintenanceEntry(entry *ent.MaintenanceEntry) MaintenanceEntry {
	return MaintenanceEntry{
		ID:            entry.ID,
		CompletedDate: types.Date(entry.Date),
		ScheduledDate: types.Date(entry.ScheduledDate),
		Name:          entry.Name,
		Description:   entry.Description,
		Cost:          entry.Cost,
	}
}

func (r *MaintenanceEntryRepository) GetScheduled(ctx context.Context, GID uuid.UUID, dt types.Date) ([]MaintenanceEntry, error) {
	entries, err := r.db.MaintenanceEntry.Query().
		Where(
			maintenanceentry.HasItemWith(
				item.HasGroupWith(group.ID(GID)),
			),
			maintenanceentry.ScheduledDate(dt.Time()),
			maintenanceentry.Or(
				maintenanceentry.DateIsNil(),
				maintenanceentry.DateEQ(time.Time{}),
			),
		).
		All(ctx)

	if err != nil {
		return nil, err
	}

	return mapEachMaintenanceEntry(entries), nil
}

func (r *MaintenanceEntryRepository) Create(ctx context.Context, groupID, itemID uuid.UUID, input MaintenanceEntryCreate) (MaintenanceEntry, error) {
	exists, err := r.db.Item.Query().Where(item.ID(itemID), item.HasGroupWith(group.ID(groupID))).Exist(ctx)
	if err != nil {
		return MaintenanceEntry{}, err
	}
	if !exists {
		return MaintenanceEntry{}, &ent.NotFoundError{}
	}
	item, err := r.db.MaintenanceEntry.Create().
		SetItemID(itemID).
		SetDate(input.CompletedDate.Time()).
		SetScheduledDate(input.ScheduledDate.Time()).
		SetName(input.Name).
		SetDescription(input.Description).
		SetCost(input.Cost).
		Save(ctx)

	return mapMaintenanceEntryErr(item, err)
}

func (r *MaintenanceEntryRepository) Update(ctx context.Context, groupID, itemID, ID uuid.UUID, input MaintenanceEntryUpdate) (MaintenanceEntry, error) {
	updated, err := r.db.MaintenanceEntry.Update().Where(
		maintenanceentry.ID(ID), maintenanceentry.ItemID(itemID),
		maintenanceentry.HasItemWith(item.HasGroupWith(group.ID(groupID))),
	).
		SetDate(input.CompletedDate.Time()).
		SetScheduledDate(input.ScheduledDate.Time()).
		SetName(input.Name).
		SetDescription(input.Description).
		SetCost(input.Cost).
		Save(ctx)
	if err != nil {
		return MaintenanceEntry{}, err
	}
	if updated != 1 {
		return MaintenanceEntry{}, &ent.NotFoundError{}
	}
	return mapMaintenanceEntryErr(r.db.MaintenanceEntry.Get(ctx, ID))
}

type MaintenanceLogQuery struct {
	Completed bool `json:"completed" schema:"completed"`
	Scheduled bool `json:"scheduled" schema:"scheduled"`
}

func (r *MaintenanceEntryRepository) GetLog(ctx context.Context, groupID, itemID uuid.UUID, query MaintenanceLogQuery) (MaintenanceLog, error) {
	log := MaintenanceLog{
		ItemID: itemID,
	}

	q := r.db.MaintenanceEntry.Query().Where(
		maintenanceentry.ItemID(itemID),
		maintenanceentry.HasItemWith(
			item.HasGroupWith(group.IDEQ(groupID)),
		),
	)

	if query.Completed {
		q = q.Where(maintenanceentry.And(
			maintenanceentry.DateNotNil(),
			maintenanceentry.DateNEQ(time.Time{}),
		))
	} else if query.Scheduled {
		q = q.Where(maintenanceentry.And(
			maintenanceentry.Or(
				maintenanceentry.DateIsNil(),
				maintenanceentry.DateEQ(time.Time{}),
			),
			maintenanceentry.ScheduledDateNotNil(),
			maintenanceentry.ScheduledDateNEQ(time.Time{}),
		))
	}

	entries, err := q.Order(ent.Desc(maintenanceentry.FieldDate)).
		All(ctx)
	if err != nil {
		return MaintenanceLog{}, err
	}

	log.Entries = mapEachMaintenanceEntry(entries)

	monthlyTotals := make(map[string]float64)
	for _, entry := range log.Entries {
		date := entry.CompletedDate.Time()
		if date.IsZero() {
			date = entry.ScheduledDate.Time()
		}

		log.CostTotal += entry.Cost
		month := date.Format("2006-01")
		monthlyTotals[month] += entry.Cost
	}

	if len(monthlyTotals) > 0 {
		log.CostAverage = log.CostTotal / float64(len(monthlyTotals))
	}
	return log, nil
}

func (r *MaintenanceEntryRepository) Delete(ctx context.Context, groupID, itemID, ID uuid.UUID) error {
	deleted, err := r.db.MaintenanceEntry.Delete().Where(
		maintenanceentry.ID(ID), maintenanceentry.ItemID(itemID),
		maintenanceentry.HasItemWith(item.HasGroupWith(group.ID(groupID))),
	).Exec(ctx)
	if err != nil {
		return err
	}
	if deleted != 1 {
		return &ent.NotFoundError{}
	}
	return nil
}
