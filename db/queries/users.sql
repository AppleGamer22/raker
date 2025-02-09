-- name: UserUser :exec
INSERT INTO Users (
	username,
	hash,
	instagram_session_id,
	instagram_user_id,
	network,
	categories
) VALUES ($1, $2, $3, $4, 'instagram', $5);

-- name: UserUpdateInstagramSession :exec
UPDATE Users SET instagram_session_id = $2, instagram_user_id = $3 where username = $1;

-- name: UserUpdateHash :exec
UPDATE Users SET hash = $2 where username = $1;

-- name: UserCategoryAdd :exec
UPDATE Users SET categories = array(
	select unnest(array_append(categories, $2)) AS c ORDER BY c
) where username = $1;

-- name: UserCategoryRemove :exec
UPDATE Users SET categories = array_remove(categories, $2) where username = $1;

-- name: UserGet :one
SELECT * FROM Users WHERE username = $1;