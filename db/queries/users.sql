-- name: UserUser :exec
INSERT INTO Users (
	username,
	hash,
	instagram_session_id,
	instagram_user_id,
	network
) VALUES (
	sqlc.arg(username),
	sqlc.arg(hash),
	sqlc.arg(instagram_session_id),
	sqlc.arg(instagram_user_id),
	'instagram'
);

-- name: UserUpdateInstagramSession :exec
UPDATE Users SET instagram_session_id = sqlc.arg(instagram_session_id), instagram_user_id = sqlc.arg(instagram_user_id) where username = sqlc.arg(username);

-- name: UserUpdateHash :exec
UPDATE Users SET hash = sqlc.arg(hash) where username = sqlc.arg(username);

-- name: UserCategoryAdd :exec
insert into UserCategories (username, category) values (sqlc.arg(username), sqlc.arg(category));

-- name: UserCategoryRemove :exec
delete from UserCategories where username = sqlc.arg(username) and category = sqlc.arg(category);

-- name: UserGet :one
SELECT * FROM Users WHERE username = sqlc.arg(username);

-- name: UserCategoriesGet :one
SELECT category FROM UserCategories WHERE username = sqlc.arg(username);