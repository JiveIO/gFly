-- Add UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
-- For more information, please visit:
-- https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
SET TIMEZONE="GMT-0";


-- -----------------------------------------------------
-- Table notifications
-- -----------------------------------------------------
CREATE TABLE notifications (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    type VARCHAR (255) NOT NULL,
    notifiable_group VARCHAR (255) NOT NULL DEFAULT 'default', -- Group name. Example `user`
    notifiable_id UUID NOT NULL,   -- The id. Example `b9900364-c506-43eb-b1e3-830b5b8370ee`
    data VARCHAR(500) NULL, -- JSON data.
    read_at TIMESTAMP NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

-- Add indexes
CREATE INDEX active_notifications ON notifications (id) WHERE type IS NOT NULL;


