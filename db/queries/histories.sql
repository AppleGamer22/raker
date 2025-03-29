-- name: HistoryAdd :one
INSERT INTO Histories (
	username,
	type,
	owner,
	post,
	date,
	files,
	categories
) VALUES (
	sqlc.arg(username)::text,
	sqlc.arg(type)::post_type,
	sqlc.arg(owner)::text,
	sqlc.arg(post)::text,
	NOW(),
	sqlc.arg(files)::text[],
	sqlc.arg(categories)::text[]
) RETURNING *;

-- name: HistoryUpdateCategories :exec
UPDATE Histories
SET categories = sqlc.slice(categories)::text[]
WHERE post = sqlc.arg(type)::post_type
	AND post = sqlc.arg(post)::text
	AND username = sqlc.arg(username)::text;

-- name: UpdateHistoryRemoveFile :exec
UPDATE Histories
SET files = array_remove(files, sqlc.arg(file)::text)
WHERE post = sqlc.arg(type)::post_type
	AND post = sqlc.arg(post)::text
	AND username = sqlc.arg(username)::text;

-- name: HistoryUpdateOwner :exec
UPDATE Histories
SET owner = sqlc.arg(old_owner)::text
WHERE post = sqlc.arg(type)::post_type
	AND owner = sqlc.arg(new_owner)::text
	AND username = sqlc.arg(username)::text;

-- name: HistoryGet :one
SELECT * FROM Histories
WHERE type = sqlc.arg(type)::post_type
	AND post = sqlc.arg(post)::text
	AND username = sqlc.arg(username)::text;

-- https://docs.sqlc.dev/en/stable/howto/select.html#passing-a-slice-as-a-parameter-to-a-query
-- https://docs.sqlc.dev/en/stable/howto/named_parameters.html
-- name: HistoryGetInclusive :many
SELECT * FROM Histories
WHERE type = ANY(sqlc.slice(types)::post_type[])
	AND categories <@ sqlc.slice(categories)::text[]
	AND OWNER LIKE sqlc.arg(owner)::text
	AND username = sqlc.arg(username)::text;

-- name: HistoryGetExclusive :many
SELECT * FROM Histories
WHERE type = ANY(sqlc.slice(types)::post_type[])
	AND categories = sqlc.slice(categories)::text[]
	AND OWNER LIKE sqlc.arg(owner)::text
	AND username = sqlc.arg(username)::text;

-- name: HistoryRemove :exec
DELETE FROM Histories
WHERE type = sqlc.arg(type)::post_type
	AND owner = sqlc.arg(owner)::text
	AND post = sqlc.arg(post)::text
	AND username = sqlc.arg(username)::text;