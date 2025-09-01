package store_test

import (
	"database/sql"
	"testing"

	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/store"
	"github.com/IlyasYOY/tasks-assistant-tg-bot/internal/task"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestSQLStore_Get_EmptyWhenDBEmpty(t *testing.T) {
	s := NewStore(t)

	got, gotErr := s.Get(789)

	require.NoError(t, gotErr)
	require.Empty(t, got)
}

func TestSQLStore_Get_ItemPresentWhenSetBefore(t *testing.T) {
	s := NewStore(t)
	err := s.Set(789, task.TasksText("test text"))
	require.NoError(t, err)

	got, gotErr := s.Get(789)

	require.NoError(t, gotErr)
	require.EqualValues(t, "test text", got)
}

func TestSQLStore_Set_ItemOverridenWhenSetBefore(t *testing.T) {
	s := NewStore(t)
	err := s.Set(789, task.TasksText("test text"))
	require.NoError(t, err)

	gotSetErr := s.Set(789, task.TasksText("test text update"))
	require.NoError(t, gotSetErr)
	got, gotErr := s.Get(789)

	require.NoError(t, gotErr)
	require.EqualValues(t, "test text update", got)
}

func TestSQLStore_Set_OkWhenDBEmpty(t *testing.T) {
	s := NewStore(t)

	gotErr := s.Set(7, task.TasksText("task A"))

	require.NoError(t, gotErr)
}

func NewStore(t *testing.T) *store.SQLStore {
	t.Helper()

	db := NewTestDB(t)
	s, newErr := store.NewSQLStore(db)
	require.NoError(t, newErr, "create sql store")

	t.Cleanup(func() {
		closeErr := s.Close()
		require.NoError(t, closeErr, "close store")
	})

	return s
}

func NewTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, openErr := sql.Open("sqlite", "file::memory:?cache=shared")
	require.NoError(t, openErr, "open in‑memory sqlite")

	migrationErr := store.Migrate(db, "migrations")
	require.NoError(t, migrationErr, "run migrations")

	t.Cleanup(func() {
		closeErr := db.Close()
		require.NoError(t, closeErr, "close db")
	})
	return db
}
