-- -----------------------------------------------------
-- Table notifications
-- -----------------------------------------------------
CREATE TABLE notifications (
    id varchar(36) DEFAULT(UUID()) PRIMARY KEY,
    type VARCHAR (255) NOT NULL,
    notifiable_group VARCHAR (255) NOT NULL DEFAULT 'default', -- Group name. Example `user`
    notifiable_id varchar(36) NOT NULL,   -- The id. Example `b9900364-c506-43eb-b1e3-830b5b8370ee`
    data VARCHAR(500) NULL, -- JSON data.
    read_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL
);

-- Add indexes
CREATE INDEX active_notifications ON notifications (id);


