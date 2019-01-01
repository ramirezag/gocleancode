CREATE TABLE files (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    content_type VARCHAR(255),
    created_dt TIMESTAMP NOT NULL -- created date time
);