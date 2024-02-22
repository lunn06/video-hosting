package models

var Schema = `
CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY NOT NULL,
	email TEXT UNIQUE NOT NULL,
	channel_name TEXT NOT NULL,
	password TEXT NOT NULL,
	data_reg TIMESTAMP NOT NULL,
	role TEXT NOT NULL REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS roles (
	id UUID PRIMARY KEY,
	role TEXT UNIQUE NOT NULL,
	can_remove_users BOOL NOT NULL,
	can_remove_others_videos BOOL NOT NULL
);

CREATE TABLE IF NOT EXISTS videos (
    id UUID PRIMARY KEY,
	title TEXT NOT NULL,	
	localization TEXT,
	upload_date TIMESTAMP NOT NULL ,
	path_to_file TEXT NOT NULL,
	count_likes INTEGER,
	views INTEGER NOT NULL                                 
);

CREATE TABLE IF NOT EXISTS users_videos (
	video_id UUID REFERENCES videos(id),
	user_id UUID REFERENCES users(id)	
);

CREATE TABLE IF NOT EXISTS tags (
	id UUID,
	tag TEXT UNIQUE NOT NULL
);


CREATE TABLE IF NOT EXISTS videos_tags (
	video_id UUID REFERENCES videos(id),
	tag_id UUID REFERENCES tags(id)	
);
`
