package repository

import (
	"context"

	"github.com/griggsjared/getsit/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXUrlEntryStore struct {
	db *pgxpool.Pool
}

func NewPGXUrlEntryStore(db *pgxpool.Pool) *PGXUrlEntryStore {
	return &PGXUrlEntryStore{
		db: db,
	}
}

type urlEntry struct {
	Token      string
	Url        string
	VisitCount int
}

func (s *PGXUrlEntryStore) SaveUrl(ctx context.Context, url entity.Url) (*entity.UrlEntry, error) {

	//if the url already exists return the entry
	entry, err := s.GetFromUrl(ctx, url)
	if err == nil {
		return entry, nil
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

	//insert the new url entry
	query := `
		INSERT INTO url_entries (url, token)
		VALUES ($1, $2)
	`
	_, err = s.db.Query(ctx, query, url, token.String())
	if err != nil {
		return nil, err
	}

	return &entity.UrlEntry{
		Token:      token,
		Url:        entity.Url(url),
		VisitCount: 0,
	}, nil
}

func (s *PGXUrlEntryStore) SaveVisit(ctx context.Context, token entity.UrlToken) error {

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

func (s *PGXUrlEntryStore) GetFromUrl(ctx context.Context, url entity.Url) (*entity.UrlEntry, error) {

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

func (s *PGXUrlEntryStore) GetFromToken(ctx context.Context, token entity.UrlToken) (*entity.UrlEntry, error) {

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
