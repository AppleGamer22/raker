CREATE TABLE Users (
	username TEXT PRIMARY KEY,
	hash TEXT NOT NULL,
	instagram_session_id TEXT NOT NULL,
	instagram_user_id TEXT NOT NULL,
	network network_type NOT NULL,
	categories TEXT[]
);