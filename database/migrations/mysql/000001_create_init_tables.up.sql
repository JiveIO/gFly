-- -----------------------------------------------------
-- Table users
-- -----------------------------------------------------
CREATE TABLE users (
                       id varchar(36) DEFAULT(UUID()) PRIMARY KEY,
                       email VARCHAR (255) NOT NULL UNIQUE,
                       password_hash VARCHAR (255) NOT NULL,
                       fullname VARCHAR (255) NULL,
                       phone VARCHAR(20) NULL,
                       token VARCHAR (100) NULL,
                       user_status INT NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NULL,
                       verified_at TIMESTAMP NULL,
                       blocked_at TIMESTAMP NULL,
                       deleted_at TIMESTAMP NULL,
                       last_access_at TIMESTAMP NULL
);

-- Add indexes
CREATE INDEX active_users ON users (id);
CREATE UNIQUE INDEX email_users ON users (email ASC);


-- -----------------------------------------------------
-- Table roles
-- -----------------------------------------------------
CREATE TABLE roles (
                       id varchar(36) DEFAULT(UUID()) PRIMARY KEY,
                       name VARCHAR(100) NOT NULL,
                       slug VARCHAR(100) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP NULL
);

-- Add indexes
CREATE INDEX active_roles ON roles (id);

-- Insert data
INSERT INTO roles (id, name, slug) VALUES(UUID(), 'Guest', 'guest');
INSERT INTO roles (id, name, slug) VALUES(UUID(), 'User', 'user');
INSERT INTO roles (id, name, slug) VALUES(UUID(), 'Moderator', 'moderator');
INSERT INTO roles (id, name, slug) VALUES(UUID(), 'Admin', 'admin');

-- -----------------------------------------------------
-- Table user_roles
-- -----------------------------------------------------
CREATE TABLE user_roles (
                            id varchar(36) DEFAULT(UUID()) PRIMARY KEY,
                            role_id varchar(36) NOT NULL,
                            user_id varchar(36) NOT NULL,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            CONSTRAINT fk_user_roles_roles
                                FOREIGN KEY (role_id)
                                    REFERENCES roles (id)
                                    ON DELETE CASCADE,
                            CONSTRAINT fk_user_roles_users
                                FOREIGN KEY (user_id)
                                    REFERENCES users (id)
                                    ON DELETE CASCADE
);

