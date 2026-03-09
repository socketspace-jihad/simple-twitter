package postgresql

import (
	"simple_twitter/internal/models"

	"github.com/rs/zerolog/log"
)

func (p *PostgreSQL) SavePost(post *models.Post) error {
	tx, err := p.DB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec("INSERT INTO posts (content,user_id,image_url) VALUES ($1,$2,$3)", post.Content, post.User.ID, post.ImageURL); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (p *PostgreSQL) UpdatePost(post *models.Post) error {
	tx, err := p.DB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	log.Debug().Str("id", post.ID.String()).Str("content", post.Content).Msg("debug body")
	if _, err := tx.Exec("UPDATE posts SET content = $1 WHERE id = $2", post.Content, post.ID); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (p *PostgreSQL) GetPost(post *models.Post) error {
	return nil
}

func (p *PostgreSQL) DeletePost(post *models.Post) error {
	tx, err := p.DB.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec("DELETE FROM posts WHERE id=$1", post.ID); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p *PostgreSQL) ListPosts() ([]models.Post, error) {
	rows, err := p.DB.Query("SELECT posts.id, content, COALESCE(posts.image_url,''), posts.created_at, u.display_name, u.username, u.id FROM posts LEFT JOIN users AS u ON u.id = posts.user_id ORDER BY created_at DESC LIMIT 500")
	if err != nil {
		return nil, err
	}
	posts := []models.Post{}
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Content, &post.ImageURL, &post.CreatedAt, &post.DisplayName, &post.Username, &post.User.ID); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
