package store

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type UserPostNew struct {
	UserId int
	Title  string
	URL    string
	Body   string
}

func (s *Store) UserPostNew(post *UserPostNew) error {
	post.Body = strings.ReplaceAll(post.Body, "\r\n", "\n")
	query := `
		INSERT INTO user_post
		(user_id, title, url, body)
		VALUES (?, ?, ?, ?)`
	_, err := s.db.Exec(query, post.UserId, post.Title, post.URL, post.Body)
	return err
}

type UserPostGet struct {
	Id            int
	Score         int
	UserUpvoted   bool
	UserDownvoted bool
	Author        string
	Title         string
	URL           string
	Body          string
	Date          *time.Time
}

func (s *Store) UserPostGet(userId, postId int) *UserPostGet {
	var post UserPostGet
	query := fmt.Sprintf(`
		SELECT
			user_post.id,
			score,
			IIF(user_post_upvote.id, 1, 0),
			IIF(user_post_downvote.id, 1, 0),
			username,
			title,
			url,
			body,
			user_post.created_at
		FROM user_post
		JOIN user ON user.id = user_post.user_id
		LEFT JOIN user_post_upvote ON
			user_post_upvote.post_id = user_post.id
			AND user_post_upvote.user_id = %[1]d
		LEFT JOIN user_post_downvote ON
			user_post_downvote.post_id = user_post.id
			AND user_post_downvote.user_id = %[1]d
		WHERE user_post.id = %d`, userId, postId)
	row := s.db.QueryRow(query)
	if err := row.Scan(
		&post.Id,
		&post.Score,
		&post.UserUpvoted,
		&post.UserDownvoted,
		&post.Author,
		&post.Title,
		&post.URL,
		&post.Body,
		&post.Date); err != nil {
		log.Printf("UserPostGet: %v", err)
		return nil
	}
	return &post
}

type UserPostListResult struct {
	Id            int
	Score         int
	Comments      int
	UserUpvoted   bool
	UserDownvoted bool
	Author        string
	Title         string
	URL           string
	Body          string
	Date          *time.Time
}

func (s *Store) UserPostList(userId, limit, offset int) func(func(*UserPostListResult, error) bool) {
	return func(yield func(*UserPostListResult, error) bool) {
		// Proof that a SELECT statement is turing complete.
		query := fmt.Sprintf(`
			SELECT
				user_post.id,
				score,
				comments,
				IIF(user_post_upvote.id, 1, 0),
				IIF(user_post_downvote.id, 1, 0),
				username,
				title,
				url,
				body,
				user_post.created_at
			FROM user_post
			JOIN user ON user.id = user_post.user_id
			LEFT JOIN user_post_upvote ON
				user_post_upvote.post_id = user_post.id
				AND user_post_upvote.user_id = %[1]d
			LEFT JOIN user_post_downvote ON
				user_post_downvote.post_id = user_post.id
				AND user_post_downvote.user_id = %[1]d
			WHERE user_post.created_at > datetime('now', '-1 day')
			ORDER BY score DESC, user_post.created_at DESC
			LIMIT %d OFFSET %d`, userId, limit, offset)
		rows, err := s.db.Query(query)
		defer rows.Close()
		if err != nil {
			yield(nil, err)
			return
		}
		for rows.Next() {
			var post UserPostListResult
			err := rows.Scan(
				&post.Id,
				&post.Score,
				&post.Comments,
				&post.UserUpvoted,
				&post.UserDownvoted,
				&post.Author,
				&post.Title,
				&post.URL,
				&post.Body,
				&post.Date)
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
