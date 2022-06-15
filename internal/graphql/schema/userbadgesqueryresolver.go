package schema

import (
	"context"

	"github.com/firstcontributions/backend/internal/storemanager"
)

type UserBadgesInput struct {
	First  *int32
	Last   *int32
	After  *string
	Before *string
}

func (n *User) Badges(ctx context.Context, in *UserBadgesInput) (*BadgesConnection, error) {
	var first, last *int64
	if in.First != nil {
		tmp := int64(*in.First)
		first = &tmp
	}
	if in.Last != nil {
		tmp := int64(*in.Last)
		last = &tmp
	}
	store := storemanager.FromContext(ctx)
	data, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err := store.UsersStore.GetBadges(
		ctx,
		nil,
		n.ref,
		in.After,
		in.Before,
		first,
		last,
	)
	if err != nil {
		return nil, err
	}
	return NewBadgesConnection(data, hasNextPage, hasPreviousPage, &firstCursor, &lastCursor), nil
}
