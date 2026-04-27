-- name: HistoryAddFromArchive :one
INSERT INTO Histories(
		username,
		post_type,
		post_owner,
		post,
		post_date,
		files,
		categories
	)
VALUES (
		sqlc.arg(username)::text,
		sqlc.arg(post_type)::post_type,
		sqlc.arg(post_owner)::text,
		sqlc.arg(post)::text,
		sqlc.arg(post_date)::TIMESTAMPTZ,
		sqlc.arg(files)::text [],
		sqlc.arg(categories)::text []
	)
RETURNING *;

-- name: HistoryAdd :one
INSERT INTO Histories(
		username,
		post_type,
		post_owner,
		post,
		post_date,
		files,
		categories
	)
VALUES (
		sqlc.arg(username)::text,
		sqlc.arg(post_type)::post_type,
		sqlc.arg(post_owner)::text,
		sqlc.arg(post)::text,
		NOW(),
		sqlc.arg(files)::text [],
		sqlc.arg(categories)::text []
	)
RETURNING *;

-- name: HistoryUpdateCategories :one
UPDATE Histories
SET categories = sqlc.slice(categories)::text []
WHERE post_type = sqlc.arg(post_type)::post_type
	AND post = sqlc.arg(post)::text
	AND post_owner = sqlc.arg(post_owner)::text
	AND username = sqlc.arg(username)::text
RETURNING *;

-- name: UpdateHistoryRemoveFile :one
UPDATE Histories
SET files = array_remove(files, sqlc.arg(file)::text)
WHERE post_type = sqlc.arg(post_type)::post_type
	AND post = sqlc.arg(post)::text
	AND post_owner = sqlc.arg(post_owner)::text
	AND username = sqlc.arg(username)::text
RETURNING *;

-- name: HistoryUpdateOwner :exec
UPDATE Histories
SET post_owner = sqlc.arg(new_owner)::text
WHERE post_type = sqlc.arg(post_type)::post_type
	AND post_owner = sqlc.arg(old_owner)::text
	AND username = sqlc.arg(username)::text;

-- name: HistoriesCategoryRename :exec
UPDATE Histories
SET categories = (
		SELECT array_agg(
				c
				ORDER BY c
			)
		FROM unnest(
				array_replace(
					categories,
					sqlc.arg(old_category)::text,
					sqlc.arg(new_category)::text
				)
			) AS c
	)
WHERE username = sqlc.arg(username)::text
	AND sqlc.arg(old_category)::text = ANY(categories);

-- name: HistoryGet :one
SELECT *
FROM Histories
WHERE post_type = sqlc.arg(post_type)::post_type
	AND post = sqlc.arg(post)::text
	AND username = sqlc.arg(username)::text;

-- name: HistoryGetByOwner :one
SELECT *
FROM Histories
WHERE post_type = sqlc.arg(post_type)::post_type
	AND post_owner = sqlc.arg(post_owner)::text
	AND post = sqlc.arg(post)::text
	AND username = sqlc.arg(username)::text;

-- name: HistoryCountByFile :one
SELECT count(*)
FROM Histories
WHERE post_type = sqlc.arg(post_type)::post_type
	AND post_owner = sqlc.arg(post_owner)::text
	AND sqlc.arg(file)::text = ANY(files)
	AND username = sqlc.arg(username)::text;

-- https://docs.sqlc.dev/en/stable/howto/select.html#passing-a-slice-as-a-parameter-to-a-query
-- https://docs.sqlc.dev/en/stable/howto/named_parameters.html
-- name: HistoryGetInclusive :many
SELECT *
FROM Histories
WHERE post_type = ANY (sqlc.slice(post_types)::post_type [])
	AND categories <@ sqlc.slice(categories)::text []
	AND post_owner LIKE sqlc.arg(post_owner)::text
	AND username = sqlc.arg(username)::text
LIMIT sqlc.arg(page_size)::int OFFSET sqlc.arg(page)::int;

-- name: HistoryGetExclusive :many
SELECT *
FROM Histories
WHERE post_type = ANY (sqlc.slice(post_types)::post_type [])
	AND categories = sqlc.slice(categories)::text []
	AND post_owner LIKE sqlc.arg(post_owner)::text
	AND username = sqlc.arg(username)::text
LIMIT sqlc.arg(page_size)::int OFFSET sqlc.arg(page)::int;

-- name: HistoryGetPage :many
SELECT *
FROM Histories
WHERE post_type = ANY (sqlc.slice(post_types)::post_type [])
	AND (
		(
			sqlc.arg(exclusive)::boolean
			and categories = sqlc.slice(categories)::text []
		)
		or (
			not sqlc.arg(exclusive)::boolean
			and categories <@ sqlc.slice(categories)::text []
		)
	)
	AND post_owner LIKE FORMAT('%%%s%%', sqlc.arg(post_owner)::text)
	AND username = sqlc.arg(username)::text
order by post_date DESC
LIMIT sqlc.arg(page_size)::int OFFSET sqlc.arg(page)::int;

-- name: HistoryCount :one
select count(*)
from Histories
WHERE post_type = ANY (sqlc.slice(post_types)::post_type [])
	AND (
		(
			sqlc.arg(exclusive)::boolean
			and categories = sqlc.slice(categories)::text []
		)
		or (
			not sqlc.arg(exclusive)::boolean
			and categories <@ sqlc.slice(categories)::text []
		)
	)
	AND post_owner LIKE FORMAT('%%%s%%', sqlc.arg(post_owner)::text)
	AND username = sqlc.arg(username)::text;

-- name: HistoryOwners :many
select distinct post_owner
from Histories
WHERE post_type = ANY (sqlc.slice(post_types)::post_type [])
	AND (
		(
			sqlc.arg(exclusive)::boolean
			and categories = sqlc.slice(categories)::text []
		)
		or (
			not sqlc.arg(exclusive)::boolean
			and categories <@ sqlc.slice(categories)::text []
		)
	)
	AND post_owner LIKE FORMAT('%%%s%%', sqlc.arg(post_owner)::text)
	AND username = sqlc.arg(username)::text;

-- name: HistoryRemove :exec
DELETE FROM Histories
WHERE post_type = sqlc.arg(post_type)::post_type
	AND post_owner = sqlc.arg(post_owner)::text
	AND post = sqlc.arg(post)::text
	AND username = sqlc.arg(username)::text;