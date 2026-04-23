package database

import (
	"database/sql"
	"os"
	"testing"

	"github.com/douglasvolcato/binary-code-processor/task_service/test"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func NewTestDB() (*sql.DB, error) {
	file := "./test.db"
	os.Remove(file)

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE tasks (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT FALSE
	);
	`
	tx, err := OpenTransaction(db)
	if err != nil {
		return nil, err
	}
	if err := ExecuteSQL(tx, createTableSQL); err != nil {
		return nil, err
	}
	if err := CommitTransaction(tx); err != nil {
		return nil, err
	}

	return db, nil
}

func TestShouldExecuteCommit(t *testing.T) {
	db, err := NewTestDB()
	assert.NoError(t, err)
	defer db.Close()

	faker := test.FakeData{}
	id := faker.ID()
	title := faker.Phrase()

	tx, err := OpenTransaction(db)
	assert.NoError(t, err)

	insertSQL := `INSERT INTO tasks (id, title, completed) VALUES (?, ?, ?)`
	err = ExecuteSQL(tx, insertSQL, id, title, false)
	assert.NoError(t, err)

	err = CommitTransaction(tx)
	assert.NoError(t, err)

	querySQL := `SELECT id, title, completed FROM tasks WHERE id = ?`
	rows, err := QuerySQL(db, querySQL, id)
	assert.NoError(t, err)
	defer rows.Close()

	var retrievedID, retrievedTitle string
	var retrievedCompleted bool
	count := 0
	for rows.Next() {
		count++
		err := rows.Scan(&retrievedID, &retrievedTitle, &retrievedCompleted)
		assert.NoError(t, err)
		assert.Equal(t, id, retrievedID)
		assert.Equal(t, title, retrievedTitle)
		assert.Equal(t, false, retrievedCompleted)
	}
	assert.Equal(t, 1, count)
}

func TestShouldExecuteRollback(t *testing.T) {
	db, err := NewTestDB()
	assert.NoError(t, err)
	defer db.Close()

	faker := test.FakeData{}
	id := faker.ID()
	title := faker.Phrase()

	tx, err := OpenTransaction(db)
	assert.NoError(t, err)

	insertSQL := `INSERT INTO tasks (id, title, completed) VALUES (?, ?, ?)`
	err = ExecuteSQL(tx, insertSQL, id, title, false)
	assert.NoError(t, err)

	err = RollbackTransaction(tx)
	assert.NoError(t, err)

	querySQL := `SELECT id, title, completed FROM tasks WHERE id = ?`
	rows, err := QuerySQL(db, querySQL, id)
	assert.NoError(t, err)
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}
	assert.Equal(t, 0, count)
}

func TestShouldReturnErrorOnInvalidSQL(t *testing.T) {
	db, err := NewTestDB()
	assert.NoError(t, err)
	defer db.Close()

	faker := test.FakeData{}
	id := faker.ID()

	tx, err := OpenTransaction(db)
	assert.NoError(t, err)

	invalidSQL := `INSERT INTO non_existing_table (id) VALUES (?)`
	err = ExecuteSQL(tx, invalidSQL, id)
	assert.Error(t, err)
}

func TestShouldDeleteTestDB(t *testing.T) {
	file := "./test.db"
	os.Remove(file)
}
