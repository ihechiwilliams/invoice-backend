

CREATE TABLE users (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   name VARCHAR(255) NOT NULL,
   email VARCHAR(255) UNIQUE NOT NULL,
   password_hash VARCHAR(255) NOT NULL,
   role VARCHAR(50) DEFAULT 'user' NOT NULL, -- Possible roles: admin, user, etc.
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
   deleted_at TIMESTAMP NULL
);
