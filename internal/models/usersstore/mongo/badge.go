// Code generated by github.com/firstcontributions/matro. DO NOT EDIT.

package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/firstcontributions/backend/internal/models/usersstore"
	"github.com/firstcontributions/backend/internal/models/utils"
	"github.com/firstcontributions/backend/pkg/cursor"
	"github.com/gokultp/go-mongoqb"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func badgeFiltersToQuery(filters *usersstore.BadgeFilters) *mongoqb.QueryBuilder {
	qb := mongoqb.NewQueryBuilder()
	if len(filters.Ids) > 0 {
		qb.In("_id", filters.Ids)
	}
	if filters.User != nil {
		qb.Eq("user_id", filters.User.Id)
	}
	return qb
}
func (s *UsersStore) CreateBadge(ctx context.Context, badge *usersstore.Badge) (*usersstore.Badge, error) {
	now := time.Now()
	badge.TimeCreated = now
	badge.TimeUpdated = now
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	badge.Id = uuid.String()
	if _, err := s.getCollection(CollectionBadges).InsertOne(ctx, badge); err != nil {
		return nil, err
	}
	return badge, nil
}

func (s *UsersStore) GetBadgeByID(ctx context.Context, id string) (*usersstore.Badge, error) {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	var badge usersstore.Badge
	if err := s.getCollection(CollectionBadges).FindOne(ctx, qb.Build()).Decode(&badge); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &badge, nil
}

func (s *UsersStore) GetOneBadge(ctx context.Context, filters *usersstore.BadgeFilters) (*usersstore.Badge, error) {
	qb := badgeFiltersToQuery(filters)
	var badge usersstore.Badge
	if err := s.getCollection(CollectionBadges).FindOne(ctx, qb.Build()).Decode(&badge); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &badge, nil
}

func (s *UsersStore) CountBadges(ctx context.Context, filters *usersstore.BadgeFilters) (
	int64,
	error,
) {
	qb := badgeFiltersToQuery(filters)

	count, err := s.getCollection(CollectionBadges).CountDocuments(ctx, qb.Build())
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *UsersStore) GetBadges(
	ctx context.Context,
	filters *usersstore.BadgeFilters,
	after *string,
	before *string,
	first *int64,
	last *int64,
	sortBy *string,
	sortOrder *string,
) (
	[]*usersstore.Badge,
	bool,
	bool,
	string,
	string,
	error,
) {
	qb := badgeFiltersToQuery(filters)

	limit, order, cursorStr := utils.GetLimitAndSortOrderAndCursor(first, last, after, before)
	var c *cursor.Cursor
	if cursorStr != nil {
		c = cursor.FromString(*cursorStr)
		if c != nil {
			if order == 1 {
				qb.Or(
					qb.And(
						qb.Eq(c.SortBy, c.OffsetValue),
						qb.Gt("_id", c.ID),
					),
					qb.Gt(c.SortBy, c.OffsetValue),
				)
			} else {
				qb.Or(
					qb.And(
						qb.Eq(c.SortBy, c.OffsetValue),
						qb.Lt("_id", c.ID),
					),
					qb.Lt(c.SortBy, c.OffsetValue),
				)
			}
		}
	}
	// incrementing limit by 2 to check if next, prev elements are present
	limit += 2
	options := &options.FindOptions{
		Limit: &limit,
		Sort:  utils.GetSortOrder(sortBy, sortOrder, order),
	}

	var firstCursor, lastCursor string
	var hasNextPage, hasPreviousPage bool

	var badges []*usersstore.Badge
	mongoCursor, err := s.getCollection(CollectionBadges).Find(ctx, qb.Build(), options)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err
	}
	err = mongoCursor.All(ctx, &badges)
	if err != nil {
		return nil, hasNextPage, hasPreviousPage, firstCursor, lastCursor, err
	}
	count := len(badges)
	if count == 0 {
		return badges, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
	}

	// check if the cursor element present, if yes that can be a prev elem
	if c != nil && badges[0].Id == c.ID {
		hasPreviousPage = true
		badges = badges[1:]
		count--
	}

	// check if actual limit +1 elements are there, if yes trim it to limit
	if count >= int(limit)-1 {
		hasNextPage = true
		badges = badges[:limit-2]
		count = len(badges)
	}

	if count > 0 {
		firstCursor = cursor.NewCursor(badges[0].Id, "time_created", badges[0].TimeCreated).String()
		lastCursor = cursor.NewCursor(badges[count-1].Id, "time_created", badges[count-1].TimeCreated).String()
	}
	if order < 0 {
		hasNextPage, hasPreviousPage = hasPreviousPage, hasNextPage
		firstCursor, lastCursor = lastCursor, firstCursor
		badges = utils.ReverseList(badges)
	}
	return badges, hasNextPage, hasPreviousPage, firstCursor, lastCursor, nil
}

func (s *UsersStore) UpdateBadge(ctx context.Context, id string, badgeUpdate *usersstore.BadgeUpdate) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)

	now := time.Now()
	badgeUpdate.TimeUpdated = &now

	u := mongoqb.NewUpdateMap().
		SetFields(badgeUpdate)

	um, err := u.BuildUpdate()
	if err != nil {
		return err
	}
	if _, err := s.getCollection(CollectionBadges).UpdateOne(ctx, qb.Build(), um); err != nil {
		return err
	}
	return nil
}

func (s *UsersStore) DeleteBadgeByID(ctx context.Context, id string) error {
	qb := mongoqb.NewQueryBuilder().
		Eq("_id", id)
	if _, err := s.getCollection(CollectionBadges).DeleteOne(ctx, qb.Build()); err != nil {
		return err
	}
	return nil
}
