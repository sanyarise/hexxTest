	CREATE TABLE IF NOT EXISTS users
    (
		id varchar UNIQUE NOT NULL,
		created_at timestamptz NOT NULL,
		deleted_at timestamptz,
		name varchar NOT NULL
    );