-- name: UserAdd :one
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
	sqlc.arg(categories)::text[])
RETURNING
	*;

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

-- name: UserGetByUsername :one
SELECT
	*
FROM
	Users
WHERE
	username = sqlc.arg(username)::text;

-- name: UserGetByID :one
SELECT
	*
FROM
	Users
WHERE
	id = sqlc.arg(user_id);

-- name: UserGetPasskeysByID :many
SELECT
	*
FROM
	passkeys
WHERE
	user_id = sqlc.arg(user_id);

-- name: UserCreatePasskey :exec
INSERT INTO Passkeys(
	id,
	user_id,
	name,
	public_key,
	attestation_type,
	aaguid,
	sign_count,
	transports,
	backup_eligible,
	backup_state)
VALUES (
	sqlc.arg(passkey_id),
	sqlc.arg(user_id),
	sqlc.arg(name),
	sqlc.arg(public_key),
	sqlc.arg(attestation_type),
	sqlc.arg(aaguid),
	sqlc.arg(sign_count),
	sqlc.arg(transports),
	sqlc.arg(backup_eligible),
	sqlc.arg(backup_state));

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

