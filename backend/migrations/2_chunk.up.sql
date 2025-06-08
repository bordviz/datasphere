CREATE TABLE IF NOT EXISTS chunk (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    file_id VARCHAR(100) NOT NULL,
    chunk_number INTEGER NOT NULL,
    file_key INTEGER NOT NULL,
    FOREIGN KEY (file_key) REFERENCES file(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_chunk_number_chunk ON chunk(chunk_number);
CREATE INDEX IF NOT EXISTS idx_file_key_chunk ON chunk(file_key);