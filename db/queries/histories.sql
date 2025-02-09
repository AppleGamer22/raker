-- name: HistoryAdd :exec
INSERT INTO Histories (
	username,
	type,
	owner,
	post,
	date,
	files,
	categories
) VALUES ($1, $2, $3, $4, NOW(), $5, $6);

-- name: HistoryUpdateCategories :exec
UPDATE Histories SET categories = $3 WHERE type = $1 AND post = $2;

-- name: UpdateHistoryRemoveFile :exec
UPDATE Histories SET files = array_remove(files, $3) WHERE type = $1 AND post = $2;

-- name: HistoryUpdateOwner :exec
UPDATE Histories SET owner = $3 WHERE type = $1 AND owner = $2;

-- name: HistoryGet :one
SELECT * FROM Histories WHERE type = $1 AND post = $2;

-- https://github.com/sqlc-dev/sqlc/issues/895#issuecomment-785553456
-- name: HistoryGetInclusive :many
SELECT * FROM Histories WHERE type = ANY($1::TEXT[]) AND categories <@ $2 AND OWNER LIKE $3;

-- name: HistoryGetExclusive :many
SELECT * FROM Histories WHERE type = ANY($1::TEXT[]) AND categories = $2 AND OWNER LIKE $3;

-- name: HistoryRemove :exec
DELETE FROM Histories where type = $1 AND owner = $2 AND post = $3;