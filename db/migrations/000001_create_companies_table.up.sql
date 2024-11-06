CREATE TABLE companies (
                           id SERIAL PRIMARY KEY,
                           name VARCHAR(255) NOT NULL,
                           description TEXT,
                           employees INT,
                           registered BOOLEAN DEFAULT FALSE,
                           type VARCHAR(50),
                           created_at TIMESTAMPTZ DEFAULT NOW(),
                           updated_at TIMESTAMPTZ DEFAULT NOW()
);
