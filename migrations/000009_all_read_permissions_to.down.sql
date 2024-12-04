DELETE FROM users_permissions
WHERE permission_id = (SELECT id FROM permissions WHERE CODE = 'reviews:read');

DELETE FROM users_permissions
WHERE permission_id = (SELECT id FROM permissions WHERE CODE = 'books:read');

DELETE FROM users_permissions
WHERE permission_id = (SELECT id FROM permissions WHERE CODE = 'reading_list:read');

DELETE FROM users_permissions
WHERE permission_id = (SELECT id FROM permissions WHERE CODE = 'users:read');
