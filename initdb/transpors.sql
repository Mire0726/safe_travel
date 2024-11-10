CREATE DATABASE IF NOT EXISTS db;
USE db;

CREATE TABLE IF NOT EXISTS transports (
    id VARCHAR(255) PRIMARY KEY,
    event_id VARCHAR(255) NOT NULL,
    transport_type ENUM('plane', 'train', 'ship', 'bus') NOT NULL,
    memo VARCHAR(255) NOT NULL,
    start_location VARCHAR(255) NOT NULL,
    arrive_location VARCHAR(255) NOT NULL,
    start_at DATETIME NOT NULL,
    arrive_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id)
);