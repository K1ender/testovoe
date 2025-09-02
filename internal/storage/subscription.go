package storage

import (
	"context"
	"database/sql"
	"errors"
	"testovoe/internal/models"
	"testovoe/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
)

const (
	MaxFutureDate = 3000
	MaxLimit      = 1000
)

//go:generate mockgen -source=subscription.go -destination=mocks/subscription.go
type SubscriptionStorage interface {
	Create(ctx context.Context, sub *models.Subscription) (int, error)
	Get(ctx context.Context, id int) (*models.Subscription, error)
	Update(ctx context.Context, id int, sub *models.Subscription) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, userID, serviceName string, limit, offset int) ([]models.Subscription, error)
	TotalForPeriod(
		ctx context.Context,
		periodStart, periodEnd time.Time,
		userID uuid.UUID,
		serviceName string,
	) (int64, error)
}

type PostgresSubscriptionStorage struct {
	db *sql.DB
}

func NewPostgresSubscriptionStorage(db *sql.DB) SubscriptionStorage {
	return &PostgresSubscriptionStorage{
		db: db,
	}
}

func (s *PostgresSubscriptionStorage) Create(ctx context.Context, sub *models.Subscription) (int, error) {
	var id int
	// query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	query, args := sqlbuilder.PostgreSQL.NewInsertBuilder().InsertInto("subscriptions").
		Cols("service_name", "price", "user_id", "start_date", "end_date").
		Values(
			sub.ServiceName,
			sub.Price,
			sub.UserID,
			sub.StartDate,
			sub.EndDate,
		).Returning("id").Build()

	err := s.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *PostgresSubscriptionStorage) Get(ctx context.Context, id int) (*models.Subscription, error) {
	query, args := sqlbuilder.PostgreSQL.NewSelectBuilder().Select("id", "service_name", "price", "user_id", "start_date", "end_date").
		From("subscriptions").
		Where(sqlbuilder.NewCond().Equal("id", id)).
		Build()
	row := s.db.QueryRowContext(ctx, query, args...)

	var sub models.Subscription
	if err := row.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &sub, nil
}

func (s *PostgresSubscriptionStorage) Update(ctx context.Context, id int, sub *models.Subscription) error {
	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder().Update("subscriptions")
	ub.Set(
		ub.Assign("service_name", sub.ServiceName),
		ub.Assign("price", sub.Price),
		ub.Assign("user_id", sub.UserID),
		ub.Assign("start_date", sub.StartDate),
		ub.Assign("end_date", sub.EndDate),
	).Where(ub.Equal("id", id))
	q, args := ub.Build()

	res, err := s.db.ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresSubscriptionStorage) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresSubscriptionStorage) List(ctx context.Context, userID, serviceName string, limit, offset int) ([]models.Subscription, error) {
	if limit > MaxLimit {
		limit = MaxLimit
	}

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("id", "service_name", "price", "user_id", "start_date", "end_date").From("subscriptions")

	var conds []string
	if userID != "" {
		conds = append(conds, sb.Equal("user_id", userID))
	}
	if serviceName != "" {
		conds = append(conds, sb.Equal("service_name", serviceName))
	}
	if len(conds) > 0 {
		sb.Where(sb.And(conds...))
	}

	sb.OrderBy("id").Desc()
	if limit > 0 {
		sb.Limit(limit)
	}
	if offset > 0 {
		sb.Offset(offset)
	}

	q, args := sb.Build()

	var out []models.Subscription
	rows, err := s.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sub models.Subscription
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate); err != nil {
			return nil, err
		}
		out = append(out, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *PostgresSubscriptionStorage) TotalForPeriod(
	ctx context.Context,
	periodStart, periodEnd time.Time,
	userID uuid.UUID, // Изменить тип
	serviceName string,
) (int64, error) {
	peEnd := periodEnd.AddDate(0, 1, 0).Add(-time.Nanosecond)

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("start_date", "end_date", "price").
		From("subscriptions").
		Where(sb.LessEqualThan("start_date", peEnd)).
		Where(sb.Or(sb.IsNull("end_date"), sb.GreaterEqualThan("end_date", periodStart))).
		Desc()

	if userID != uuid.Nil {
		sb.Where(sb.Equal("user_id", userID.String()))
	}

	if serviceName != "" {
		sb.Where(sb.Equal("service_name", serviceName))
	}

	sqlStr, args := sb.Build()

	rows, err := s.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var total int64
	for rows.Next() {
		var sub models.Subscription
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate); err != nil {
			return 0, err
		}

		subEnd := time.Date(MaxFutureDate, 1, 1, 0, 0, 0, 0, time.UTC)
		if sub.EndDate.Valid {
			subEnd = sub.EndDate.Time
		}

		months := utils.MonthsOverlap(sub.StartDate, subEnd, periodStart, periodEnd)
		total += int64(months) * int64(sub.Price)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	return total, nil
}
