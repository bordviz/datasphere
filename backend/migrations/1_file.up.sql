CREATE TABLE IF NOT EXISTS file (
    id SERIAL PRIMARY KEY UNIQUE,
    filename VARCHAR(50) NOT NULL,
    size INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    chunks_count INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_filename_file ON file(filename);
CREATE INDEX IF NOT EXISTS idx_name_file ON file(name);