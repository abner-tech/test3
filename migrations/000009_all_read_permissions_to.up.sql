INSERT INTO users_permissions
SELECT id, (SELECT id FROM permissions WHERE CODE = 'reviews:read')
FROM users;

INSERT INTO users_permissions
SELECT id, (SELECT id FROM permissions WHERE CODE = 'books:read')
FROM users;

INSERT INTO users_permissions
SELECT id, (SELECT id FROM permissions WHERE CODE = 'reading_list:read')
FROM users;

INSERT INTO users_permissions
SELECT id, (SELECT id FROM permissions WHERE CODE = 'users:read')
FROM users;