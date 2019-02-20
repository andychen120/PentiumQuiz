package function

import (
	"database/sql"
	"encoding/json"
	_ "fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/openfaas-incubator/go-function-sdk"
)

var db *sql.DB

type PetOrder struct {
	Id       int    `json:"id"`
	Petid    int    `json:"petid"`
	Quantity int    `json:"quantity"`
	Shipdate string `json:"shipDate"`
	Status   string `json:"status"`
	Complete bool   `json:"complete"`
}

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

func addOrder(req handler.Request) (handler.Response, error) {
	var err error
	message := ""
	order := PetOrder{}
	status := http.StatusOK

	err = json.Unmarshal(req.Body, &order)
	if err != nil {
		// Check if input is folloed format.
		return handler.Response{
			Body:       []byte("Invalid Order"),
			StatusCode: http.StatusBadRequest,
		}, nil
	} else {
		_, err := db.Exec("INSERT INTO pet_order (id, petid, quantity, shipdate, status, complete) VALUES (?, ?, ?, ?, ?, ?)",
			order.Id, order.Petid, order.Quantity, order.Shipdate, order.Status, strconv.FormatBool(order.Complete))
		if err != nil {
			// Could not insert to DB.
			return handler.Response{
				Body:       []byte("Invalid Order"),
				StatusCode: http.StatusBadRequest,
			}, nil
		} else {
			message = string(req.Body)
		}
	}

	return handler.Response{
		Body:       []byte(message),
		StatusCode: status,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}, nil
}

func getOrder(req handler.Request) (handler.Response, error) {
	var err error
	message := ""

	// Check if id is integer.
	if _, err = strconv.Atoi(string(req.Body)); err != nil {
		return handler.Response{
			Body:       []byte("Invalid ID supplied"),
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	rows := db.QueryRow("SELECT * FROM pet_order WHERE id = ?", req.Body)

	result := PetOrder{}
	var complete string
	err = rows.Scan(&result.Id, &result.Petid, &result.Quantity, &result.Shipdate, &result.Status, &complete)
	if complete == "true" {
		result.Complete = true
	} else {
		result.Complete = false
	}

	// Check if id is existed.
	if err != nil {
		return handler.Response{
			Body:       []byte("Order not found"),
			StatusCode: http.StatusNotFound,
		}, nil
	} else {
		// Get results, marchal it to json string.
		jresult, err := json.Marshal(result)
		if err != nil {
			message = err.Error()
		} else {
			message = string(jresult)
		}
	}

	return handler.Response{
		Body:       []byte(message),
		StatusCode: http.StatusOK,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
	}, err
}

func deleteOrder(req handler.Request) (handler.Response, error) {
	var err error

	// Check if id is integer.
	if _, err = strconv.Atoi(string(req.Body)); err != nil {
		return handler.Response{
			Body:       []byte("Invalid ID supplied"),
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	// Check if id is existed
	var exists bool
	if err := db.QueryRow("SELECT exists (SELECT id FROM pet_order WHERE id = ?)", req.Body).Scan(&exists); err != nil || !exists {
		return handler.Response{
			Body:       []byte("Order not found"),
			StatusCode: http.StatusNotFound,
		}, nil
	}

	db.Exec("DELETE FROM pet_order WHERE id = ?", req.Body)

	return handler.Response{
		Body:       []byte(""),
		StatusCode: http.StatusOK,
	}, err
}

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	var err error
	var response handler.Response

	err = connectDB()
	if err != nil {
		return handler.Response{
			Body:       []byte(err.Error()),
			StatusCode: http.StatusOK,
		}, err
	}
	defer db.Close()

	switch req.Method {
	case "GET":
		response, err = getOrder(req)
	case "POST":
		response, err = addOrder(req)
	case "DELETE":
		response, err = deleteOrder(req)
	default:
	}
	//message := fmt.Sprintf("Method was: %s, body was: %s", string(req.Method), string(req.Body))

	return response, err
}
