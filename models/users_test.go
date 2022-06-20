package models

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func errorResponse(values ...interface{}) string {
	return fmt.Sprintf("Expected %s, got %s", values[0], values[1])
}

func getUser() *User {
	return &User{
		Username: "silkroad",
		Wallet: "0x0000000000000000000000000000000000000000",
		Email: "example@example.com",
		Password: "test",
	}
}

func testingUserService() (*UserService, error) {
	psqlInfo, err := GetPsqlInfo("silkroad_test"); if err != nil {
		return nil, err
	}

	us, err := NewUserService(psqlInfo); if err != nil {
		return nil, err
	}
	// clear the users table between test
	err = us.DestructiveReset(); if err != nil {
		return us, err
	}
	return us, nil
}

func createNewUser(us *UserService, user *User) error {
	err := us.CreateNewUser(user); if err != nil {
		return err
	}
	return nil
}

func TestHashPassword(t *testing.T) {
	us, err := testingUserService(); if err != nil {
		t.Errorf(err.Error())
	}
	// close database connection
	defer us.Close()
	// create new user
	user := getUser()
	password := user.Password

	if err := us.HashPassword(user); err != nil {
		t.Error(err)
	}
	// assert user.Password is cleared
	if user.Password != "" {
		t.Error(errorResponse("password to be null", user.Password))
	}
	// assert that user.PasswordHash is hashed correctly
	if user.PasswordHash == password {
		t.Error(errorResponse("hashed password", "plain text password"))
	}
}

func TestVerifyHashedPassword(t *testing.T) {
	us, err := testingUserService(); if err != nil {
		t.Errorf(err.Error())
	}
	// close database connection
	defer us.Close()
	// create new user
	user := getUser()
	password := user.Password

	if err := us.HashPassword(user); err != nil {
		t.Error(err)
	}

	// assert that user.PasswordHash is hashed correctly
	if err := us.VerifyHashedPassword(password, user.PasswordHash); err != nil {
		t.Error(err)
	}
}

func TestCreateUser(t *testing.T) {
	us, err := testingUserService(); if err != nil {
		panic(err)
	}
	// close database connection
	defer us.Close()

	// create new user
	user := getUser()
	if err := createNewUser(us, user); err != nil {
		t.Error(err)
	}
	// assert that user.ID is 1
	if user.ID == 0 {
		t.Errorf(errorResponse(1, user.ID))
	}
	// assert that time.CreatedAt is recent and less than 5 seconds
	if time.Since(user.CreatedAt) > time.Duration(5 * time.Second) {
		t.Error(errorResponse("soon", user.CreatedAt))
	}
	// assert that time.UpdatedAt is recent and less than 5 seconds
	if time.Since(user.UpdatedAt) > time.Duration(5 * time.Second) {
		t.Error(errorResponse("soon", user.UpdatedAt))
	}
	// assert that user.Remember is set properly
	if user.Remember == "" {
		t.Error(errorResponse("", user.Remember))
	}
	// assert that user.RememberHash is set properly
	if user.RememberHash == "" {
		t.Error(errorResponse("", user.RememberHash))
	}

}

func TestGetAllUsers(t *testing.T) {
	us, err := testingUserService(); if err != nil {
		t.Errorf(err.Error())
	}
	// close database connection
	defer us.Close()
	
	// create new user
	if err := createNewUser(us, getUser()); err != nil {
		t.Errorf(err.Error())
	}

	users, err := us.GetAllUsers(); if err != nil {
		t.Errorf(err.Error())
	}

	if len(users) != 1 {
		t.Errorf("Expected users of array of length: 1, got users of length: %d", len(users))
	}
}

func TestDeleteUserById(t *testing.T) {
	us, err := testingUserService(); if err != nil {
		t.Errorf(err.Error())
	}
	// close database connection
	defer us.Close()

	// create new user
	user := User{
		Username: "silkroad",
		Wallet: "0x0000000000000000000000000000000000000000",
		Email: "example@test.com",
		Password: "test",
	}

	if err := us.CreateNewUser(&user); err != nil {
		t.Errorf(err.Error())
	}

	// delete newly created user
	if err := us.DeleteUserById(user.ID); err != nil {
		t.Errorf(err.Error())
	}

	// assert that user does not exist
	_, err = us.GetUserById(user.ID); if err != nil {
		if !errors.Is(err, ErrNotFound) {
			t.Errorf(err.Error())
		}
	}
}

func TestIsExistingUser(t *testing.T) {
	us, err := testingUserService(); if err != nil {
		t.Errorf(err.Error())
	}
	// close database connection
	defer us.Close()
	user := User{
		Username: "silkroad",
		Wallet: "0x0000000000000000000000000000000000000000",
		Email: "example@test.com",
		Password: "test",
	}
	// create new user
	if err := us.CreateNewUser(&user); err != nil {
		t.Errorf(err.Error())
	}

	// should pass since user has already been created above
	if err := us.IsExistingUser(user); err != nil {
		t.Error(err)
	}

	// should return err for arbitrary user ID
	if err := us.IsExistingUser(User{ID: 1000}); !errors.Is(err, ErrNotFound) {
		t.Errorf(errorResponse(ErrNotFound.Error(), err.Error()))
	}
}

func TestAuthenticate(t *testing.T) {
	us, err := testingUserService(); if err != nil {
		t.Errorf(err.Error())
	}
	// close database connection
	defer us.Close()

	// create new user
	user := User{
		Username: "silkroad",
		Wallet: "0x0000000000000000000000000000000000000000",
		Email: "example@example.com",
		Password: "test",
	}
	password := user.Password
	if err := createNewUser(us, &user); err != nil {
		t.Error(err)
	}
	// authenticate user
	_, err = us.Authenticate(user.Email, password); if err != nil {
		t.Errorf(err.Error())
	}
}