CREATE TYPE post_type AS ENUM(
	'instagram',
	'highlight',
	'story',
	'tiktok',
	'vsco'
);

ALTER TYPE post_type
	ADD VALUE 'snapchat';

ALTER TYPE post_type
	ADD VALUE 'threads';

CREATE TYPE network_type AS ENUM(
	'instagram',
	'tiktok',
	'vsco'
);

