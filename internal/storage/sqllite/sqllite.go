package sqllite

import (
	"database/sql"
	"fmt"

	"github.com/black-devil-369/student-api/internal/config"
	"github.com/black-devil-369/student-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type DataBase struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*DataBase, error) {
	db, err := sql.Open("sqlite3", cfg.Storage_Path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY,
	name TEXT,
	age INTEGER,
	email TEXT
	)`)
	if err != nil {
		return nil, err
	}
	return &DataBase{
		Db: db,
	}, nil
}

func (s *DataBase) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students(name,email,age) VALUES(?,?,?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// print the data in termila which are storge in Sqllite Database
	fmt.Println(" ")
	fmt.Printf("Name\tEmail\tAge\n")
	fmt.Println(" ")
	fmt.Printf("%s\t%s\t%d\n", name, email, age)
	fmt.Println("")
	// end here
	return lastId, nil
}

func (s *DataBase) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()
	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no Student found with id %s", fmt.Sprint(id))
		}
		return types.Student{}, fmt.Errorf("querry Error %w", err)
	}
	return student, nil
}
