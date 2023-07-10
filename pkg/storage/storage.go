package storage

type Post struct {
	//ID      primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Title   string `json:"title" bson:"title"`
	Content string `json:"content" bson:"content"`
	PubTime int64  `json:"published_at" bson:"published_at"`
	Link    string `json:"link" bson:"link"`
}

type Interface interface {
	Posts(n int) ([]Post, error) // получение n последних публикаций
	//Posts() ([]Post, error) // получение всех публикаций
	AddPost(Post) error    // создание новой публикации
	UpdatePost(Post) error // обновление публикации
	DeletePost(Post) error // удаление публикации по ID
}
