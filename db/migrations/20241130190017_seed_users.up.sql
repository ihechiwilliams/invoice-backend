INSERT INTO users (id, name, email, phone, password, created_at, updated_at)
VALUES
    (uuid_generate_v4(), 'John Doe', 'johndoe@example.com', '+1234567890', crypt('password123', gen_salt('bf')), NOW(), NOW()),
    (uuid_generate_v4(), 'Jane Smith', 'janesmith@example.com', '+1234567891', crypt('password123', gen_salt('bf')), NOW(), NOW()),
    (uuid_generate_v4(), 'Alice Johnson', 'alicejohnson@example.com', '+1234567892', crypt('password123', gen_salt('bf')), NOW(), NOW()),
    (uuid_generate_v4(), 'Bob Williams', 'bobwilliams@example.com', '+1234567893', crypt('password123', gen_salt('bf')), NOW(), NOW()),
    (uuid_generate_v4(), 'Charlie Brown', 'charliebrown@example.com', '+1234567894', crypt('password123', gen_salt('bf')), NOW(), NOW());
