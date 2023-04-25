create database c4_transacciones 
use c4_transacciones 

CREATE TABLE Users (
  user_id INT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(50) NOT NULL,
  email VARCHAR(100) NOT NULL,
  password VARCHAR(50) NOT NULL
);

CREATE TABLE Wallets (
  wallet_id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  balance FLOAT NOT NULL DEFAULT 0,
  FOREIGN KEY (user_id) REFERENCES Users(user_id)
);

CREATE TABLE Providers (
  provider_id INT PRIMARY KEY AUTO_INCREMENT,
  provider_name VARCHAR(50) NOT NULL,
  provider_contact VARCHAR(100) NOT NULL
);

CREATE TABLE Requests (
  request_id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  provider_id INT NOT NULL,
  request_details VARCHAR(200) NOT NULL,
  request_date DATETIME NOT NULL,
  is_used BOOLEAN NOT NULL DEFAULT FALSE,
  FOREIGN KEY (user_id) REFERENCES Users(user_id),
  FOREIGN KEY (provider_id) REFERENCES Providers(provider_id)
);
CREATE TABLE Transactions (
  transaction_id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  provider_id INT NOT NULL,
  request_id INT NOT NULL,
  transaction_type VARCHAR(20) NOT NULL,
  transaction_date DATETIME NOT NULL,
  amount FLOAT NOT NULL,
  transaction_description VARCHAR(200) NOT NULL,
  FOREIGN KEY (user_id) REFERENCES Users(user_id),
  FOREIGN KEY (provider_id) REFERENCES Providers(provider_id),
  FOREIGN KEY (request_id) REFERENCES Requests(request_id)
);

CREATE TABLE ServiceDelivery (
  service_delivery_id INT PRIMARY KEY AUTO_INCREMENT,
  request_id INT NOT NULL,
  provider_id INT NOT NULL,
  delivery_date DATETIME NOT NULL,
  delivery_status VARCHAR(20) NOT NULL,
  FOREIGN KEY (request_id) REFERENCES Requests(request_id),
  FOREIGN KEY (provider_id) REFERENCES Providers(provider_id)
);


