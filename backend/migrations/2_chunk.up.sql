CREATE TABLE IF NOT EXISTS chunk (
    id SERIAL PRIMARY KEY UNIQUE,
    file_id VARCHAR(100) NOT NULL,
    chunk_number INTEGER NOT NULL,
    file INTEGER NOT NULL REFERENCES file(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_chunk_number_chunk ON chunk(chunk_number);
CREATE INDEX IF NOT EXISTS idx_file_chunk ON chunk(file);