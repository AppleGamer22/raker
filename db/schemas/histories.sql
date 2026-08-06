CREATE TABLE Histories(
	username text REFERENCES Users(username),
	post_type post_type NOT NULL,
	post_owner text NOT NULL,
	post text NOT NULL,
	PRIMARY KEY (username, post_type, post_owner, post),
	post_date timestamptz NOT NULL,
	files text[],
	categories text[]
);

ALTER TABLE Histories
	ADD COLUMN incognito boolean NOT NULL DEFAULT FALSE;

ALTER TABLE Histories
	ALTER COLUMN categories SET NOT NULL;

ALTER TABLE Histories
	ADD COLUMN coordinates point DEFAULT NULL;

UPDATE
	histories
SET
	coordinates = NULL;

SELECT DISTINCT
	coordinates::text
FROM
	histories;

