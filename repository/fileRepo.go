package repository

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gocleancode/db"
	"time"
)

type FileRepo interface {
	SaveFile(file File) (int64, error)
	GetFileById(id int64) (File, error)
	TxDeleteFileById(id int64, tx *sql.Tx) error
}

type fileRepo struct {
	Db db.DB
}

type File struct {
	Id          *int64
	FileName    *string
	FilePath    *string
	ContentType *string
	CreatedDt   *time.Time
}

func NewFileRepo(db db.DB) FileRepo {
	return fileRepo{Db: db}
}

func (repo fileRepo) SaveFile(file File) (int64, error) {
	var generatedId int64
	if file.CreatedDt == nil {
		now := time.Now()
		file.CreatedDt = &now
	}
	stmt, err := repo.Db.Prepare("INSERT INTO files(file_name, file_path, content_type, created_dt) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Error(err)
		return generatedId, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(file.FileName, file.FilePath, file.ContentType, file.CreatedDt)
	if err != nil {
		log.Error(err)
		return generatedId, err
	}
	generatedId, err = res.LastInsertId()
	if err != nil {
		log.Error(err)
		return generatedId, err
	}
	return generatedId, nil
}

func (repo fileRepo) GetFileById(id int64) (File, error) {
	file := File{}
	row := repo.Db.QueryRow("SELECT id, file_name, file_path, content_type, created_dt from files where id = ?", id)
	err := row.Scan(&file.Id, &file.FileName, &file.FilePath, &file.ContentType, &file.CreatedDt)
	if err != nil {
		return file, err
	}
	return file, nil
}

func (repo fileRepo) TxDeleteFileById(id int64, tx *sql.Tx) error {
	stmt, err := tx.Prepare("DELETE from files where id = ?")
	defer stmt.Close()
	if err != nil {
		log.Error(err)
		return err
	}
	res, err := stmt.Exec(id)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debug(fmt.Sprintf("TxDeleteFileById response = %v", res))
	return nil
}
