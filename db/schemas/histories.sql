CREATE OR ALTER TABLE Histories (
	username TEXT REFERENCES Users(username),
	[type] TEXT NOT NOT,
	[owner] TEXT NOT NULL,
	post TEXT NOT NULL,
	PRIMARY KEY([type], [owner], [post]),
	[date] TIMESTAMPTZ NOT NULL,
	urls TEXT[],
	categories TEXT[]
);