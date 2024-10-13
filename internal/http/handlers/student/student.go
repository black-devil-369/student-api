package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/black-devil-369/student-api/internal/response"
	"github.com/black-devil-369/student-api/internal/storage"
	"github.com/black-devil-369/student-api/internal/types"
	"github.com/go-playground/validator/v10"
)

func NewUser(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GenralError(fmt.Errorf("empty body")))
			return
		}
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GenralError(err))
			return
		}
		// request validation
		if err := validator.New().Struct(student); err != nil {
			validatorErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidatationError(validatorErrs))
			return
		}
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		slog.Info("user created successfully", slog.String("userId", fmt.Sprint(lastId)))
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting a student", slog.String("id", id))

		IntId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GenralError(err))
			return
		}

		student, err := storage.GetStudentById(IntId)
		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GenralError(err))
			return
		}
		response.WriteJson(w, http.StatusOK, student)
		// id := r.URL.Query().Get("id")
		// student, err := storage.GetStudent(id)
		// if err != nil {
		// 	response.WriteJson(w, http.StatusNotFound, err)
		// 	return
		// }
		// response.WriteJson(w, http.StatusOK, student)
	}
}
