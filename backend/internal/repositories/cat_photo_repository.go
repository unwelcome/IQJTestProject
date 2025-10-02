package repositories

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/unwelcome/iqjtest/internal/entities"
	"path/filepath"
)

type CatPhotoRepository struct {
	db          *sql.DB
	minioClient *minio.Client
	endpoint    string
	bucketName  string
}

func NewCatPhotoRepository(db *sql.DB, minioClient *minio.Client, endpoint, bucketName string) *CatPhotoRepository {
	return &CatPhotoRepository{
		db:          db,
		minioClient: minioClient,
		endpoint:    endpoint,
		bucketName:  bucketName,
	}
}

func (r *CatPhotoRepository) AddCatPhoto(ctx context.Context, catID int, req *entities.CatPhotoUploadRequest) (*entities.CatPhotoUploadResponse, error) {
	// Генерируем уникальное имя файла
	filename := generateFilename(catID, req.FileName)

	// Сохраняем файл в Minio
	_, err := r.minioClient.PutObject(
		ctx,
		r.bucketName,
		filename,
		req.File,
		req.FileSize,
		minio.PutObjectOptions{
			ContentType: req.MimeType,
		})
	if err != nil {
		return nil, err
	}

	// Создаем тело ответа
	res := &entities.CatPhotoUploadResponse{FileName: filename}

	// Создаем публичный url, формат: http://localhost:9000/bucket-name/filename
	res.Url = fmt.Sprintf("http://%s/%s/%s", r.endpoint, r.bucketName, filename)

	// Сохраняем фото в бд
	query := `INSERT INTO cat_photos (cat_id, url, filename, filesize, mime_type, is_primary) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`
	err = r.db.QueryRowContext(ctx, query, catID, res.Url, filename, req.FileSize, req.MimeType, req.IsPrimary).Scan(&res.ID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *CatPhotoRepository) DeleteCatPhoto(ctx context.Context, photoID int) error {
	// Удаляем файл из бд и получаем filename
	var filename string
	query := `DELETE FROM cat_photos WHERE id = $1 RETURNING filename;`
	err := r.db.QueryRowContext(ctx, query, photoID).Scan(&filename)
	if err != nil {
		return err
	}

	// Удаляем файл из minio
	err = r.minioClient.RemoveObject(ctx, r.bucketName, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (r *CatPhotoRepository) SetCatPhotoPrimary(ctx context.Context, catID, photoID int) error {
	// Создаем транзакцию
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx error: %w", err)
	}
	defer tx.Rollback()

	// Убираем у всех фото конкретного кота is_primary
	_, err = tx.ExecContext(ctx, `UPDATE cat_photos SET is_primary = false WHERE cat_id = $1;`, catID)
	if err != nil {
		return fmt.Errorf("remove all is_primary error: %w", err)
	}

	// Устанавливаем is_primary для выбранной фото
	_, err = tx.ExecContext(ctx, `UPDATE cat_photos SET is_primary = true WHERE id = $1;`, photoID)
	if err != nil {
		return fmt.Errorf("set is_primary by id error: %w", err)
	}

	// Коммитим транзакцию
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit tx error: %w", err)
	}
	return nil
}

func (r *CatPhotoRepository) GetAllCatPhotos(ctx context.Context, catID int) ([]*entities.CatPhotoUrl, error) {
	query := `SELECT id, url, is_primary FROM cat_photos WHERE cat_id = $1;`

	rows, err := r.db.QueryContext(ctx, query, catID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var catPhotos []*entities.CatPhotoUrl
	for rows.Next() {
		catPhoto := &entities.CatPhotoUrl{}
		err = rows.Scan(&catPhoto.ID, &catPhoto.Url, &catPhoto.IsPrimary)
		if err != nil {
			return nil, err
		}
		catPhotos = append(catPhotos, catPhoto)
	}

	return catPhotos, nil
}

func (r *CatPhotoRepository) GetCatPhotoByID(ctx context.Context, photoID int) (*entities.CatPhoto, error) {
	query := `SELECT id, url, filename, filesize, mime_type, is_primary, created_at FROM cat_photos WHERE id = $1;`

	catPhoto := &entities.CatPhoto{ID: photoID}
	err := r.db.QueryRowContext(ctx, query, photoID).Scan(&catPhoto.Url, &catPhoto.FileName, &catPhoto.FileSize, &catPhoto.MimeType, &catPhoto.IsPrimary, &catPhoto.CreatedAt)
	if err != nil {
		return nil, err
	}
	return catPhoto, nil
}

func generateFilename(catID int, fileName string) string {
	// Извлекаем расширение файла
	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".jpg"
	}

	// Генерируем уникальный ID для имени файла
	uniqueID := generateUniqueID()

	// Формируем название файла cat/{catID}/{uuid}{ext}
	return fmt.Sprintf("cat/%d/%s%s", catID, uniqueID, ext)
}

func generateUniqueID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}
