package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/Gin-postgres-docker/config"
	"github.com/Deo-Mugabe/GOLANG/tree/main/RestAPi/Gin-postgres-docker/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func SetupTestDB() {
	var err error
	config.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	config.DB.AutoMigrate(&models.User{})
}

func TestGetUsers(t *testing.T) {
	SetupTestDB()
	r := SetupRouter()
	r.GET("/users", GetUsers)

	// Create a test user
	config.DB.Create(&models.User{Name: "John Doe", Email: "john@example.com", Password: "password"})

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var users []models.User
	err := json.Unmarshal(w.Body.Bytes(), &users)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, "John Doe", users[0].Name)
	assert.Equal(t, "john@example.com", users[0].Email)
}

func TestCreateUser(t *testing.T) {
	SetupTestDB()
	r := SetupRouter()
	r.POST("/users", CreateUser)

	user := models.User{Name: "Jane Doe", Email: "jane@example.com", Password: "password"}
	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var createdUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	assert.Nil(t, err)
	assert.Equal(t, "Jane Doe", createdUser.Name)
	assert.Equal(t, "jane@example.com", createdUser.Email)
	assert.NotZero(t, createdUser.Id)
}

func TestDeleteUser(t *testing.T) {
	SetupTestDB()
	r := SetupRouter()
	r.DELETE("/users/:id", DeleteUser)

	// Create a test user
	user := models.User{Name: "John Doe", Email: "john@example.com", Password: "password"}
	config.DB.Create(&user)

	req, _ := http.NewRequest("DELETE", "/users/"+strconv.Itoa(user.Id), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var deletedUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &deletedUser)
	assert.Nil(t, err)
	assert.Equal(t, user.Id, deletedUser.Id)

	// Ensure user is deleted
	var users []models.User
	config.DB.Find(&users)
	assert.Equal(t, 0, len(users))
}

func TestUpdateUser(t *testing.T) {
	SetupTestDB()
	r := SetupRouter()
	r.PUT("/users/:id", UpdateUser)

	// Create a test user
	user := models.User{Name: "John Doe", Email: "john@example.com", Password: "password"}
	config.DB.Create(&user)

	updatedUser := models.User{Name: "John Smith", Email: "johnsmith@example.com", Password: "newpassword"}
	jsonValue, _ := json.Marshal(updatedUser)
	req, _ := http.NewRequest("PUT", "/users/"+strconv.Itoa(user.Id), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var newUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &newUser)
	assert.Nil(t, err)
	assert.Equal(t, "John Smith", newUser.Name)
	assert.Equal(t, "johnsmith@example.com", newUser.Email)

	// Ensure user is updated in DB
	config.DB.First(&user, user.Id)
	assert.Equal(t, "John Smith", user.Name)
	assert.Equal(t, "johnsmith@example.com", user.Email)
}

func TestGetUserById(t *testing.T) {
	SetupTestDB()
	r := SetupRouter()
	r.GET("/users/:id", GetUserById)

	// Create a test user
	user := models.User{Name: "John Doe", Email: "john@example.com", Password: "password"}
	config.DB.Create(&user)

	req, _ := http.NewRequest("GET", "/users/"+strconv.Itoa(user.Id), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var fetchedUser models.User
	err := json.Unmarshal(w.Body.Bytes(), &fetchedUser)
	assert.Nil(t, err)
	assert.Equal(t, user.Id, fetchedUser.Id)
	assert.Equal(t, user.Name, fetchedUser.Name)
	assert.Equal(t, user.Email, fetchedUser.Email)
}
