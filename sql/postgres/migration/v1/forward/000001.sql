-- DROP SCHEMA IF EXISTS v1 CASCADE;

CREATE SCHEMA IF NOT EXISTS v1;

CREATE TABLE IF NOT EXISTS v1.author (
	author_id serial4 NOT NULL,
	created_ts timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	first_name text NOT NULL,
	last_name text NOT NULL,
	CONSTRAINT author_pk PRIMARY KEY (author_id)
);

CREATE TABLE IF NOT EXISTS v1.collection (
	title text NOT NULL,
	created_ts timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT collection_pk PRIMARY KEY (title)
);

CREATE TABLE IF NOT EXISTS v1.book (
	book_id serial4 NOT NULL,
	author_id serial4 NOT NULL,
	created_ts timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	title text NOT NULL,
	publish_date timestamp NULL,
	edition int2 NULL,
	description text NULL,
	genre_id serial4 NOT NULL,
	CONSTRAINT book_pk PRIMARY KEY (book_id),
	CONSTRAINT book_author_fk FOREIGN KEY (author_id) REFERENCES v1.author(author_id)
);

CREATE TABLE IF NOT EXISTS v1.collection_books (
	title text NOT NULL,
	book_id serial4 NOT NULL,
	CONSTRAINT collection_books_un UNIQUE (title, book_id),
	CONSTRAINT collection_books_book_fk FOREIGN KEY (book_id) REFERENCES v1.book(book_id) on delete cascade,
	CONSTRAINT collection_books_collection_fk FOREIGN KEY (title) REFERENCES v1.collection(title) on delete cascade
);
