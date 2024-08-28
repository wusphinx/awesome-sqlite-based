package ferretdbdemo

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/FerretDB/FerretDB/ferretdb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName   = "one"
	collName = "colone"
)

type Tea struct {
	Item        string    `bson:"item,omitempty"`
	Rating      int       `bson:"rating,omitempty"`
	DateOrdered time.Time `bson:"date_ordered,omitempty"`
}

func TestEmbedded(t *testing.T) {
	socketPath := "/tmp/dummy.sock"
	// 删除已存在的 socket 文件
	os.Remove(socketPath)

	ctx, cancel := context.WithCancel(context.Background())

	f, err := ferretdb.New(&ferretdb.Config{
		Listener: ferretdb.ListenerConfig{
			Unix: socketPath,
		},
		Handler:   "sqlite",
		SQLiteURL: "file:./",
	})
	require.NoError(t, err)

	done := make(chan struct{})

	go func() {
		require.NoError(t, f.Run(ctx))
		close(done)
	}()

	uri := f.MongoDBURI()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.NoError(t, err)

	err = client.Ping(ctx, nil)
	require.NoError(t, err)

	uv7, err := uuid.NewV7()
	require.NoError(t, err)
	item := uv7.String()
	rating := 10
	now := time.Now()
	coll := client.Database(dbName).Collection(collName)

	// 增
	result, err := coll.InsertOne(ctx, Tea{Item: item, Rating: rating, DateOrdered: now})
	require.NoError(t, err)
	t.Logf("result is: %v", result.InsertedID)

	// 查
	filter := bson.D{{"_id", result.InsertedID}}
	var data Tea
	err = coll.FindOne(ctx, filter).Decode(&data)
	require.NoError(t, err)
	assert.Equal(t, item, data.Item)
	t.Logf("sr is: %+v", data)

	// 改
	rating2 := rating + 1
	update := bson.D{{"$set", bson.D{{"rating", rating2}}}}
	_, err = coll.UpdateByID(ctx, result.InsertedID, update)
	assert.Nil(t, err)

	// 删
	deleteResult, err := coll.DeleteOne(ctx, filter)
	assert.Nil(t, err)

	assert.EqualValues(t, 1, deleteResult.DeletedCount)

	cancel()
	<-done
}
