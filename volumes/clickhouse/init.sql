CREATE TABLE logs
(
  "time" String,
  "level" String,
  "message" String
) Engine = MergeTree()
  ORDER BY tuple();