package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"user-service/internal/testutil"
	"user-service/internal/util/liberr"
	"user-service/module/auth/entity"
	"user-service/module/auth/internal/repository"
	"user-service/module/auth/testutil/fixtures"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	userAllAttributes = []string{
		"id",
		"name",
		"email",
		"phone",
		"password",
		"created_at",
		"updated_at",
	}

	userAllColumnsStr = strings.Join(userAllAttributes, ", ")
)

func TestUserRepository_GetByUsername(t *testing.T) {
	columns := userAllColumnsStr
	rows := userAllAttributes
	expectedQuery := fmt.Sprintf("SELECT %s FROM users WHERE (phone = ? OR email = ?)", columns)
	dummyUser := fixtures.NewUser(fixtures.User)

	type input struct {
		ctx      context.Context
		username string
	}

	testCases := []struct {
		name           string
		in             input
		mockDependency func(*testutil.RepositoryDependency, input)
		assertFn       func(*entity.User, error)
	}{
		{
			name: "Success on GetByEmail",
			in: input{
				ctx:      context.TODO(),
				username: "jhon.doe@test.com",
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(in.username, in.username).
					WillReturnRows(
						sqlmock.
							NewRows(rows).
							AddRow(fixtures.GetUserRow(dummyUser)...),
					).RowsWillBeClosed()
			},
			assertFn: func(result *entity.User, err error) {
				assert.Nil(t, err)
				assert.Equal(t, dummyUser, result)
			},
		},
		{
			name: "Error on Execute Query with Not Found Row",
			in: input{
				ctx:      context.TODO(),
				username: "jhon.doe@test.com",
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(in.username, in.username).
					WillReturnError(sql.ErrNoRows)
			},
			assertFn: func(result *entity.User, err error) {
				assert.NotNil(t, err)

				berr, ok := err.(*liberr.BaseError)
				assert.True(t, ok)

				assert.Equal(t, entity.ErrorUserNotFound, berr.GetDetails()[0])
				assert.Nil(t, result)
			},
		},
		{
			name: "Error on Execute Query",
			in: input{
				ctx:      context.TODO(),
				username: "jhon.doe@test.com",
			},
			mockDependency: func(dependency *testutil.RepositoryDependency, in input) {
				dependency.MockedSQL.
					ExpectQuery(regexp.QuoteMeta(expectedQuery)).
					WithArgs(in.username, in.username).
					WillReturnError(sqlmock.ErrCancelled).
					RowsWillBeClosed()
			},
			assertFn: func(result *entity.User, err error) {
				assert.NotNil(t, err)
				assert.Nil(t, result)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repositoryDependency := testutil.NewRepositoryDependency()
			repo := repository.NewUserRepository(repositoryDependency.MockedDB)

			defer ctrl.Finish()

			tc.mockDependency(&repositoryDependency, tc.in)
			tc.assertFn(repo.GetByUsername(tc.in.ctx, tc.in.username))
		})
	}
}
