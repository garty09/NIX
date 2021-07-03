CREATE DATABASE JSONPlaceholder;

CREATE TABLE public_posts
(
    user_id integer NOT NULL,
    id integer NOT NULL PRIMARY KEY,
    title text NOT NULL,
    body text NOT NULL

);

CREATE TABLE public_comments
(
    post_id integer NOT NULL,
    id integer NOT NULL PRIMARY KEY,
    name varchar(100) NOT NULL,
    email varchar(40) NOT NULL,
    body text NOT NULL
);

