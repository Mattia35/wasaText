package database

var userTable = `CREATE TABLE IF NOT EXISTS userTable (
	userId INTEGER NOT NULL,
	username STRING NOT NULL UNIQUE,
	photo STRING,
	PRIMARY KEY (userId)
);`

var groupTable = `CREATE TABLE IF NOT EXISTS groupTable (
	groupId INTEGER NOT NULL,
	username STRING NOT NULL,
	photo STRING,
	PRIMARY KEY (groupId)
);`

var convTable = `CREATE TABLE IF NOT EXISTS convTable (
	convId INTEGER NOT NULL,
	groupId INTEGER,
	lastMessageId INTEGER,
	PRIMARY KEY (convId)
	CONSTRAINT conversation 
		FOREIGN KEY (groupId) REFERENCES groupTable (groupId)
			ON DELETE CASCADE
);`

var messTable = `CREATE TABLE IF NOT EXISTS messTable (
	messId INTEGER NOT NULL,
	dateTime DATETIME DEFAULT CURRENT_TIMESTAMP,
	text TEXT NOT NULL,
	status BOOLEAN,
	convId INTEGER NOT NULL,
	Photo STRING,
	senderId INTEGER NOT NULL,
	PRIMARY KEY (messId, convId)
	CONSTRAINT message
		FOREIGN KEY (convId) REFERENCES convTable (convId)
			ON DELETE CASCADE
		FOREIGN KEY (senderId) REFERENCES userTable (userId)
			ON DELETE CASCADE
);`

var checkMessTable = `CREATE TABLE IF NOT EXISTS checkMessTable (
	messId INTEGER NOT NULL,
	userId INTEGER NOT NULL,
	convId INTEGER NOT NULL,
	PRIMARY KEY (messId, userId, convId)
	CONSTRAINT checkMess
		FOREIGN KEY (messId, convId) REFERENCES messTable (messId, convId)
			ON DELETE CASCADE
		FOREIGN KEY (userId) REFERENCES userTable (userId)
			ON DELETE CASCADE
);`   

var usersGroupTable = `CREATE TABLE IF NOT EXISTS usersGroupTable (
	groupId INTEGER NOT NULL,
	userId INTEGER NOT NULL,
	PRIMARY KEY (groupId, userId)
	CONSTRAINT usersGroup 
		FOREIGN KEY (groupId) REFERENCES groupTable (groupId)
			ON DELETE CASCADE
		FOREIGN KEY (userId) REFERENCES userTable (userId)
			ON DELETE CASCADE
);`

var usersConvTable = `CREATE TABLE IF NOT EXISTS usersConvTable (
	convId INTEGER NOT NULL,
	userId INTEGER NOT NULL,
	PRIMARY KEY (convId, userId)
	CONSTRAINT usersConv
		FOREIGN KEY (convId) REFERENCES convTable (convId)
			ON DELETE CASCADE
		FOREIGN KEY (userId) REFERENCES userTable (userId)
			ON DELETE CASCADE
);`

var commentTable = `CREATE TABLE IF NOT EXISTS commentTable (
	commId INTEGER NOT NULL UNIQUE,
	messId INTEGER NOT NULL,
	content TEXT NOT NULL,
	senderId INTEGER NOT NULL,
	convId INTEGER NOT NULL,
	PRIMARY KEY (messId, senderId, convId)
	CONSTRAINT message
		FOREIGN KEY (messId, convId) REFERENCES messTable (messId, convId)
			ON DELETE CASCADE
		FOREIGN KEY (senderId) REFERENCES userTable (userId)
			ON DELETE CASCADE
);`
