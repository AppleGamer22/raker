CREATE TYPE post_type AS ENUM (
	'instagram',
	'highlight',
	'story',
	'tiktok',
	'vsco'
);

ALTER TYPE post_type ADD VALUE 'snapchat';

CREATE TYPE network_type AS ENUM ('instagram', 'tiktok', 'vsco');