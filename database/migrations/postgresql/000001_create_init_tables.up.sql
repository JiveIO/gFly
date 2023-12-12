-- Add UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
-- For more information, please visit:
-- https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
SET TIMEZONE="GMT-0";


-- -----------------------------------------------------
-- Table users
-- -----------------------------------------------------
CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    email VARCHAR (255) NOT NULL UNIQUE,
    password_hash VARCHAR (255) NOT NULL,
    fullname VARCHAR (255) NULL,
    phone VARCHAR(20) NULL,
    token VARCHAR (100) NULL,
    user_status INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP NULL,
    verified_at TIMESTAMP NULL,
    blocked_at TIMESTAMP NULL,
    deleted_at TIMESTAMP NULL,
    last_access_at TIMESTAMP NULL
);

-- Add indexes
CREATE INDEX active_users ON users (id) WHERE user_status = 1;
CREATE UNIQUE INDEX email_users ON users (email ASC) WHERE user_status = 1;


-- -----------------------------------------------------
-- Table roles
-- -----------------------------------------------------
CREATE TABLE roles (
  id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
  updated_at TIMESTAMP NULL
);

-- Add indexes
CREATE INDEX active_roles ON roles (id);

-- Insert data
INSERT INTO roles VALUES(uuid_in(md5(random()::text || random()::text)::cstring), 'Guest', 'guest');
INSERT INTO roles VALUES(uuid_in(md5(random()::text || random()::text)::cstring), 'User', 'user');
INSERT INTO roles VALUES(uuid_in(md5(random()::text || random()::text)::cstring), 'Moderator', 'moderator');
INSERT INTO roles VALUES(uuid_in(md5(random()::text || random()::text)::cstring), 'Admin', 'admin');

-- -----------------------------------------------------
-- Table user_roles
-- -----------------------------------------------------
CREATE TABLE user_roles (
  id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
  role_id UUID NOT NULL,
  user_id UUID NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW (),
  CONSTRAINT fk_user_roles_roles
    FOREIGN KEY (role_id)
    REFERENCES roles (id)
    ON DELETE CASCADE,
  CONSTRAINT fk_user_roles_users
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON DELETE CASCADE
);





