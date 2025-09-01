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
	store := NewStore(t)

	got, gotErr := store.Get(789)

	require.NoError(t, gotErr)
	require.Empty(t, got)
}

func TestSQLStore_Get_ItemPresentWhenSetBefore(t *testing.T) {
	store := NewStore(t)
	err := store.Set(789, task.TasksText("test text"))
	require.NoError(t, err)

	got, gotErr := store.Get(789)

	require.NoError(t, gotErr)
	require.EqualValues(t, "test text", got)
}

func TestSQLStore_Set_ItemOverridenWhenSetBefore(t *testing.T) {
	store := NewStore(t)
	err := store.Set(789, task.TasksText("test text"))
	require.NoError(t, err)

	gotSetErr := store.Set(789, task.TasksText("test text update"))
	require.NoError(t, gotSetErr)
	got, gotErr := store.Get(789)

	require.NoError(t, gotErr)
	require.EqualValues(t, "test text update", got)
}

func TestSQLStore_Set_OkWhenDBEmpty(t *testing.T) {
	store := NewStore(t)

	gotErr := store.Set(7, task.TasksText("task A"))

	require.NoError(t, gotErr)
}

func TestSQLStore_Close(t *testing.T) {
	db := NewTestDB(t)
	store, gotErr := store.NewSQLStore(db)
	require.NoError(t, gotErr)

	gotErr = store.Close()

	require.NoError(t, gotErr)
}

func NewStore(t *testing.T) *store.SQLStore {
	t.Helper()

	db := NewTestDB(t)
	store, gotErr := store.NewSQLStore(db)
	require.NoError(t, gotErr)

	t.Cleanup(func() {
		closeErr := store.Close()
		require.NoError(t, closeErr, "close store")
	})

	return store
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
