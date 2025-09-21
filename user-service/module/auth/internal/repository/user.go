package repository

import (
	"context"
	"database/sql"
	"time"
	"user-service/internal/util/liberr"
	"user-service/module/auth/entity"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

var (
	userTable = "users"

	userColumns = []string{"id", "name", "email", "phone", "password", "created_at", "updated_at"}
)

type UserRepository struct {
	db *sqlx.DB
}

type userObject struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Email     *string   `db:"email"`
	Phone     *string   `db:"phone"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (o *userObject) toEntity() *entity.User {
	user := &entity.User{
		ID:        o.ID,
		Name:      o.Name,
		Password:  o.Password,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}

	if o.Email != nil {
		user.Email = *o.Email
	}
	if o.Phone != nil {
		user.Phone = *o.Phone
	}

	return user
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select(userColumns...)
	sb.From(userTable)
	sb.Where(
		sb.Or(
			sb.Equal("phone", username),
			sb.Equal("email", username),
		),
	)

	query, args := sb.Build()

	row := u.db.QueryRowxContext(ctx, query, args...)
	obj := &userObject{}

	if err := row.StructScan(obj); err != nil {
		if err == sql.ErrNoRows {
			return nil, liberr.NewBaseError(entity.ErrorUserNotFound)
		}
		return nil, liberr.NewTracer("Error when StructScan on user.GetByEmail").Wrap(err)
	}

	return obj.toEntity(), nil
}
