CREATE TYPE image_status AS ENUM ('uploading','ready','processing','deleted');

CREATE TYPE image_format AS ENUM ('jpg', 'jpeg', 'png', 'gif', 'webp');

CREATE TABLE IF NOT EXISTS images (
    id SERIAL PRIMARY KEY,
    owner_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    original_name VARCHAR(255) NOT NULL,
    format image_format NOT NULL,
    status image_status NOT NULL DEFAULT 'uploading',
    storage_path VARCHAR(255) NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT FALSE,
    width INT NOT NULL DEFAULT 0,
    height INT NOT NULL DEFAULT 0,
    file_size INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_images_owner_id
    ON images(owner_id, created_at DESC)
    WHERE status != 'deleted';
