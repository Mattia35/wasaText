/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"progetto.wasa/service/api/structions"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	NameControl(username string) (bool, error)
	GetUserByName(username string) (structions.User, error)
	GetGroupByGroupId(groupId int) (structions.Group, error)
	CreateUser(u structions.User) (structions.User, error)
	UsernameModify(userId int, username string) error
	UserControlByGroup(userId int, groupId int) (bool, error)
	GroupnameModify(groupId int, groupname string) error
	CreateGroup(gr structions.Group, userId int) (structions.Group, error)
	CreateConversation(conv structions.Conversation) (structions.Conversation, error)
	UserControlByUsername(username string) (structions.User, error)
	AddUserToConv(userId int, convId int) error
	AddUserToGroup(userId int, groupId int) error
	CreateMessage(mes structions.Message) (structions.Message, error)
	AddMessageToConv(MessageId int, ConvId int) error
	GetUserById(userId int) (structions.User, error)
	GetConvByUsers(userId int, destId int) (int, error)
	GetConversationsByUserId(userId int) ([]structions.Conversation, error)
	GetConvById(convId int) (structions.Conversation, error)
	GetUserByConv(convId int, userId int) (structions.User, error)
	GetMessageById(LastMessage int, ConvId int) (structions.Message, error)
	GetUsersByGroupId(groupId int) ([]structions.User, error)
	IsUserInGroup(userId int, groupId int) (bool, error)
	IsUserInConv(userId int, convId int) (bool, error)
	RemoveUserFromGroup(userId int, groupId int) error
	GetConvByGroupId(groupId int) (structions.Conversation, error)
	RemoveUserFromConv(userId int, convId int) error
	RemoveGroup(groupId int) error
	RemoveConv(convId int) error
	GetMaxMessageId(convId int) (int, error)
	RemoveMessage(convId int) error
	CheckMessageSender(messId int, userId int) (bool, error)
	GetMessagesByConvId(convId int) ([]structions.Message, error)
	AddUserToListOfReadersOfMess(messId int, userId int, convId int) error
	CheckIfUserHasReadMess(messId int, userId int) (bool, error)
	CheckAllUsersHaveReadMess(messId int) (bool, error)
	UpdateMessageStatus(messId int) error
	GetCommentsByMessId(messId int, convId int) ([]structions.Comment, error)
	GetUsersByConvId(convId int) ([]structions.User, error)
	AddUserToListOfAlreadyReadersOfMess(messId int, userId int, convId int) error
	CheckIfUserHasAlreadyCommented(messId, userId, convId int) (bool, error)
	CreateComment(com structions.Comment) (structions.Comment, error)
	UpdateComment(userId int, messId int, convId int, emoji string) error
	IsCommentInConv(commId int, messId int, convId int) (bool, error)
	GetCommentById(commId int, messId int, convId int) (structions.Comment, error)
	RemoveComment(commId int, messId int, convId int) error
	SetGroupPhoto(groupId int, base64 string) error
	SetUserPhoto(userId int, base64 string) error
	SearchUsers(query string) ([]structions.User, error)
	DoConversationExist(ConvId int) (bool, error)
	GetConversationByUsers(userId int, destId int) (int, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}
	/// Check if the database is empty
	var tableSQL uint8
	err := db.QueryRow("SELECT COUNT(name) FROM sqlite_master WHERE type='table'").Scan(&tableSQL)
	if err != nil {
		return nil, fmt.Errorf("error checking if database is empty: %w", err)
	}

	// Check of the number of table is corret (there are 5 tables)
	// if the number of table is not 5, we creating missing tables
	if tableSQL != 8 {

		// Creation of the user tabel
		_, err = db.Exec(userTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure user: %w", err)
		}

		// Creation of the message table
		_, err = db.Exec(messTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure message: %w", err)
		}

		// Creation of the group table
		_, err = db.Exec(groupTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure group: %w", err)
		}

		// Creation of the user_group table
		_, err = db.Exec(usersGroupTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure user and group: %w", err)
		}

		// Creation of the user_conv table
		_, err = db.Exec(usersConvTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure user and conv: %w", err)
		}

		// Creation of the conversation table
		_, err = db.Exec(convTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure conversation: %w", err)
		}

		// Creation of the check_message table
		_, err = db.Exec(checkMessTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure check message: %w", err)
		}

		// Creation of the comment table
		_, err = db.Exec(commentTable)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure check message: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
