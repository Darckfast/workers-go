-- Migration number: 0001 	 2025-08-19T17:48:33.037Z
CREATE TABLE IF NOT EXISTS "Testing" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "data" TEXT,
  "created_at" INTEGER NOT NULL DEFAULT ( unixepoch('subsec') * 1000 ),
  "updated_at" INTEGER NOT NULL DEFAULT ( unixepoch('subsec') * 1000 )
);

