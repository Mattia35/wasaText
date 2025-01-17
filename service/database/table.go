package database

var userTable = `CREATE TABLE IF NOT EXISTS userTable (
	userId INTEGER NOT NULL,
	username STRING NOT NULL UNIQUE,
	PRIMARY KEY (userId)
);`

var groupTable = `CREATE TABLE IF NOT EXISTS groupTable (
	groupId INTEGER NOT NULL,
	username STRING NOT NULL,
	PRIMARY KEY (groupId)
);`

var convTable = `CREATE TABLE IF NOT EXISTS convTable (
	convId INTEGER NOT NULL,
	groupId INTEGER,
	otherUserId INTEGER,
	senderId INTEGER NOT NULL,
	lastMessageId INTEGER,
	lastMessageConvId INTEGER,
	lastMessageSenderId INTEGER,
	PRIMARY KEY (convId, senderId)
	CONSTRAINT conversation 
		FOREIGN KEY (senderId) REFERENCES userTable (userId)
			ON DELETE CASCADE
		FOREIGN KEY (otherUserId) REFERENCES userTable (userId)
			ON DELETE CASCADE
		FOREIGN KEY (groupId) REFERENCES groupTable (groupId)
			ON DELETE CASCADE
		FOREIGN KEY (lastMessageId, lastMessageConvId, lastMessageSenderId) REFERENCES messTable (messId, convId, senderId)
			ON DELETE CASCADE
);`

var messTable = `CREATE TABLE IF NOT EXISTS messTable (
	messId INTEGER NOT NULL,
	dateTime DATETIME DEFAULT CURRENT_TIMESTAMP,
	text TEXT NOT NULL,
	status BOOLEAN,
	convId INTEGER NOT NULL,
	convSenderId INTEGER,
	senderId INTEGER NOT NULL,
	PRIMARY KEY (messId, convId, senderId)
	CONSTRAINT message
		FOREIGN KEY (convId, convSenderId) REFERENCES convTable (convId, senderId)
			ON DELETE CASCADE
		FOREIGN KEY (senderId) REFERENCES userTable (userId)
			ON DELETE CASCADE
);`

var usersGroupTable = `CREATE TABLE IF NOT EXISTS usersGroupTable (
	groupId INTEGER NOT NULL UNIQUE,
	userId INTEGER NOT NULL,
	PRIMARY KEY (groupId, userId)
	CONSTRAINT usersGroup 
		FOREIGN KEY (groupId) REFERENCES groupTable (groupId)
			ON DELETE CASCADE
		FOREIGN KEY (userId) REFERENCES userTable (userId)
			ON DELETE CASCADE
);`

var usersConvTable = `CREATE TABLE IF NOT EXISTS usersConvTable (
	convId INTEGER NOT NULL UNIQUE,
	userId INTEGER NOT NULL,
	PRIMARY KEY (convId, userId)
	CONSTRAINT usersConv
		FOREIGN KEY (convId) REFERENCES convTable (convId)
			ON DELETE CASCADE
		FOREIGN KEY (userId) REFERENCES userTable (userId)
			ON DELETE CASCADE
);`
