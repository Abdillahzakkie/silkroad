package models

import (
	"errors"

	"github.com/abdillahzakkie/silkroad/hmac"
	"github.com/abdillahzakkie/silkroad/rand"
	"github.com/jackc/pgconn"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const passwordPepper = "test-pepper"
const hmacSecretKey = "secret-hmac-key"

// UserDB is used to interact with the users table
// For Single user queries, any error but not ErrNotFound
// should be handled as 500 ErrInternalServerError
type UserDB interface {
	// Methods for querying users
	GetUserById(id uint) (User, error)
	GetAllUsers() ([]User, error)
	GetUserByRememberHash(token string) (User, error)
	GenerateRememberToken() (string, error)

	// Methods for altering users
	CreateNewUser(user *User) error
	UpdateUser(user User) error
	DeleteUserById(id uint) error

	// Used to close a DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error

	// helper functions
	getUser(user *User) error
	verifyHashedPassword(password, passwordHash string) error
	isExistingUser(user User) error
	hashPassword(user *User) error
}

// UserService is a set of methods used to
// work with user model
type UserService interface {
	// Authenticate will verify the provided email and password correct
	// if they are correct the user corresponding to that email will be returned
	// otherwise it will  return either:
	// ErrNotFound, ErrInvalidCredentials or another error
	// if something goes bad
	Authenticate(email, password string) (User, error)
	UserDB
}

var _ UserDB = &userService{}
type userService struct {
	UserDB
}

var _ UserDB = &userValidator{}
type userValidator struct {
	UserDB
}

var _ UserDB = &userGorm{}
type userGorm struct {
	db *gorm.DB
	hmac hmac.HMAC
}

/*
	NewUserService creates a new user service
	returns a pointer to the newly created UserService instance as it's first value
	and an error if one occurred as it's second value
*/
func NewUserService(psqlInfo string) (UserService, error) {
	ug, err := newUserGorm(psqlInfo); if err != nil {
		return nil, err
	}

	// auto migrate table
	if err := ug.AutoMigrate(); err != nil {
		return nil, err
	}

	us := userService{
		UserDB: &userValidator {
			UserDB: ug,
		},
	}

	return &us, nil
}

/*
	NewUserService creates a new user service
	returns a pointer to the newly created UserService instance as it's first value
	and an error if one occurred as it's second value
*/
func newUserGorm(psqlInfo string) (*userGorm, error) {
	db, err := ConnectDatabase(psqlInfo); if err != nil {
		return nil, err
	}

	ug := userGorm{
		hmac: hmac.NewHMAC(hmacSecretKey),
		db: db,
	}

	return &ug, nil
}

/*
	Close closes the database connection
	returns error if unable to close connection
*/
func (us *userGorm) Close() error {
	sql, err := us.db.DB(); if err != nil {
		return err
	}
	sql.Close()
	return nil
}

/*
	AutoMigrate auto migrates the users table
	returns error if unable to create table
*/
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}); err != nil {
		return errors.New("error while creating new 'users' table")
	}
	return nil
}

/*
	DestructiveReset drops the users table and creates a new one
	returns error if unable to drop table or re-create it
*/
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.Migrator().DropTable("users"); err != nil {
		return errors.New("unable to delete 'users' records")
	}
	// create new tables and index
	if err := ug.AutoMigrate(); err != nil {
		return errors.New("error while creating new 'users' table")
	}
	return nil
}

/*
	Authenticate authenticates the user and returns the user object
	returns error if unable to authenticate user
*/
func (us *userService) Authenticate(email, password string) (User, error) {
	user := User {
		Email: email,
	}

	// find user by  Username or Email
	if err := us.getUser(&user); err != nil {
		switch err {
			case ErrNotFound:
				return User{}, ErrInvalidCredentials
			default:
				return User{}, err
		}
	}
	// verify user's password
	if err := us.verifyHashedPassword(password, user.PasswordHash); err != nil {
		return User{}, ErrInvalidCredentials
	}
	return user, nil
}

/*
	GetUserById gets user by id
	returns error if unable to get user
*/
func (ug *userGorm) GetUserById(id uint) (User, error) {
	user := User { ID: id }

	if err := ug.getUser(&user); err != nil {
		return User{}, err
	}
	return user, nil
}

/*
	GetAllUsers gets all users
	returns error if unable to get users
*/
func (ug *userGorm) GetAllUsers() ([]User, error) {
	var users []User
	if err := ug.db.Find(&users).Error; err != nil {
		return nil, ErrInternalServerError
	}
	return users, nil
}

/*
	GetUserByRememberHash gets user by remember token
	this method will handle hashing the token and
	returns error if unable to get user
*/
func (ug *userGorm) GetUserByRememberHash(token string) (User, error) {
	hashedToken, err := ug.hmac.Hash(token); if err != nil {
		return User{}, err
	}
	user := User{
		RememberHash: hashedToken,
	}

	if err := ug.getUser(&user); err != nil {
		return User{}, err
	}

	return user, nil
}

/*
	CreateNewUser creates a new user
	returns error if unable to create user
*/
func (ug *userGorm) CreateNewUser(user *User) error {
	if err := ug.hashPassword(user); err != nil {
		return err
	}
	// creates new Remember Token
	var err error
	user.Remember, err = rand.RememberToken(); if err != nil {
		return err
	}
	// generate new RememberHash Token
	user.RememberHash, err = ug.hmac.Hash(user.Remember); if err != nil {
		return err
	}
	// save new user to database
	if err := ug.db.Create(&user).Error; err != nil {
		switch err.(type) {
			case *pgconn.PgError:
				return ErrAlreadyExists
			default:
				return ErrInternalServerError
		}
	}
	return nil
}

/*
	UpdateUser checks if user exists
	update user's credentials with the one passed in
	returns error if user does not exist
*/
func (ug *userGorm) UpdateUser(user User) error {
	if err := ug.db.Model(&user).Updates(user).Error; err != nil {
		switch err {
			case gorm.ErrRecordNotFound:
				return ErrNotFound
			default:
				return err
		}
	}
	return nil
}

/*
	DeleteUserById - deletes user by id
	returns error if unable to delete user
*/
func (ug *userGorm) DeleteUserById(id uint) error {
	// checks if user exists
	user, err := ug.GetUserById(id); if err != nil {
		return err
	}
	// delete record from database
	if err := ug.db.Where(user).Delete(&user).Error; err != nil {
		return ErrInternalServerError
	}
	return nil
}

/*
	GenerateRememberToken generates login token
	return ErrInternalServerError if an error is encountered
*/
func (ug *userGorm) GenerateRememberToken() (string, error) {
	token, err := rand.RememberToken(); if err != nil {
		return "", ErrInternalServerError
	}
	return token, nil
}

// Helpers functions
/*
	getUser gets user by the provided "user" data
	returns error if unable to get user
*/
func (ug *userGorm) getUser(user *User) error {
	// checks if user exists
	if err := ug.isExistingUser(*user); err != nil {	
		switch err {
			case ErrNotFound:
				return err
			default:
				return ErrInternalServerError
		}
	}

	if err := ug.db.Where(user).First(&user).Error; err != nil {
		return err
	}
	return nil
}

/*
	IsExistingUser checks if user exists
	returns error if user does not exist
*/
func (ug *userGorm) isExistingUser(user User) error {
	if err := ug.db.Where(user).First(&user).Error; err != nil {
		switch err {
			case gorm.ErrRecordNotFound:
				return ErrNotFound
			default:
				return err
		}
	}
	return nil
}

/*
	hashPassword hashes the provided password
	returns error if unable to hash password
*/
func (ug *userGorm) hashPassword(user *User) error {
	passwordByte := []byte(user.Password + passwordPepper)
	passwordHash, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = ""
	user.PasswordHash = string(passwordHash)
	return nil
}

/*
	verifyHashedPassword verifies the provided password against the hashed password
	returns error if unable to verify password
*/
func (ug *userGorm) verifyHashedPassword(password, passwordHash string) error {
	passwordByte := []byte(password + passwordPepper)
	password = ""
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), passwordByte); err != nil {
		switch err {
			case bcrypt.ErrMismatchedHashAndPassword:
				return ErrInvalidCredentials
			default:
				return err
		}
	}
	return nil
}