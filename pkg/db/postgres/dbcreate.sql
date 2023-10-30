DROP TABLE IF EXISTS comments;

CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY, -- первичный ключ
    postid BIGINT DEFAULT 0,
    commentid BIGINT DEFAULT 0,
	content TEXT
);
