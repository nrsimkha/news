package postgres

import (
	"context"
	storage "goNews/pkg/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

var ctx context.Context = context.Background()

// Конструктор БД.
func New(constr string) (*Store, error) {
	db, err := pgxpool.New(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}
	return &s, nil
}

// Получить массив новостей из БД в количестве amount
func (store Store) News(amount int) ([]storage.Post, error) {

	rows, err := store.db.Query(ctx, `
	SELECT * FROM news ORDER BY id LIMIT $1; 
`, amount)
	if err != nil {
		return nil, err
	}
	var posts []storage.Post

	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)

	}
	return posts, rows.Err()
}

// Добавить массив новостей в БД
func (store Store) AddNews(news []storage.Post) error {
	_, err := store.db.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS news (
			ID SERIAL PRIMARY KEY, 
			Title TEXT NOT NULL,
			Content TEXT NOT NULL,
			PubTime BIGINT NOT NULL DEFAULT extract(epoch from now()),
			Link TEXT NOT NULL
		);
	`)
	if err != nil {
		return err
	}
	for _, post := range news {
		_, err = store.db.Exec(ctx, `INSERT INTO news (Title, Content, PubTime, Link) VALUES ($1, $2, $3, $4);`,
			post.Title, post.Content, post.PubTime, post.Link)
		if err != nil {
			return err
		}
	}

	return nil
}
