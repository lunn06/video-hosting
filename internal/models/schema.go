package models

var Schema = `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	email TEXT UNIQUE NOT NULL,
	channel_name TEXT UNIQUE NOT NULL,
	password TEXT NOT NULL,
	registration_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS jwt_tokens (
	uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	token TEXT UNIQUE,
	creation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users_tokens (
	user_id INTEGER UNIQUE REFERENCES users(id),
	token_uuid UUID UNIQUE REFERENCES jwt_tokens(uuid)
);

CREATE TABLE IF NOT EXISTS roles (
	id INTEGER PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	can_remove_users BOOL NOT NULL,
	can_remove_others_videos BOOL NOT NULL
);

CREATE TABLE IF NOT EXISTS users_roles (
	user_id INTEGER UNIQUE REFERENCES users(id),
	role_id INTEGER REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS videos (
	uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	title TEXT NOT NULL,
	localization TEXT NOT NULL,
	upload_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	file_path TEXT UNIQUE NOT NULL,
	likes_count INTEGER NOT NULL,
	views_count INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS users_videos (
	user_id SERIAL REFERENCES users(id),
	video_uuid UUID UNIQUE REFERENCES videos(uuid)
);

CREATE TABLE IF NOT EXISTS tags (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS videos_tags (
	video_uuid UUID REFERENCES videos(uuid),
	tag_id SERIAL REFERENCES tags(id)
);
`
