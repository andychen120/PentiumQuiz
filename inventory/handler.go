package function

import (
	"database/sql"
	"encoding/json"
	_ "fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/openfaas-incubator/go-function-sdk"
)

type Data struct {
	status   string
	quantity int
	count    int
}

var db *sql.DB

func connectDB() error {
	var err error
	db, err = sql.Open("mysql", "myuser:mypass@tcp(34.80.17.17:3306)/pentium")
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return err
}

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	message := ""
	err := connectDB()
	if err != nil {
		return handler.Response{
			Body:       []byte(err.Error()),
			StatusCode: http.StatusOK,
		}, err
	}
	defer db.Close()

	datas := make([]*Data, 0)
	rows, _ := db.Query("SELECT status, quantity, COUNT(*) as count FROM pet_order GROUP BY status, quantity")
	for rows.Next() {
		data := new(Data)

		rows.Scan(&data.status, &data.quantity, &data.count)
		datas = append(datas, data)
	}

	var result map[string]int

	result = make(map[string]int)

	for _, data := range datas {
		result[data.status] = result[data.status] + (data.quantity * data.count)
	}

	jresult, err := json.Marshal(result)
	if err != nil {
		message = err.Error()
	} else {
		message = string(jresult)
	}

	return handler.Response{
		Body:       []byte(message),
		StatusCode: http.StatusOK,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}, err
}
