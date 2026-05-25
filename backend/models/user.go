package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"backend/utils"

	beego "github.com/beego/beego/v2/server/web"
)

// User represents an application user stored in CSV.
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

var (
	userCSVHeader           = []string{"id", "name", "email", "password", "created_at"}
	userCSVFilePathOverride string
)

// GetAllUsers returns all users from the configured CSV file.
func GetAllUsers() ([]User, error) {
	filePath, err := getUserCSVFilePath()
	if err != nil {
		return nil, err
	}

	if err := utils.EnsureFileExists(filePath, userCSVHeader); err != nil {
		return nil, err
	}

	rows, err := utils.ReadCSV(filePath)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)
	for index, row := range rows {
		if index == 0 || isEmptyCSVRow(row) {
			continue
		}

		user, err := parseUserRow(row, index+1)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// GetUserByEmail returns a user matching the provided email, or nil when not found.
func GetUserByEmail(email string) (*User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}

	for index := range users {
		if strings.EqualFold(strings.TrimSpace(users[index].Email), strings.TrimSpace(email)) {
			return &users[index], nil
		}
	}

	return nil, nil
}

// GetUserByID returns a user matching the provided ID, or nil when not found.
func GetUserByID(id int) (*User, error) {
	users, err := GetAllUsers()
	if err != nil {
		return nil, err
	}

	for index := range users {
		if users[index].ID == id {
			return &users[index], nil
		}
	}

	return nil, nil
}

// CreateUser appends a user to the configured CSV file.
func CreateUser(user *User) error {
	if user == nil {
		return errors.New("user is required")
	}

	filePath, err := getUserCSVFilePath()
	if err != nil {
		return err
	}

	if err := utils.EnsureFileExists(filePath, userCSVHeader); err != nil {
		return err
	}

	if user.ID == 0 {
		user.ID = GetNextID()
	}
	if user.CreatedAt == "" {
		user.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	row := []string{
		strconv.Itoa(user.ID),
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
	}

	return utils.AppendCSV(filePath, row)
}

// GetNextID returns the next available user ID.
func GetNextID() int {
	users, err := GetAllUsers()
	if err != nil {
		return 1
	}

	maxID := 0
	for _, user := range users {
		if user.ID > maxID {
			maxID = user.ID
		}
	}

	return maxID + 1
}

func getUserCSVFilePath() (string, error) {
	if userCSVFilePathOverride != "" {
		return userCSVFilePathOverride, nil
	}

	filePath, err := beego.AppConfig.String("csv_user_file")
	if err != nil {
		return "", err
	}

	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return "", errors.New("csv_user_file is not configured")
	}

	return filePath, nil
}

func parseUserRow(row []string, rowNumber int) (User, error) {
	if len(row) < len(userCSVHeader) {
		return User{}, fmt.Errorf("invalid user row %d: expected %d columns, got %d", rowNumber, len(userCSVHeader), len(row))
	}

	id, err := strconv.Atoi(strings.TrimSpace(row[0]))
	if err != nil {
		return User{}, fmt.Errorf("invalid user id at row %d: %w", rowNumber, err)
	}

	return User{
		ID:        id,
		Name:      row[1],
		Email:     row[2],
		Password:  row[3],
		CreatedAt: row[4],
	}, nil
}

func isEmptyCSVRow(row []string) bool {
	for _, value := range row {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}

	return true
}
