-- name: HistoryAdd :exec
INSERT INTO Histories (
	username,
	type,
	owner,
	post,
	date
) VALUES (
	sqlc.arg(username),
	sqlc.arg(type),
	sqlc.arg(owner),
	sqlc.arg(post),
	NOW()
);

-- name: HistoryAddCategory :exec
insert into HistoryCategories (
	username,
	type,
	owner,
	post,
	category
) values (
	sqlc.arg(username),
	sqlc.arg(type),
	sqlc.arg(owner),
	sqlc.arg(post),
	sqlc.arg(category)
);

-- name: HistoryRemoveCategory :exec
delete from HistoryCategories
WHERE post = sqlc.arg(post)
	AND post = sqlc.arg(post)
	AND owner = sqlc.arg(owner)
	AND category = sqlc.arg(category);

-- name: HistoryAddFile :exec
insert into HistoryFiles (
	username,
	type,
	owner,
	post,
	file
) values (
	sqlc.arg(username),
	sqlc.arg(type),
	sqlc.arg(owner),
	sqlc.arg(post),
	sqlc.arg(file)
);

-- name: UpdateHistoryRemoveFile :exec
delete from HistoryFiles
WHERE post = sqlc.arg(post)
	AND post = sqlc.arg(post)
	AND owner = sqlc.arg(owner)
	AND username = sqlc.arg(username);

-- name: HistoryUpdateOwner :exec
UPDATE Histories
SET owner = sqlc.arg(old_owner)
WHERE post = sqlc.arg(post)
	AND owner = sqlc.arg(new_owner)
	AND username = sqlc.arg(username);

-- name: HistoryGet :one
SELECT * FROM Histories
WHERE type = sqlc.arg(type)
	AND post = sqlc.arg(post)
	AND owner = sqlc.arg(owner)
	AND username = sqlc.arg(username)
limit 30;

-- https://docs.sqlc.dev/en/stable/howto/select.html#passing-a-slice-as-a-parameter-to-a-query
-- https://docs.sqlc.dev/en/stable/howto/named_parameters.html
-- name: HistoryGetInclusive :many
SELECT * FROM Histories inner JOIN HistoryCategories
	on Histories.username = HistoryCategories.username
	and Histories.type = HistoryCategories.type
	and Histories.owner = HistoryCategories.owner
	and Histories.post = HistoryCategories.post
WHERE Histories.type in (sqlc.slice(types))
	AND category in (sqlc.slice(categories))
	AND Histories.owner like sqlc.arg(owner)
	AND Histories.username = sqlc.arg(username)
limit 30;

-- name: HistoryGetExclusive :many
SELECT * FROM Histories inner JOIN HistoryCategories
	on Histories.username = HistoryCategories.username
	and Histories.type = HistoryCategories.type
	and Histories.owner = HistoryCategories.owner
	and Histories.post = HistoryCategories.post
where Histories.type in (sqlc.slice(types))
	AND category in (sqlc.slice(categories))
	AND Histories.owner like sqlc.arg(owner)
	AND Histories.username = sqlc.arg(username)
GROUP BY HistoryCategories.post
having count(HistoryCategories.post) = sqlc.arg(category_count)
limit 30;

-- name: HistoryRemove :exec
DELETE FROM Histories
WHERE type = sqlc.arg(type)
	AND owner = sqlc.arg(owner)
	AND post = sqlc.arg(post)
	AND username = sqlc.arg(username);