-- Create the database
CREATE DATABASE airline;

-- Use the created database
USE airline;

-- Create the users table
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- Create the trips table
CREATE TABLE trips (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(20) NOT NULL
);

-- Create the seats table
CREATE TABLE seats (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(10) NOT NULL,
    trip_id INT NOT NULL,
    assigned_to INT,
    FOREIGN KEY (trip_id) REFERENCES trips(id),
    FOREIGN KEY (assigned_to) REFERENCES users(id)
);


