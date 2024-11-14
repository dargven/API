package sqlite

import (
	storage "API/internal/Storage"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// Пока не рабочий калл который надо переделывать

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage" // Имя текущей функции для логов и ошибок
	//Подключапемся к бд
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	// В последствие есть необходимость в использовании миграций для создания и контроля версий за БД.
	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
		    id INTEGER PRIMARY KEY,
		    alias TEXT NOT NULL UNIQUE,
		    url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);

`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

// Методы SaveUrl и getUrl

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	//Подготавливаем запрос
	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) valus(?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	// выполняем запрос
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}
		return 0, fmt.Errorf("%s:execute statement: %w", op, err)
	}
	// Получаем id только что созданной записи

	//Примечание: в postgresql что то подобное поддерживается и будет работать
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failet to get last insert id : %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s:prepare statement: %w", op, err)
	}
	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}
	return resURL, nil
}

//Метод DeleteUrl будет написан позже
