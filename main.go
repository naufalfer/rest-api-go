package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest-api-go/auth"
	"rest-api-go/database"
	"rest-api-go/models"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

func main() {
	database.Connect()
	router := mux.NewRouter()

	router.HandleFunc("/login", Login).Methods("POST")
	router.HandleFunc("/register", Register).Methods("POST")

	router.HandleFunc("/karyawan", getAllKaryawan).Methods("GET")
	router.HandleFunc("/karyawan", createKaryawan).Methods("POST")
	router.HandleFunc("/karyawan/{id}", getKaryawan).Methods("GET")
	router.HandleFunc("/karyawan/{id}", updateKaryawan).Methods("PUT")
	router.HandleFunc("/karyawan/{id}", deleteKaryawan).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	handler := c.Handler(router)

	http.ListenAndServe(":8080", handler)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	user, token, err := auth.Authenticate(user.Username, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	createdUser, err := auth.Register(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(createdUser)
}

func createKaryawan(w http.ResponseWriter, r *http.Request) {
	var karyawan models.Karyawan
	if err := json.NewDecoder(r.Body).Decode(&karyawan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	karyawan.CreatedAt = time.Now()

	if err := database.DB.Create(&karyawan).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(karyawan)
}

func getKaryawan(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var karyawan models.Karyawan

	if err := database.DB.First(&karyawan, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Karyawan not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(karyawan)
}

func getAllKaryawan(w http.ResponseWriter, r *http.Request) {
	var karyawan []models.Karyawan

	result := database.DB.Find(&karyawan)
	if result.Error != nil {
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(karyawan)
	if err != nil {
		http.Error(w, "Gagal membuat response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func updateKaryawan(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var karyawan models.Karyawan
	if err := database.DB.First(&karyawan, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Karyawan not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var updatedKaryawan models.Karyawan
	if err := json.NewDecoder(r.Body).Decode(&updatedKaryawan); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedKaryawan.ID = karyawan.ID
	updatedKaryawan.CreatedAt = karyawan.CreatedAt
	updatedKaryawan.UpdatedAt = time.Now()

	if err := database.DB.Save(&updatedKaryawan).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedKaryawan)
}

func deleteKaryawan(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var karyawan models.Karyawan
	if err := database.DB.First(&karyawan, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Karyawan not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := database.DB.Delete(&karyawan).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Karyawan with ID %s has been deleted", id)
}
