package repository

import (
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/mrizkisaputra/expenses-api/internal/expense"
	"github.com/mrizkisaputra/expenses-api/internal/expense/model"
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
	repository expense.PostgresRepository
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
	repository = NewExpensePgRepository(db)

	code := m.Run()

	fmt.Println("=================================================================")
	fmt.Println("TEST END")
	fmt.Println("=================================================================")
	os.Exit(code)
}

func TestExpensePostgresRepository_Create(t *testing.T) {
	userId := uuid.New()
	amount := 100.00
	e := &model.Expense{
		UserId:      userId,
		Description: "dinner with family",
		Amount:      &amount,
		Category:    "Dinner",
		CreatedAt:   time.Now().UnixMilli(),
		UpdatedAt:   time.Now().UnixMilli(),
	}

	// define expect sql
	expectedSQLCreateExpense := regexp.QuoteMeta(`INSERT INTO "expenses" ("id_user","description","amount","category","created_at","updated_at","deleted_at")
	VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)
	mock.ExpectBegin()
	mock.ExpectQuery(expectedSQLCreateExpense).
		WithArgs(e.UserId, e.Description, e.Amount, e.Category, e.CreatedAt, e.UpdatedAt, e.DeletedAt).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
	mock.ExpectCommit()

	err := repository.Create(context.Background(), e)
	require.Nil(t, err)
	require.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestExpensePostgresRepository_FindByIdAndUserId(t *testing.T) {
	userId := uuid.New()
	id := uuid.New()
	expenses := new(model.Expense)

	// define expect sql
	expectedSQLFindByIdAnUserId := regexp.QuoteMeta(`SELECT * FROM "expenses" WHERE (id = $1 AND id_user = $2) AND "expenses"."deleted_at" IS NULL LIMIT $3`)
	rows := sqlmock.NewRows([]string{"id", "id_user", "description", "amount", "category", "created_at", "updated_at", "deleted_at"}).
		AddRow(id, userId, "dinner with family", 100.00, "dinner", time.Now().UnixMilli(), time.Now().UnixMilli(), nil)
	mock.ExpectQuery(expectedSQLFindByIdAnUserId).WithArgs(id, userId, 1).
		WillReturnRows(rows)

	err := repository.FindByIdAndUserId(context.Background(), expenses, id.String(), userId.String())
	require.Nil(t, err)
	require.NoError(t, err)
	require.NotNil(t, expenses)
	require.Equal(t, id, expenses.Id)
	require.Equal(t, userId, expenses.UserId)
	require.Equal(t, "dinner with family", expenses.Description)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestExpensePostgresRepository_Remove(t *testing.T) {
	userId := uuid.New()
	id := uuid.New()
	amount := 20.00
	expenses := &model.Expense{
		Id:          id,
		UserId:      userId,
		Description: "Breakfast with family",
		Amount:      &amount,
		Category:    "Breakfast",
	}

	// define expect sql
	expectedSQLRemove := regexp.QuoteMeta(`UPDATE "expenses" SET "deleted_at"=$1 WHERE "expenses"."id" = $2 AND "expenses"."deleted_at" IS NULL`)
	mock.ExpectBegin()
	mock.ExpectExec(expectedSQLRemove).WithArgs(sqlmock.AnyArg(), id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repository.Remove(context.Background(), expenses)
	require.Nil(t, err)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations, got error: %v", err)
	}
}

func TestExpensePostgresRepository_Update(t *testing.T) {

}

func TestExpensePostgresRepository_FindAll(t *testing.T) {

}

func TestExpensePostgresRepository_FindAllByDateRange(t *testing.T) {

}
