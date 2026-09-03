package pg

import (
	"VideoBot/internal/link"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type LinkRepo struct{}

func NewLinkRepo() *LinkRepo {
	return &LinkRepo{}
}

func (r *LinkRepo) Save(ctx context.Context, lnk link.Link) error {
	tx, err := TxFromContext(ctx)
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO links (id, link, link_type, chat_id)
		VALUES (@id, @link, @type, @chat_id)
	`

	args := pgx.NamedArgs{
		"id":      lnk.ID(),
		"link":    lnk.Link(),
		"type":    lnk.LinkType().Int(),
		"chat_id": lnk.ChatId(),
	}

	if _, err := tx.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("вставка в бд: %w", err)
	}

	return nil
}
