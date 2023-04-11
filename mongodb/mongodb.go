package mongodb

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/minhnvc/go_corelib/utils"
)

type MongoDBClient struct {
	client *mongo.Client
}

func NewMongoDBClient() (*MongoDBClient, error) {
	// Setup the mgm default config
	clientOptions := options.Client().ApplyURI(utils.GetConfig("MONGO_URL"))
	clientOptions.SetMaxPoolSize(100)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	fmt.Println("Mongo DB connected")
	// Ping the primary
	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}
	return &MongoDBClient{client: client}, nil
}

func (c *MongoDBClient) Close() {
	if c.client != nil {
		c.client.Disconnect(context.Background())
	}
}

func (c *MongoDBClient) GetCollection(name string) *mongo.Collection {
	return c.client.Database(utils.GetConfig("DATABASE_NAME")).Collection(name)
}

type ModelInterface interface {
	Save() string
}

func getId(model interface{}) primitive.ObjectID {
	objectId, _ := primitive.ObjectIDFromHex(reflect.Indirect(reflect.ValueOf(model)).FieldByName("Id").String())
	return objectId
}

func getObjectIdFromId(id string) primitive.ObjectID {
	objectId, _ := primitive.ObjectIDFromHex(id)
	return objectId
}

// public functions
func InsertOne(collection string, model interface{}) (string, error) {
	db, err := NewMongoDBClient()
	if err != nil {
		return "", err
	}
	defer db.Close()
	result, err := db.GetCollection(collection).InsertOne(context.Background(), model)

	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func Update(collection string, model interface{}) (string, error) {

	db, err := NewMongoDBClient()
	if err != nil {
		return "", err
	}
	defer db.Close()

	filter := bson.M{"_id": getId(model)}

	//remove _id field before update
	reflect.Indirect(reflect.ValueOf(model)).FieldByName("Id").SetString("")

	updates, _ := utils.ToDoc(model)
	_, err = db.GetCollection(collection).UpdateOne(context.Background(), filter,
		bson.D{{Key: "$set", Value: updates}})

	if err != nil {
		return "", err
	}

	return getId(model).Hex(), nil
}

func UpdateOne(collection string, id string, data bson.D) (string, error) {
	db, err := NewMongoDBClient()
	if err != nil {
		return "", err
	}
	defer db.Close()

	filter := bson.M{"_id": getObjectIdFromId(id)}

	_, err = db.GetCollection(collection).UpdateOne(context.Background(), filter,
		bson.D{{Key: "$set", Value: data}})

	if err != nil {
		return "", err
	}
	return id, nil
}

func GetById(collection string, model interface{}, id string) {
	db, _ := NewMongoDBClient()
	defer db.Close()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}

	result := db.GetCollection(collection).FindOne(context.Background(), filter)

	if result != nil {
		err := result.Decode(model)
		if err != nil {
			fmt.Println("Mongo", err, collection, id)
		}
	}
}

func Count(collection string, filter primitive.M) int64 {
	db, _ := NewMongoDBClient()
	defer db.Close()

	result, _ := db.GetCollection(collection).CountDocuments(context.Background(), filter)
	return result
}

func FindOne(collection string, results interface{}, filter primitive.M) {
	db, _ := NewMongoDBClient()
	defer db.Close()

	result := db.GetCollection(collection).FindOne(context.Background(), filter)
	if result != nil {
		err := result.Decode(results)
		if err != nil {
			fmt.Println("Mongo", err, collection, filter)
		}
	}
}
func Find(collection string, results interface{}, filter primitive.M) {
	db, _ := NewMongoDBClient()
	defer db.Close()

	result, _ := db.GetCollection(collection).Find(context.Background(), filter)
	defer result.Close(context.Background())
	if result != nil {
		err := result.All(context.Background(), results)
		if err != nil {
			fmt.Println("Mongo", err)
		}
	}
}

func All(collection string, results interface{}) {
	db, _ := NewMongoDBClient()
	defer db.Close()

	result, _ := db.GetCollection(collection).Find(context.Background(), bson.M{})
	defer result.Close(context.Background())
	if result != nil {
		err := result.All(context.Background(), results)
		if err != nil {
			fmt.Println("Mongo", err)
		}
	}
}

func Aggregate(collection string, results interface{}, filter []primitive.D) {
	db, _ := NewMongoDBClient()
	defer db.Close()

	result, _ := db.GetCollection(collection).Aggregate(context.Background(), filter)
	defer result.Close(context.Background())
	if result != nil {
		err := result.All(context.Background(), results)
		if err != nil {
			fmt.Println("Mongo", err.Error())
		}
	}
}

func CreateIndex(collectionName string, field string, value string, unique bool) bool {
	db, _ := NewMongoDBClient()
	defer db.Close()

	mod := mongo.IndexModel{
		Keys:    bson.M{field: value}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(unique),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.GetCollection(collectionName)

	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		fmt.Println("Mongo", err.Error())
		return false
	}

	return true
}
