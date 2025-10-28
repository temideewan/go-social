package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(ctx context.Context, post *Post) error
		GetById(ctx context.Context, id int64) (*Post, error)
		DeleteById(ctx context.Context, id int64) error
		GetAllPosts(ctx context.Context) ([]Post, error)
		UpdatePost(ctx context.Context, post *Post) error
		GetUserFeed(ctx context.Context, userId int64, query PaginatedFeedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(ctx context.Context, user *User) error
		GetById(ctx context.Context, id int64) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string) error
	}
	Comments interface {
		GetByPostId(ctx context.Context, id int64) ([]Comment, error)
		Create(ctx context.Context, comment *Comment) error
	}

	Followers interface {
		Follow(ctx context.Context, followerId int64, followeeId int64) error
		Unfollow(ctx context.Context, followerId int64, followeeId int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}
