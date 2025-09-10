package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (content, title, user_id, tags)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostStore) GetById(ctx context.Context, id int64) (*Post, error) {
	var post Post
	query := `
	SELECT id,title,content,created_at,updated_at,user_id, tags FROM posts WHERE id = $1
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.UserID,
		pq.Array(&post.Tags),
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *PostStore) GetAllPosts(ctx context.Context) ([]Post, error) {
	query := `
	SELECT p.id,p.title,p.content,p.created_at,p.updated_at,p.user_id, p.tags FROM posts p
	ORDER BY p.id
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	posts := []Post{}

	for rows.Next() {
		post := Post{}

		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt, &post.UserID, pq.Array(&post.Tags))

		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (s *PostStore) DeleteById(ctx context.Context, id int64) error {
	query := `
	DELETE FROM posts WHERE id=$1
	`
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
func (s *PostStore) UpdatePost(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts 
	SET title=$1, content=$2 
	WHERE id=$3
	RETURNING id, user_id, created_at, updated_at, tags
	`
	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID).Scan(&post.ID, &post.UserID, &post.CreatedAt, &post.UpdatedAt, pq.Array(&post.Tags))
	if err != nil {
		return err
	}
	return nil
}
