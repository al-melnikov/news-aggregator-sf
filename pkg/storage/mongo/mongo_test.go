package mongo

import (
	"context"
	"news-aggregator-sf/pkg/storage"
	"testing"
	"time"
)

// TODO надо бы как-нибудь запускать тестовую базу а не ту же самую
func TestNew(t *testing.T) {
	ctx := context.Background()
	_, err := New(ctx, "mongodb://localhost:27017/")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStorage_AddPost(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := New(ctx, "mongodb://localhost:27017/")
	post := storage.Post{
		Title:   "тест",
		Content: "Текст",
		PubTime: 1232,
		Link:    "ссылка",
	}

	db.AddPost(post)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("addPost tested")
}

func TestStorage_UpdatePost(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := New(ctx, "mongodb://localhost:27017/")
	post := storage.Post{
		Title:   "тест",
		Content: "новый текст",
		PubTime: 1232,
		Link:    "ссылка",
	}
	db.UpdatePost(post)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("updatePost tested")
}

func TestStorage_DeletePost(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := New(ctx, "mongodb://localhost:27017/")
	post := storage.Post{
		Title:   "тест",
		Content: "новый текст",
		PubTime: 1232,
		Link:    "ссылка",
	}
	db.DeletePost(post)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("deletePost tested")
}
