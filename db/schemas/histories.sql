CREATE TABLE Histories (
	username TEXT,
	foreign key (username) REFERENCES Users(username),
	type ENUM (
		'instagram',
		'highlight',
		'story',
		'tiktok',
		'vsco'
	),
	owner TEXT NOT NULL,
	post TEXT NOT NULL,
	PRIMARY KEY(username, type, owner, post),
	date timestamp NOT NULL
);

create table HistoryFiles (
	username TEXT,
	foreign key (username) REFERENCES Users(username),
	type ENUM (
		'instagram',
		'highlight',
		'story',
		'tiktok',
		'vsco'
	),
	owner TEXT,
	post TEXT,
	FOREIGN KEY (owner, post) REFERENCES Histories(owner, post),
	file TEXT,
	primary key (username, type, owner, post, file)
);

create table HistoryCategories (
	username TEXT,
	FOREIGN KEY (username) REFERENCES UserCategories(username),
	type ENUM (
		'instagram',
		'highlight',
		'story',
		'tiktok',
		'vsco'
	),
	owner TEXT,
	post TEXT,
	FOREIGN KEY (owner, post) REFERENCES Histories(owner, post),
	category TEXT,
	FOREIGN KEY (category) REFERENCES UserCategories(category),
	primary key (username, type, owner, post, category)
);