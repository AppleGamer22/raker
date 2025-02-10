-- name: UserUser :exec
INSERT INTO Users (
	username,
	hash,
	instagram_session_id,
	instagram_user_id,
	network,
	categories
) VALUES (
	sqlc.arg(username)::text,
	sqlc.arg(hash)::text,
	sqlc.arg(instagram_session_id)::text,
	sqlc.arg(instagram_user_id)::text,
	'instagram',
	sqlc.arg(categories)::text[]
);

-- name: UserUpdateInstagramSession :exec
UPDATE Users SET instagram_session_id = sqlc.arg(instagram_session_id)::text, instagram_user_id = sqlc.arg(instagram_user_id)::text where username = sqlc.arg(username)::text;

-- name: UserUpdateHash :exec
UPDATE Users SET hash = sqlc.arg(hash)::text where username = sqlc.arg(username)::text;

-- name: UserCategoryAdd :exec
UPDATE Users SET categories = array(
	select unnest(array_append(categories, sqlc.arg(category)::text)) AS c ORDER BY c
) where username = sqlc.arg(username)::text;

-- name: UserCategoryRemove :exec
UPDATE Users SET categories = array_remove(categories, sqlc.arg(category)::text) where username = sqlc.arg(username)::text;

-- name: UserGet :one
SELECT * FROM Users WHERE username = sqlc.arg(username)::text;