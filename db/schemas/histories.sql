CREATE TABLE Histories (
	username TEXT REFERENCES Users(username),
	type TEXT NOT NULL,
	owner TEXT NOT NULL,
	post TEXT NOT NULL,
	PRIMARY KEY(type, owner, post),
	date TIMESTAMPTZ NOT NULL,
	files TEXT[],
	categories TEXT[]
);