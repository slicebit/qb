package qb

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestLogger(t *testing.T) {
// 	// engine, err := New("default", ":memory:")
// 	// metadata := MetaData()
// 	// actors := Table("actors",
// 	// 	Column("id", BigInt()).NotNull(),
// 	// 	PrimaryKey("id"),
// 	// )
// 	// metadata.AddTable(actors)
// 	// metadata.CreateAll(engine)
// 	defer metadata.DropAll(engine)
// 	logCapture := &TestingLogWriter{t, nil}
// 	defer logCapture.Flush()
// 	engine.SetLogger(&DefaultLogger{LQuery | LBindings, log.New(logCapture, "", log.LstdFlags)})
// 	engine.Logger().SetLogFlags(LQuery)

// 	_, err = engine.Exec(actors.Insert().Values(map[string]interface{}{"id": 5}))
// 	assert.Nil(t, err)

// 	engine.Logger().SetLogFlags(LQuery | LBindings)
// 	_, err = engine.Exec(actors.Insert().Values(map[string]interface{}{"id": 10}))
// 	assert.Nil(t, err)

// 	assert.Equal(t, engine.Logger().LogFlags(), LQuery|LBindings)
// }

func TestLoggerFlags(t *testing.T) {
	logger := DefaultLogger{LDefault, log.New(os.Stdout, "", -1)}

	logger.SetLogFlags(LBindings)

	assert.Equal(t, logger.LogFlags(), LBindings)
	// engine, err := New("sqlite3", ":memory:")
	// assert.Equal(t, nil, err)

	// before setting flags, this is on the default
	// assert.Equal(t, engine.Logger().LogFlags(), LDefault)

	// engine.SetLogFlags(LBindings)
	// after setting flags, we have the expected value
	// assert.Equal(t, engine.Logger().LogFlags(), LBindings)
}
