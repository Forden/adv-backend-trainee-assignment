package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"adv-backend-trainee-assignment/src/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreSQLManager struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewPostgreSQLManager(dbUrl string) (*PostgreSQLManager, error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, dbUrl)
	if err != nil {
		return nil, fmt.Errorf("can't open postgresql db: %s", err)
	} else {
		return &PostgreSQLManager{pool: db, ctx: ctx}, nil
	}
}

func (postgre PostgreSQLManager) Close() error {
	postgre.pool.Close()
	return nil
}

func (postgre PostgreSQLManager) NewAd(adData models.CreatingAd) (string, error) {
	adID := uuid.New().String()
	marshalledPhotoLinks, err := json.Marshal(adData.PhotoLinks)
	if err != nil {
		return "", err
	} else {
		_, err := postgre.pool.Exec(postgre.ctx, "INSERT INTO ads (ad_id, title, description, price, photo_links, created_at) VALUES ($1, $2, $3, $4, $5, $6)", adID, adData.Title, adData.Description, adData.Price, marshalledPhotoLinks, time.Now().UTC().Unix())
		if err != nil {
			return "", err
		}
		return adID, nil
	}
}

func (postgre PostgreSQLManager) SelectAd(adID string) (*models.DbAd, error) {
	rows, err := postgre.pool.Query(postgre.ctx, "SELECT * FROM ads WHERE ad_id = $1", adID)
	if err != nil {
		return nil, err
	} else {
		defer rows.Close()
		for rows.Next() {
			var res models.DbAd
			var tmp string
			var createdAt int64
			err := rows.Scan(&res.AdID, &res.Title, &res.Description, &res.Price, &tmp, &createdAt)
			if err != nil {
				return nil, err
			} else {
				res.CreatedAt = time.Unix(createdAt, 0)
				err := json.Unmarshal([]byte(tmp), &res.PhotoLinks)
				if err != nil {
					return nil, err
				} else {
					return &res, nil
				}
			}
		}
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return nil, nil
}

func (postgre PostgreSQLManager) GetAllAds(sortBy string, sortOrder string, page int, perPage int) ([]*models.DbAd, error) {
	rows, err := postgre.pool.Query(postgre.ctx, fmt.Sprintf("SELECT * FROM ads ORDER BY %s %s LIMIT %d OFFSET %d", sortBy, sortOrder, perPage, (page-1)*perPage))
	if err != nil {
		return nil, err
	} else {
		defer rows.Close()
		var result []*models.DbAd
		for rows.Next() {
			var res models.DbAd
			var tmp string
			var createdAt int64
			err := rows.Scan(&res.AdID, &res.Title, &res.Description, &res.Price, &tmp, &createdAt)
			if err != nil {
				return nil, err
			} else {
				res.CreatedAt = time.Unix(createdAt, 0)
				err := json.Unmarshal([]byte(tmp), &res.PhotoLinks)
				if err != nil {
					return nil, err
				} else {
					result = append(result, &res)
				}
			}
		}
		return result, err
	}
}
