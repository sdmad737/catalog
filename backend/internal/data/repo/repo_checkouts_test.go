package repo

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckoutReturnAndMissingLifecycle(t *testing.T) {
	ctx := context.Background()
	asset := useItems(t, 1)[0]
	borrower, err := tRepos.People.Create(ctx, tGroup.ID, PersonCreate{
		Name: "Jordan Lee", Identifier: "STU-014", Email: "jordan@example.test",
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_, _ = tClient.Checkout.Delete().Exec(ctx)
		_ = tClient.Person.DeleteOneID(borrower.ID).Exec(ctx)
	})

	due := time.Now().AddDate(0, 0, 14)
	checkedOut, err := tRepos.Checkouts.Checkout(ctx, tGroup.ID, asset.ID, CheckoutInput{
		PersonID: borrower.ID, DueAt: &due, Condition: "good", Note: "Includes case and mouthpiece",
	}, "Alex Staff")
	require.NoError(t, err)
	assert.Equal(t, "checked_out", checkedOut.Status)
	require.NotNil(t, checkedOut.ActiveCheckout)
	assert.Equal(t, borrower.ID, checkedOut.ActiveCheckout.Person.ID)
	assert.Equal(t, "Jordan Lee", checkedOut.ActiveCheckout.Person.Name)
	assert.Equal(t, "good", checkedOut.Condition)

	_, err = tRepos.Checkouts.Checkout(ctx, tGroup.ID, asset.ID, CheckoutInput{
		PersonID: borrower.ID, Condition: "good",
	}, "Alex Staff")
	assert.ErrorIs(t, err, ErrItemUnavailable)

	missing, err := tRepos.ItemEvents.MarkMissing(ctx, tGroup.ID, asset.ID, MissingItemInput{Note: "Not in its assigned locker"}, "Alex Staff")
	require.NoError(t, err)
	assert.Equal(t, "missing", missing.Status)

	found, err := tRepos.ItemEvents.MarkFound(ctx, tGroup.ID, asset.ID, FoundItemInput{Condition: "fair", Note: "Found backstage"}, "Alex Staff")
	require.NoError(t, err)
	assert.Equal(t, "checked_out", found.Status, "finding a checked-out item must not erase its borrower assignment")
	require.NotNil(t, found.ActiveCheckout)

	returned, err := tRepos.Checkouts.Return(ctx, tGroup.ID, asset.ID, ReturnInput{
		LocationID: asset.Location.ID, Condition: "fair", Note: "Returned without damage",
	}, "Alex Staff")
	require.NoError(t, err)
	assert.Equal(t, "available", returned.Status)
	assert.Nil(t, returned.ActiveCheckout)
	assert.Equal(t, "fair", returned.Condition)

	history, err := tRepos.ItemEvents.List(ctx, tGroup.ID, asset.ID)
	require.NoError(t, err)
	kinds := make([]string, len(history))
	for i, event := range history {
		kinds[i] = event.Kind
	}
	assert.ElementsMatch(t, []string{"checkout", "missing", "found", "return"}, kinds)
}

func TestCheckoutIsOrganizationScoped(t *testing.T) {
	ctx := context.Background()
	asset := useItems(t, 1)[0]
	otherGroup, err := tRepos.Groups.GroupCreate(ctx, "other-checkout-group")
	require.NoError(t, err)
	otherPerson, err := tRepos.People.Create(ctx, otherGroup.ID, PersonCreate{Name: "Other Borrower"})
	require.NoError(t, err)
	t.Cleanup(func() { _ = tClient.Group.DeleteOneID(otherGroup.ID).Exec(ctx) })

	_, err = tRepos.Checkouts.Checkout(ctx, tGroup.ID, asset.ID, CheckoutInput{
		PersonID: otherPerson.ID, Condition: "good",
	}, "Staff")
	assert.Error(t, err)

	active, err := tRepos.Checkouts.ActiveForItem(ctx, otherGroup.ID, asset.ID)
	require.NoError(t, err)
	assert.Nil(t, active)
}
