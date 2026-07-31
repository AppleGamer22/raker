CREATE TABLE Users(
	username text PRIMARY KEY,
	password_hash text NOT NULL,
	instagram_session_id text NOT NULL,
	instagram_user_id text NOT NULL,
	network network_type NOT NULL,
	categories text[]
);

ALTER TABLE Users
	ADD COLUMN tiktok_session_id text NOT NULL DEFAULT '';

ALTER TABLE Users
	ADD COLUMN tiktok_session_id_guard text NOT NULL DEFAULT '';

ALTER TABLE Users
	ADD COLUMN id uuid NOT NULL UNIQUE DEFAULT gen_random_uuid();

CREATE TABLE Passkeys(
	id bytea PRIMARY KEY,
	name text NOT NULL,
	user_id uuid NOT NULL REFERENCES Users(id) ON DELETE CASCADE,
	public_key bytea NOT NULL,
	attestation_type varchar(64) NOT NULL,
	aaguid bytea NOT NULL,
	sign_count bigint NOT NULL DEFAULT 0,
	transports text[] NOT NULL DEFAULT '{}'
);

