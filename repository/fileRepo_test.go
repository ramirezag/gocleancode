package repository_test

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	myDb "gocleancode/db"
	"gocleancode/repository"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"regexp"
	"testing"
	"time"
)

func TestSaveFile(t *testing.T) {
	// Given
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("an error was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()
	repo := repository.NewFileRepo(myDb.DB{mockDb, "mockdb"})

	fileName := "fname"
	filePath := "/some/file/path"
	contentType := "contentType"
	createdDt := time.Now()
	expectedId := int64(1)
	sqlRegexStr := regexp.QuoteMeta("INSERT INTO files(file_name, file_path, content_type, created_dt) VALUES(?, ?, ?, ?)")
	mock.
		ExpectPrepare(sqlRegexStr).
		ExpectExec().
		WithArgs(&fileName, &filePath, &contentType, &createdDt).
		WillReturnResult(sqlmock.NewResult(expectedId, 1))
	file := repository.File{FileName: &fileName, FilePath: &filePath, ContentType: &contentType, CreatedDt: &createdDt}
	// When
	actualGeneratedId, err := repo.SaveFile(file)
	// Then
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
	}
	assert.Equal(t, expectedId, actualGeneratedId)
}

func TestGetFileById(t *testing.T) {
	// Given
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("an error was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()
	repo := repository.NewFileRepo(myDb.DB{mockDb, "mockdb"})

	id := int64(1)
	fileName := "fname"
	filePath := "/some/file/path"
	contentType := "contentType"
	createdDt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "file_name", "file_path", "content_type", "created_dt"}).
		AddRow(&id, &fileName, &filePath, &contentType, &createdDt)

	expectedFile := repository.File{Id: &id, FileName: &fileName, FilePath: &filePath, ContentType: &contentType, CreatedDt: &createdDt}

	mock.
		ExpectQuery("SELECT id, file_name, file_path, content_type, created_dt from files where id = ?").
		WithArgs(id).
		WillReturnRows(rows)
	// When
	actualFile, err := repo.GetFileById(id)
	// Then
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.Equal(t, expectedFile, actualFile)
}

func TestTxDeleteFileById(t *testing.T) {
	// Given
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("an error was not expected when opening a stub database connection", err)
	}
	defer mockDb.Close()
	mockmyDb := myDb.DB{mockDb, "mockdb"}
	id := int64(1)

	mock.ExpectBegin()
	mock.
		ExpectPrepare("DELETE from files where id = ?").
		ExpectExec().
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	// When
	err = mockmyDb.Transact(func(tx *sql.Tx) error {
		repo := repository.NewFileRepo(mockmyDb)
		return repo.TxDeleteFileById(id, tx)
	})
	// Then
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
