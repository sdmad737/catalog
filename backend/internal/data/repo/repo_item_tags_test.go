package repo

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemTagLifecycle(t *testing.T) {
	ctx := context.Background()
	inventoryItem := useItems(t, 1)[0]

	assigned, err := tRepos.ItemTags.Assign(ctx, tGroup.ID, inventoryItem.ID)
	require.NoError(t, err)
	assert.Len(t, assigned.Token, 43)
	assert.Nil(t, assigned.RevokedAt)

	preview, err := tRepos.ItemTags.Public(ctx, assigned.Token, false)
	require.NoError(t, err)
	assert.Equal(t, "active", preview.Status)
	listed, err := tRepos.ItemTags.List(ctx, tGroup.ID, inventoryItem.ID)
	require.NoError(t, err)
	require.Len(t, listed, 1)
	assert.Nil(t, listed[0].VerifiedAt, "rendering a QR code must not verify the physical setup")

	public, err := tRepos.ItemTags.Public(ctx, assigned.Token, true)
	require.NoError(t, err)
	assert.Equal(t, "active", public.Status)
	assert.Equal(t, inventoryItem.Name, public.Name)
	assert.Equal(t, tGroup.Name, public.Organization)
	encoded, err := json.Marshal(public)
	require.NoError(t, err)
	var publicPayload map[string]any
	require.NoError(t, json.Unmarshal(encoded, &publicPayload))
	for _, privateField := range []string{
		"id", "purchasePrice", "purchaseFrom", "notes", "missingNote", "lastKnownLocation",
		"attachments", "borrower", "activeCheckout", "serialNumber", "location",
	} {
		assert.NotContains(t, publicPayload, privateField, "anonymous tag payload must not serialize %s", privateField)
	}

	listed, err = tRepos.ItemTags.List(ctx, tGroup.ID, inventoryItem.ID)
	require.NoError(t, err)
	require.Len(t, listed, 1)
	assert.NotNil(t, listed[0].VerifiedAt, "opening the public URL must verify the tag")
	assert.NotNil(t, listed[0].LastScannedAt)

	_, err = tRepos.ItemTags.Assign(ctx, tGroup.ID, inventoryItem.ID)
	assert.ErrorIs(t, err, ErrActiveTagExists)

	replacement, err := tRepos.ItemTags.Replace(ctx, tGroup.ID, inventoryItem.ID)
	require.NoError(t, err)
	assert.NotEqual(t, assigned.Token, replacement.Token)

	old, err := tRepos.ItemTags.Public(ctx, assigned.Token, true)
	require.NoError(t, err)
	assert.Equal(t, "revoked", old.Status)
	assert.Empty(t, old.Name, "revoked tags must not disclose item details")

	tag, resolved, err := tRepos.ItemTags.Resolve(ctx, tGroup.ID, replacement.Token)
	require.NoError(t, err)
	assert.Equal(t, inventoryItem.ID, resolved.ID)
	assert.NotNil(t, tag.LastScannedAt)

	require.NoError(t, tRepos.ItemTags.Revoke(ctx, tGroup.ID, inventoryItem.ID, replacement.ID))
	_, _, err = tRepos.ItemTags.Resolve(ctx, tGroup.ID, replacement.Token)
	assert.Error(t, err)
}

func TestItemTagsAreGroupScoped(t *testing.T) {
	ctx := context.Background()
	inventoryItem := useItems(t, 1)[0]
	assigned, err := tRepos.ItemTags.Assign(ctx, tGroup.ID, inventoryItem.ID)
	require.NoError(t, err)

	otherGroup, err := tRepos.Groups.GroupCreate(ctx, "other-tag-group")
	require.NoError(t, err)
	t.Cleanup(func() { _ = tClient.Group.DeleteOneID(otherGroup.ID).Exec(ctx) })
	_, _, err = tRepos.ItemTags.Resolve(ctx, otherGroup.ID, assigned.Token)
	assert.Error(t, err)

	err = tRepos.ItemTags.Revoke(ctx, otherGroup.ID, inventoryItem.ID, assigned.ID)
	assert.Error(t, err)

	_, err = tRepos.ItemTags.Assign(ctx, otherGroup.ID, uuid.New())
	assert.Error(t, err)
}
