package qb

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestLogger(t *testing.T) {
	db, err := New("sqlite3", ":memory:")
	actors := Table("actors",
		Column("id", BigInt()).NotNull(),
		PrimaryKey("id"),
	)
	db.Metadata().AddTable(actors)
	db.CreateAll()
	defer db.DropAll()
	db.Engine().SetLogger(DefaultLogger{LQuery | LBindings, log.New(ioutil.Discard, "", log.LstdFlags)})
	db.Engine().Logger().SetLogFlags(LQuery)

	_, err = db.Engine().Exec(actors.Insert().Values(map[string]interface{}{"id": 5}))
	assert.Nil(t, err)

	db.Engine().Logger().SetLogFlags(LBindings)
	_, err = db.Engine().Exec(actors.Insert().Values(map[string]interface{}{"id": 10}))
	assert.Nil(t, err)

	assert.Equal(t, LQuery|LBindings, db.Engine().Logger().LogFlags())
}
