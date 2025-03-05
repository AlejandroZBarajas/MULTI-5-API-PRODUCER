package infrastructurec

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Connection struct {
	DB *sql.DB
}

var db *sql.DB
var err error

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error al cargar archivo .env: '%v'", err)
	}
}

func ConnectDB() {
	LoadEnv()

	User := os.Getenv("USER")
	Password := os.Getenv("PASSWORD")
	Host := os.Getenv("HOST")
	Port := os.Getenv("PORT")
	Name := os.Getenv("NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", User, Password, Host, Port, Name)

	db, err = sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("no se pudo conectar a mysql: '%v'", err)
	}

	fmt.Println("conectado a mysql")
}

func GetDB() *sql.DB {
	return db
}

func (c *Connection) RunQuery(query string, values ...interface{}) (sql.Result, error) {
	stmt, err := c.DB.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error: '%v'", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		return nil, fmt.Errorf("error al introducir datos: '%v'", err)
	}
	return result, nil
}

func (c *Connection) GetData(query string, values ...interface{}) (*sql.Rows, error) {
	rows, err := c.DB.Query(query, values...)
	if err != nil {
		fmt.Errorf("error: '%s'", err)
	}
	return rows, nil
}
