package models

var SchemaForUsers = `
CREATE TABLE IF NOT EXISTS users (
	id uuid UNIQUE,
	email TEXT UNIQUE,
	channel_name TEXT,
	password TEXT
);`
