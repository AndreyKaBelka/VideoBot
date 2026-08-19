package link

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) Save(ctx context.Context, tx pgx.Tx, link Link) error {
	const query = `
		INSERT INTO links (id, link, link_type, chat_id)
		VALUES (@id, @link, @type, @chat_id)
	`

	args := pgx.NamedArgs{
		"id":      link.id,
		"link":    link.link,
		"type":    link.linkType,
		"chat_id": link.chatId,
	}
	_, err := tx.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("вставка в бд: %w", err)
	}

	return nil
}
