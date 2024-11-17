package repository

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/mrizkisaputra/expenses-api/internal/user"
	"github.com/mrizkisaputra/expenses-api/internal/user/model"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"regexp"
	"testing"
	"time"
)

var (
	expectedCreateUserQuery   = regexp.QuoteMeta(`INSERT INTO "users" ("email","password","first_name","last_name","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)
	expectedUpdateUserQuery   = regexp.QuoteMeta(`UPDATE "users" SET "email"=$1,"password"=$2,"avatar"=$3,"first_name"=$4,"last_name"=$5,"city"=$6,"phone_number"=$7,"updated_at"=$8 WHERE id = $9`)
	expectedFindUserByIdQuery = regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 LIMIT $2`)
	expectedFindAlreadyEmail  = regexp.QuoteMeta(`SELECT count(*) FROM "users" WHERE "users"."email" = $1`)
)
var (
	repository user.UserPostgresRepository
	mock       sqlmock.Sqlmock
)

// Test Before and After
func TestMain(m *testing.M) {
	fmt.Println("=================================================================")
	fmt.Println("STARTING TEST")
	fmt.Println("=================================================================")

	// initialize mock sql db
	mockDb, mocking, err := sqlmock.New()
	defer mockDb.Close()
	if err != nil {
		log.Fatalf("Failed to open mock sql db, got error %v", err)
	}
	mock = mocking

	// initialize db connection
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDb}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database, got error %v", err)
	}

	// initialize repository
	repository = NewUserPostgresRepository(db)

	code := m.Run()

	fmt.Println("=================================================================")
	fmt.Println("TEST END")
	fmt.Println("=================================================================")
	os.Exit(code)
}

// Unit Test User Register
func TestUserPostgresRepository_Create(t *testing.T) {
	// prepare model data for user register
	usr := &model.User{
		Information: model.Information{
			FirstName: "muhammat",
			LastName:  "saputra",
		},
		Email:    "mrizkisaputra6@gmail.com",
		Password: "secret",
	}

	// define expected sql query
	mock.ExpectBegin()
	mock.ExpectQuery(expectedCreateUserQuery).
		WithArgs(usr.Email, usr.Password, usr.Information.FirstName, usr.Information.LastName, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
	mock.ExpectCommit()

	createdUser, err := repository.Create(context.Background(), usr)
	require.Nil(t, err)
	require.NotNil(t, createdUser)
	//fmt.Printf("%+v", createdUser)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

// Unit Test Updating existing user
func TestUserPostgresRepository_Update(t *testing.T) {
	// prepare model data update user
	usr := &model.User{
		Information: model.Information{
			FirstName:   "Jhon",
			LastName:    "Smith",
			City:        "New York",
			PhoneNumber: "123456",
		},
		Email:    "jhon@smith.com",
		Password: "secret",
		Avatar:   "https://example.com/images",
	}

	// define expected sql query
	mock.ExpectBegin()
	mock.ExpectExec(expectedUpdateUserQuery).
		WithArgs(
			usr.Email,
			usr.Password,
			usr.Avatar,
			usr.Information.FirstName,
			usr.Information.LastName,
			usr.Information.City,
			usr.Information.PhoneNumber,
			sqlmock.AnyArg(),
			sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	updatedUser, err := repository.Update(context.Background(), usr)
	require.Nil(t, err)
	require.NotNil(t, updatedUser)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

// Unit Test Find user by id
func TestUserPostgresRepository_FindById(t *testing.T) {
	testUser := &model.User{
		Id:       uuid.New(),
		Email:    "Alex@gmail.com",
		Password: "secret",
		Information: model.Information{
			FirstName:   "alex",
			LastName:    "smith",
			City:        "New York",
			PhoneNumber: "123456",
		},
		CreatedAt: time.Now().UnixMilli(),
		UpdatedAt: time.Now().UnixMilli(),
	}

	// define expected sql query
	rows := sqlmock.NewRows([]string{"id", "email", "password", "avatar", "first_name", "last_name", "city", "phone_number", "created_at", "updated_at"}).
		AddRow(testUser.Id, testUser.Email, testUser.Password, testUser.Avatar, "alex", "smith", "New York", "123456", testUser.CreatedAt, testUser.UpdatedAt)
	mock.ExpectQuery(expectedFindUserByIdQuery).WithArgs(testUser.Id, 1).WillReturnRows(rows)

	findUser, err := repository.FindById(context.Background(), testUser)
	require.Nil(t, err)
	require.NotNil(t, findUser)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestUserPostgresRepository_FindAlreadyExistByEmail(t *testing.T) {
	var expectCount int64 = 1
	testUser := &model.User{
		Email: "jonson@gmail.com",
	}

	// define expected sql query
	mock.ExpectQuery(expectedFindAlreadyEmail).WithArgs(testUser.Email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectCount))

	total, err := repository.FindAlreadyExistByEmail(context.Background(), testUser)
	require.Nil(t, err)
	require.Equal(t, expectCount, total)
	require.NotEqual(t, 0, total)
	require.NoError(t, mock.ExpectationsWereMet())
}
