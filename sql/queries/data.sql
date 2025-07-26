-- name: InsertData :exec
INSERT INTO data (url, content, created_at, updated_at) VALUES (
	?,
	?,
	datetime('now'),
	datetime('now')
);

-- name: RetrieveData :one
SELECT url, content FROM data WHERE url=?;
