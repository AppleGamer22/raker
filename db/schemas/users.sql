CREATE TABLE Users (
	username TEXT PRIMARY KEY,
	password_hash TEXT NOT NULL,
	instagram_session_id TEXT NOT NULL,
	instagram_user_id TEXT NOT NULL,
	network network_type NOT NULL,
	categories TEXT []
);

alter table Users
add column tiktok_session_id TEXT NOT NULL DEFAULT '';

alter table Users
add column tiktok_session_id_guard TEXT NOT NULL DEFAULT '';