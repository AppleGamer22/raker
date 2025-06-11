CREATE TABLE Users (
	username TEXT PRIMARY KEY,
	hash TEXT NOT NULL,
	instagram_session_id TEXT NOT NULL,
	instagram_user_id TEXT NOT NULL,
	network ENUM ('instagram', 'tiktok', 'vsco')
);

create table UserCategories (
	username TEXT,
	foreign key (username) REFERENCES Users(username),
	category TEXT,
	primary key (username, category)
);