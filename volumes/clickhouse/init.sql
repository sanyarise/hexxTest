CREATE TABLE IF NOT EXISTS logs
(
  "time" String,
  "level" String,
  "message" String
) Engine = MergeTree()
  ORDER BY tuple();