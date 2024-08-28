package redkademo

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nalgeon/redka"
	"github.com/stretchr/testify/assert"
)

/*
单实例场景下，不必额外安装 redis，可直接使用 redis 的数据结构
*/
func TestRedka(t *testing.T) {
	// Open a database.
	db, err := redka.Open("data.db", nil)
	assert.Nil(t, err)
	defer db.Close()

	// Set some string keys.
	err = db.Str().Set("name", "alice")
	assert.Nil(t, err)

	err = db.Str().Set("age", 25)
	assert.Nil(t, err)

	// Check if the keys exist.
	count, err := db.Key().Count("name", "age", "city")
	assert.Nil(t, err)
	assert.Equal(t, 2, count)

	// Get a key.
	name, err := db.Str().Get("name")
	assert.Nil(t, err)
	assert.Equal(t, "alice", name.String())

	_, err = db.Hash().Set("hs", "one", 1)
	assert.Nil(t, err)
	_, err = db.Hash().Set("hs", "two", 2)
	assert.Nil(t, err)
	_, err = db.Hash().Set("hs", "three", 3)
	assert.Nil(t, err)
	_, err = db.Hash().Set("hs", "four", 4)
	assert.Nil(t, err)

	count, err = db.Hash().Len("hs")
	assert.Nil(t, err)
	assert.Equal(t, 4, count)
}
