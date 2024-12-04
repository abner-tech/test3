INSERT INTO users_permissions
VALUES (
    (SELECT id FROM users WHERE email = 'john@example.com'),
    (SELECT id FROM permissions WHERE code = 'reviews:write')
), 
(
    (SELECT id FROM users WHERE email = 'john@example.com'),
    (SELECT id FROM permissions WHERE code = 'books:write')
), 
(
    (SELECT id FROM users WHERE email = 'john@example.com'),
    (SELECT id FROM permissions WHERE code = 'reading_list:write')
),
(
    (SELECT id FROM users WHERE email = 'john@example.com'),
    (SELECT id FROM permissions WHERE code = 'users:write')
);