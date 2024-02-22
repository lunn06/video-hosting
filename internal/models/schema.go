package models

var Schema = `
CREATE TABLE IF NOT EXISTS users (
     id TEXT PRIMARY KEY NOT NULL,
     email TEXT UNIQUE NOT NULL,
     channel_name TEXT NOT NULL,
     password TEXT NOT NULL,
     registration_time TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS roles (
     id SERIAL PRIMARY KEY,
     name TEXT UNIQUE NOT NULL,
     can_remove_users BOOL NOT NULL,
     can_remove_others_videos BOOL NOT NULL
);

CREATE TABLE IF NOT EXISTS users_roles (
	user_id TEXT REFERENCES users(id),
	role_id SERIAL REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS videos (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	localization TEXT,
	upload_time TIMESTAMP NOT NULL ,
	file_path TEXT NOT NULL,
	likes_count SERIAL,
	views_count SERIAL NOT NULL
);

CREATE TABLE IF NOT EXISTS users_videos (
	video_id TEXT REFERENCES videos(id),
	user_id TEXT REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);


CREATE TABLE IF NOT EXISTS videos_tags (
    video_id TEXT REFERENCES videos(id),
    tag_id SERIAL REFERENCES tags(id)
);
`
