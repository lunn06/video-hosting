package models

// TODO: подумать над связями и какие еще таблицы могут быть

var Schema = `
CREATE TABLE IF NOT EXISTS users (
	user_id UUID PRIMARY KEY NOT NULL,
	email TEXT UNIQUE NOT NULL,
	channel_name TEXT NOT NULL,
	password TEXT NOT NULL,
	data_reg TIMESTAMP,
	jwt TEXT NOT NULL,
	role TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS videos (
	video_id UUID PRIMARY KEY,
	user_id UUID REFERENCES users(user_id)	
);

CREATE TABLE IF NOT EXISTS video_attributes (
    video_id UUID REFERENCES videos(video_id),
	title TEXT NOT NULL,	
	tag TEXT,
	localization TEXT,
	upload_date TIMESTAMP NOT NULL ,
	path_to_file TEXT NOT NULL,
	count_likes INTEGER,
	views INTEGER NOT NULL                                 
);`
