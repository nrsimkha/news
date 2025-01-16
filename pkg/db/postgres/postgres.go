package postgres

import (
	"context"
	storage "goNews/pkg/db"

	"github.com/jackc/pgx/v5"
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

// Получить массив новостей из БД
// keystring - строка для фильтра по названию новости. Если пустая, то возвращаются все новости
// Пагинация осуществляется переменной номера страницы page
func (store Store) News(page int, limit int, keystring string) (*storage.PostsWithPagination, error) {
	var pagination storage.Pagination
	pagination.CurrentPage = page
	pagination.PostsPerPage = limit
	//Общее количество страниц
	var totalCount int
	var rows pgx.Rows
	var err error
	// Calculate the OFFSET
	offset := (page - 1) * limit

	if keystring != "" {
		store.db.QueryRow(ctx, "SELECT count(*) FROM news WHERE title ILIKE $1", "%"+keystring+"%").Scan(&totalCount)
		rows, err = store.db.Query(ctx, `
		SELECT * FROM news WHERE title ILIKE $1 ORDER BY id LIMIT $2 OFFSET $3; 
	`, "%"+keystring+"%", limit, offset)
	} else {
		store.db.QueryRow(ctx, "SELECT count(*) FROM news").Scan(&totalCount)
		rows, err = store.db.Query(ctx, `
		SELECT * FROM news ORDER BY id LIMIT $1 OFFSET $2; 
	`, limit, offset)
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	pagination.TotalPages = totalCount / pagination.PostsPerPage
	var posts []storage.Post
	defer rows.Close()
	posts, err = pgx.CollectRows(rows, pgx.RowToStructByName[storage.Post])
	if err != nil {
		return nil, err
	}
	/* 	for rows.Next() {
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

	} */
	return &storage.PostsWithPagination{
		Posts:      posts,
		Pagination: pagination,
	}, rows.Err()
}

// Получить новость из БД по ID.
func (store Store) NewsByID(id int) (*storage.Post, error) {

	row := store.db.QueryRow(ctx, `
	SELECT * FROM news WHERE id=$1; 
`, id)
	var post storage.Post
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.PubTime,
		&post.Link,
	)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

// Добавить массив новостей в БД.
// Возвращает последний ID, добавленный в таблицу.
func (store Store) AddNews(news []storage.Post) (*int, error) {
	var idAr []int
	var id int
	for _, post := range news {
		row := store.db.QueryRow(ctx, `INSERT INTO news (Title, Content, PubTime, Link) VALUES ($1, $2, $3, $4) RETURNING id;`,
			post.Title, post.Content, post.PubTime, post.Link)

		if err := row.Scan(&id); err != nil { // scan will release the connection

			return nil, err
		}
		idAr = append(idAr, id)
	}

	return &idAr[len(idAr)-1], nil
}
