package repository

import (
	"context"
	"fmt"

	"github.com/griggsjared/getsit/internal/url/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXUrlEntryRepository struct {
	db *pgxpool.Pool
}

func NewPGXUrlEntryRepository(db *pgxpool.Pool) *PGXUrlEntryRepository {
	return &PGXUrlEntryRepository{
		db: db,
	}
}

type urlEntry struct {
	Token      string
	Url        string
	VisitCount int
}

func (s *PGXUrlEntryRepository) SaveUrl(ctx context.Context, url entity.Url) (*entity.UrlEntry, error) {

	//if the url already exists escape with an error
	_, err := s.GetFromUrl(ctx, url)
	if err == nil {
		return nil, fmt.Errorf("entry already exists")
	}

	//find a unique token
	var token entity.UrlToken
	for {
		token = entity.NewUrlToken()
		_, err := s.GetFromToken(ctx, token)
		if err != nil {
			break
		}
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO url_entries (url, token)
		VALUES ($1, $2)
	`
	_, err = s.db.Exec(ctx, query, url, token.String())
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &entity.UrlEntry{
		Token:      token,
		Url:        entity.Url(url),
		VisitCount: 0,
	}, nil
}

func (s *PGXUrlEntryRepository) SaveVisit(ctx context.Context, token entity.UrlToken) error {

	query := `
		UPDATE url_entries
		SET visit_count = visit_count + 1
		WHERE token = $1
	`

	_, err := s.db.Exec(ctx, query, token)
	if err != nil {
		return err
	}

	return nil
}

func (s *PGXUrlEntryRepository) GetFromUrl(ctx context.Context, url entity.Url) (*entity.UrlEntry, error) {

	query := `
		SELECT token, url, visit_count
		FROM url_entries
		WHERE url = $1
	`

	row := s.db.QueryRow(ctx, query, url)

	var urlEntry urlEntry
	err := row.Scan(&urlEntry.Token, &urlEntry.Url, &urlEntry.VisitCount)
	if err != nil {
		return nil, err
	}

	return &entity.UrlEntry{
		Token:      entity.UrlToken(urlEntry.Token),
		Url:        entity.Url(url),
		VisitCount: urlEntry.VisitCount,
	}, nil
}

func (s *PGXUrlEntryRepository) GetFromToken(ctx context.Context, token entity.UrlToken) (*entity.UrlEntry, error) {

	query := `
		SELECT token, url, visit_count
		FROM url_entries
		WHERE token = $1
	`
	row := s.db.QueryRow(ctx, query, token)

	var urlEntry urlEntry
	err := row.Scan(&urlEntry.Token, &urlEntry.Url, &urlEntry.VisitCount)
	if err != nil {
		return nil, err
	}

	return &entity.UrlEntry{
		Token:      entity.UrlToken(token),
		Url:        entity.Url(urlEntry.Url),
		VisitCount: urlEntry.VisitCount,
	}, nil
}
