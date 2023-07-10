package mongo

import (
	"context"
	"errors"
	"log"
	"news-aggregator-sf/pkg/storage"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DatabaseName string = "go_news"
var CollectionName string = "posts"

type Storage struct {
	Db *mongo.Client
}

// Возвращает новый экземпляр базы данных
func New(ctx context.Context, constr string) (*Storage, error) {
	mongoOpts := options.Client().ApplyURI(constr)
	client, err := mongo.Connect(ctx, mongoOpts)
	if err != nil {
		log.Fatal(err)
	}
	// не забываем закрывать ресурсы
	s := Storage{
		Db: client,
	}
	return &s, nil
}

// Возвращает n последних постов из базы данных
func (mg *Storage) Posts(n int) ([]storage.Post, error) {
	collection := mg.Db.Database(DatabaseName).Collection(CollectionName)
	filter := bson.D{}
	opts := options.Find().SetSort(bson.D{{Key: "published_at", Value: -1}}).SetLimit(int64(n))
	cur, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		panic(err)
	}

	var results []storage.Post
	if err = cur.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return results, nil
}

// Добавляет новый пост в базу данных
func (mg *Storage) AddPost(p storage.Post) error {
	collection := mg.Db.Database(DatabaseName).Collection(CollectionName)

	filter := bson.D{{Key: "title", Value: p.Title}}
	var post storage.Post
	err := collection.FindOne(context.TODO(), filter).Decode(&post)
	if err == nil {
		return errors.New("post already exists")
	}

	_, err = collection.InsertOne(context.Background(), p)
	if err != nil {
		return err
	}
	return nil
}

// Удаляет пост из базы данных с совпадающим title
func (mg *Storage) DeletePost(p storage.Post) error {
	collection := mg.Db.Database(DatabaseName).Collection(CollectionName)
	filter := bson.D{{Key: "title", Value: p.Title}}
	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

// Обновляет content поста с совпадающим title
func (mg *Storage) UpdatePost(p storage.Post) error {
	collection := mg.Db.Database(DatabaseName).Collection(CollectionName)
	filter := bson.D{{Key: "title", Value: p.Title}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "content", Value: p.Content}}}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
