package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"news-aggregator-sf/pkg/api"
	"news-aggregator-sf/pkg/rss"
	"news-aggregator-sf/pkg/storage"
	"news-aggregator-sf/pkg/storage/mongo"
	"os"
	"time"
)

// конфигурация приложения
type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"request_period"`
}

func main() {

	mongo.DatabaseName = "go_news"
	mongo.CollectionName = "posts"

	// инициализация зависимостей приложения
	db, err := mongo.New(context.Background(), "mongodb://localhost:27017/")
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)

	// чтение и раскодирование файла конфигурации
	b, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}

	// запуск парсинга новостей в отдельном потоке
	// для каждой ссылки
	chPosts := make(chan []storage.Post)
	chErrs := make(chan error)
	for _, url := range config.URLS {
		go parseURL(url, db, chPosts, chErrs, config.Period)
	}

	// запись потока новостей в БД
	go func() {
		for posts := range chPosts {
			for _, post := range posts {
				db.AddPost(post)
			}
		}
	}()

	// обработка потока ошибок
	go func() {
		for err := range chErrs {
			log.Println("ошибка:", err)
		}
	}()

	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(":4567", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

// Асинхронное чтение потока RSS. Раскодированные
// новости и ошибки пишутся в каналы.
func parseURL(url string, db *mongo.Storage, posts chan<- []storage.Post, errs chan<- error, period int) {
	for {
		news, err := rss.RssToStruct(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
		time.Sleep(time.Minute * time.Duration(period))
	}
}
