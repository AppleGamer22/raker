-- name: UserAdd :exec
INSERT INTO Users(
	username,
	password_hash,
	instagram_session_id,
	instagram_user_id,
	network,
	categories)
VALUES (
	sqlc.arg(username)::text,
	sqlc.arg(password_hash)::text,
	sqlc.arg(instagram_session_id)::text,
	sqlc.arg(instagram_user_id)::text,
	'instagram',
	sqlc.arg(categories)::text[]);

-- name: UserUpdateInstagramSession :exec
UPDATE
	Users
SET
	instagram_session_id = sqlc.arg(instagram_session_id)::text,
	instagram_user_id = sqlc.arg(instagram_user_id)::text,
	password_hash = sqlc.arg(password_hash)::text
WHERE
	username = sqlc.arg(username)::text;

-- name: UserUpdateHash :exec
UPDATE
	Users
SET
	password_hash = sqlc.arg(password_hash)::text
WHERE
	username = sqlc.arg(username)::text;

-- name: UserCategoryAdd :exec
UPDATE
	Users
SET
	categories = ARRAY (
		SELECT
			unnest(array_append(categories, sqlc.arg(category)::text)) AS c
		ORDER BY
			c)
	WHERE
		username = sqlc.arg(username)::text;

-- name: UserCategoryRemove :exec
UPDATE
	Users
SET
	categories = array_remove(categories, sqlc.arg(category)::text)
WHERE
	username = sqlc.arg(username)::text;

-- name: UserGet :one
SELECT
	*
FROM
	Users
WHERE
	username = sqlc.arg(username)::text;

-- name: UserGetPasskeysByUsername :many
SELECT
	*
FROM
	passkeys
WHERE
	username = sqlc.arg(username)::text;

-- name: UserCreatePasskey :exec
INSERT INTO Passkeys(
	id,
	username,
	name,
	public_key,
	attestation_type,
	aaguid,
	transports)
VALUES (
	sqlc.arg(passkey_id)::bytea,
	sqlc.arg(username)::text,
	sqlc.arg(name)::text,
	sqlc.arg(public_key)::bytea,
	sqlc.arg(attestation_type)::varchar(64),
	sqlc.arg(aaguid)::bytea,
	sqlc.arg(transports)::text[]);

-- name: PasskeyUpdateSignCount :exec
UPDATE
	Passkeys
SET
	sign_count = sign_count + 1
WHERE
	id = sqlc.arg(passkey_id)::bytea;

-- name: PasskeyUpdateName :exec
UPDATE
	Passkeys
SET
	name = sqlc.arg(name)::text
WHERE
	id = sqlc.arg(passkey_id)::bytea;

