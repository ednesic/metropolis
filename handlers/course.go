package handlers

import (
	"encoding/json"
	"github.com/ednesic/coursemanagement/servivemanager"
	"github.com/ednesic/coursemanagement/storage"
	"github.com/ednesic/coursemanagement/types"
	"github.com/gorilla/mux"
	"net/http"
)

func GetCourse(w http.ResponseWriter, r *http.Request) error {
	name := mux.Vars(r)["name"]
	course, err := servivemanager.CourseService.FindOne(name)

	if err == storage.ErrNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return err
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	j, err := json.Marshal(course)
	_, err = w.Write(j)
	return err
}

func GetCourses(w http.ResponseWriter, _ *http.Request) error {
	courses, err := servivemanager.CourseService.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	j, err := json.Marshal(courses)
	_, err = w.Write(j)
	return err
}

func SetCourse(w http.ResponseWriter, r *http.Request) error {
	var course types.Course

	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	err = servivemanager.CourseService.Create(course)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(course)
	_, _ = w.Write(j)
	return err
}

func PutCourse(w http.ResponseWriter, r *http.Request) error {
	var course types.Course
	err := json.NewDecoder(r.Body).Decode(&course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	err = servivemanager.CourseService.Update(course)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusCreated)
	j, err := json.Marshal(course)
	_, err = w.Write(j)
	return err
}

func DelCourse(w http.ResponseWriter, r *http.Request) error {
	name := mux.Vars(r)["name"]
	err := servivemanager.CourseService.Delete(types.Course{Name:name})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	j, err := json.Marshal("")
	_, err = w.Write(j)
	return err
}