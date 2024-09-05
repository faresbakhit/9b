package store

import (
	"fmt"
	"strings"
	"time"
)

type UserPost struct {
	Id        int64
	UserId    int64
	Title     string
	Url       string
	Body      string
	Score     int
	CreatedAt *time.Time
}

func (s *Store) UserPostNew(post *UserPost) error {
	post.Body = strings.ReplaceAll(post.Body, "\r\n", "\n")
	query := `
		INSERT INTO user_post
		(user_id, title, url, body)
		VALUES (?, ?, ?, ?)
		RETURNING id, score, created_at`
	row := s.db.QueryRow(query, post.UserId, post.Title, post.Url, post.Body)
	if err := row.Scan(&post.Id, &post.Score, &post.CreatedAt); err != nil {
		return err
	}
	return nil
}

func (s *Store) UserPostListFromUser(userId int64, limit, offset int) func(func(*UserPost, error) bool) {
	return func(yield func(*UserPost, error) bool) {
		query := fmt.Sprintf(`
			SELECT * FROM user_post
			ORDER BY created_at
			LIMIT %d OFFSET %d`, limit, offset)
		rows, err := s.db.Query(query)
		defer rows.Close()
		if err != nil {
			yield(nil, err)
			return
		}
		for rows.Next() {
			var post UserPost
			err := rows.Scan(
				&post.Id,
				&post.UserId,
				&post.Title,
				&post.Url,
				&post.Body,
				&post.Score,
				&post.CreatedAt)
			if err != nil {
				yield(nil, err)
				return
			}
			if !yield(&post, nil) {
				return
			}
		}
		if err := rows.Err(); err != nil {
			yield(nil, err)
		}
	}
}
