package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/vaib1343/students-api/internal/http/util/response"
	"github.com/vaib1343/students-api/internal/storage"
	"github.com/vaib1343/students-api/internal/types"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("body is empty")))
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validationErrors))
			return
		}

		id, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": id})

	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		intId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, student)

	}
}

func GetStudents(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := storage.GetAllStudent()
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, err)
			return
		}

		response.WriteJSON(w, http.StatusOK, students)
	}
}
