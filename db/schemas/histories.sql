CREATE TABLE Histories (
	username TEXT REFERENCES Users(username),
	post_type post_type NOT NULL,
	post_owner TEXT NOT NULL,
	post TEXT NOT NULL,
	PRIMARY KEY(username, post_type, post_owner, post),
	post_date TIMESTAMPTZ NOT NULL,
	files TEXT [],
	categories TEXT []
);