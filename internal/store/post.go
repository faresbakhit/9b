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

func (s *Store) PostNew(post *UserPostNew) (id int, err error) {
	post.Body = strings.ReplaceAll(post.Body, "\r\n", "\n")
	query := `
		INSERT INTO post
		(user_id, title, url, body)
		VALUES (?, ?, ?, ?)
		RETURNING id`
	row := s.db.QueryRow(query, post.UserId, post.Title, post.URL, post.Body)
	err = row.Scan(&id)
	if err != nil {
		log.Print(err)
	}
	return
}

type PostGet struct {
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

func (s *Store) PostGet(userId, postId int) *PostGet {
	var post PostGet
	query := fmt.Sprintf(`
		SELECT
			post.id,
			score,
			IIF(post_upvote.id, 1, 0),
			IIF(post_downvote.id, 1, 0),
			username,
			title,
			url,
			body,
			post.created_at
		FROM post
		JOIN user ON user.id = post.user_id
		LEFT JOIN post_upvote ON
			post_upvote.post_id = post.id
			AND post_upvote.user_id = %[1]d
		LEFT JOIN post_downvote ON
			post_downvote.post_id = post.id
			AND post_downvote.user_id = %[1]d
		WHERE post.id = %d`, userId, postId)
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
		log.Print(err)
		return nil
	}
	return &post
}

type PostFrontpage struct {
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

func (s *Store) PostListFrontpage(userId, limit, offset int) func(func(*PostFrontpage, error) bool) {
	return func(yield func(*PostFrontpage, error) bool) {
		// Proof that a SELECT statement is turing complete.
		query := fmt.Sprintf(`
			SELECT
				post.id,
				score,
				comments,
				IIF(post_upvote.id, 1, 0),
				IIF(post_downvote.id, 1, 0),
				username,
				title,
				url,
				body,
				post.created_at
			FROM post
			JOIN user ON user.id = post.user_id
			LEFT JOIN post_upvote ON
				post_upvote.post_id = post.id
				AND post_upvote.user_id = %[1]d
			LEFT JOIN post_downvote ON
				post_downvote.post_id = post.id
				AND post_downvote.user_id = %[1]d
			WHERE post.created_at > datetime('now', '-1 day')
			ORDER BY score DESC, post.created_at DESC
			LIMIT %d OFFSET %d`, userId, limit, offset)
		rows, err := s.db.Query(query)
		defer rows.Close()
		if err != nil {
			log.Print(err)
			yield(nil, err)
			return
		}
		for rows.Next() {
			var post PostFrontpage
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
				log.Print(err)
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

func (s *Store) PostGetScore(postId int) (int, error) {
	var score int
	row := s.db.QueryRow("SELECT score FROM post WHERE id = ?", postId)
	err := row.Scan(&score)
	if err != nil {
		log.Print(err)
	}
	return score, err
}

func (s *Store) PostGetUpvotes(postId int) (int, error) {
	var score int
	row := s.db.QueryRow("SELECT COUNT() FROM post_upvote WHERE post_id = ?", postId)
	err := row.Scan(&score)
	if err != nil {
		log.Print(err)
	}
	return score, err
}

func (s *Store) PostGetDownvotes(postId int) (int, error) {
	var score int
	row := s.db.QueryRow("SELECT COUNT() FROM post_downvote WHERE post_id = ?", postId)
	err := row.Scan(&score)
	if err != nil {
		log.Print(err)
	}
	return score, err
}

func (s *Store) PostCreateUpvote(userId, postId int) (upvotes int, err error) {
	_, err = s.db.Exec("INSERT INTO post_upvote (user_id, post_id) VALUES (?, ?)", userId, postId)
	if err != nil {
		log.Print(err)
	}
	upvotes, err = s.PostGetUpvotes(postId)
	if err != nil {
		log.Print(err)
	}
	return
}

func (s *Store) PostCreateDownvote(userId, postId int) (downvotes int, err error) {
	_, err = s.db.Exec("INSERT INTO post_downvote (user_id, post_id) VALUES (?, ?)", userId, postId)
	if err != nil {
		log.Print(err)
	}
	downvotes, err = s.PostGetDownvotes(postId)
	if err != nil {
		log.Print(err)
	}
	return
}
