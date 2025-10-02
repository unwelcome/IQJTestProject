package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/unwelcome/iqjtest/internal/entities"
)

type CatRepository struct {
	db *sql.DB
}

func NewCatRepository(db *sql.DB) *CatRepository {
	return &CatRepository{db: db}
}

func (r *CatRepository) CreateCat(ctx context.Context, userID int, cat *entities.Cat) error {
	query := `INSERT INTO cats(name, age, description, created_by) VALUES ($1, $2, $3, $4) RETURNING id;`

	err := r.db.QueryRowContext(ctx, query, cat.Name, cat.Age, cat.Description, userID).Scan(&cat.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CatRepository) GetCatByID(ctx context.Context, catID int) (*entities.CatWithPhotos, error) {
	// Запрос на получение кота с left join его фото с сортировкой по is_primary и cp.id
	query := `
		SELECT 
			c.name,
			c.age,
			c.description,	
			c.created_at, 
			c.created_by,
			cp.id, 
			cp.url, 
			cp.is_primary
		FROM cats c
		LEFT JOIN cat_photos cp ON c.id = cp.cat_id
		WHERE c.id = $1
		ORDER BY cp.is_primary DESC, cp.id ASC;
	`

	// Выполняем запрос
	rows, err := r.db.QueryContext(ctx, query, catID)
	if err != nil {
		return nil, err
	}

	// cat.ID = catID - 1 флаг, что структура cat пуста
	cat := &entities.CatWithPhotos{ID: catID - 1}
	var catPhotos []*entities.CatPhotoUrl

	// Переменные для скана
	var (
		catAge, catCreatedBy                  int
		catName, catDescription, catCreatedAt string
		photoID                               sql.NullInt64
		photoUrl                              sql.NullString
		photoIsPrimary                        sql.NullBool
	)

	// Мэппинг ответа в структуру
	for rows.Next() {
		err := rows.Scan(
			&catName, &catAge, &catDescription, &catCreatedAt, &catCreatedBy,
			&photoID, &photoUrl, &photoIsPrimary,
		)
		if err != nil {
			return nil, err
		}

		// Если cat пустой заносим значения в поля кота
		if cat.ID != catID {
			cat.ID = catID
			cat.Name = catName
			cat.Age = catAge
			cat.Description = catDescription
			cat.CreatedAt = catCreatedAt
			cat.CreatedBy = catCreatedBy
		}

		// Если все поля не null то создаем объект фото
		if photoID.Valid && photoUrl.Valid && photoIsPrimary.Valid {
			catPhotos = append(catPhotos, &entities.CatPhotoUrl{
				ID:        int(photoID.Int64),
				Url:       photoUrl.String,
				IsPrimary: photoIsPrimary.Bool,
			})
		}
	}

	// Цикл не выполнился ни разу, значит кота нет
	if cat.ID == (catID - 1) {
		return nil, fmt.Errorf("cat not found")
	}

	cat.Photos = catPhotos
	return cat, nil
}

func (r *CatRepository) GetAllCats(ctx context.Context) ([]*entities.CatWithPrimePhoto, error) {
	// Запрос на получение всех котов с left join фото котов, сортируя по catID, затем по is_primary и в конце по photoID
	// Т.о. Получаем кота с первым is_primary фото либо кота с первым фото либо кота без фото
	query := `
		SELECT DISTINCT ON (c.id)
			c.id,
			c.name,
			c.age,
			cp.id AS photo_id, 
			cp.url
		FROM cats c
		LEFT JOIN cat_photos cp ON c.id = cp.cat_id
		ORDER BY c.id, cp.is_primary DESC NULLS LAST, cp.id ASC;
	`

	// Выполняем запрос
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		cats    []*entities.CatWithPrimePhoto
		photoID sql.NullInt64
		url     sql.NullString
	)

	// Мэппинг ответа в структуру
	for rows.Next() {
		cat := &entities.CatWithPrimePhoto{}

		err = rows.Scan(&cat.ID, &cat.Name, &cat.Age, &photoID, &url)
		if err != nil {
			return nil, err
		}

		// Если photoID не null
		if photoID.Valid {
			id := int(photoID.Int64)
			cat.PhotoID = &id
		}
		// Если url не null
		if url.Valid {
			urlStr := url.String
			cat.Url = &urlStr
		}

		cats = append(cats, cat)
	}

	return cats, nil
}

func (r *CatRepository) UpdateCatName(ctx context.Context, catID int, newName string) error {
	query := `UPDATE cats SET name = $1 WHERE id = $2;`

	_, err := r.db.ExecContext(ctx, query, newName, catID)
	if err != nil {
		return err
	}

	return nil
}

func (r *CatRepository) UpdateCatAge(ctx context.Context, catID int, newAge int) error {
	query := `UPDATE cats SET age = $1 WHERE id = $2;`

	_, err := r.db.ExecContext(ctx, query, newAge, catID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CatRepository) UpdateCatDescription(ctx context.Context, catID int, newDescription string) error {
	query := `UPDATE cats SET description = $1 WHERE id = $2;`

	_, err := r.db.ExecContext(ctx, query, newDescription, catID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CatRepository) UpdateCat(ctx context.Context, catID int, catUpdateRequest *entities.CatUpdateRequest) error {
	query := `UPDATE cats SET (name, age, description) = ($1, $2, $3) WHERE id = $4;`

	_, err := r.db.ExecContext(ctx, query, catUpdateRequest.Name, catUpdateRequest.Age, catUpdateRequest.Description, catID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CatRepository) DeleteCat(ctx context.Context, catID int) error {
	query := `DELETE FROM cats WHERE id = $1;`

	_, err := r.db.ExecContext(ctx, query, catID)
	if err != nil {
		return err
	}
	return nil
}
