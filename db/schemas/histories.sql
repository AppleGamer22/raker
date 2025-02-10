CREATE TABLE Histories (
	username TEXT REFERENCES Users(username),
	type post_type NOT NULL,
	owner TEXT NOT NULL,
	post TEXT NOT NULL,
	PRIMARY KEY(type, owner, post),
	date TIMESTAMPTZ NOT NULL,
	files TEXT[],
	categories TEXT[]
);