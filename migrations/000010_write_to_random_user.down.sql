DELETE FROM users_permissions
WHERE user_id = (SELECT id FROM users WHERE email = 'jochn@example.com')
AND permission_id = (SELECT id FROM permissions WHERE code = 'reviews:write');

DELETE FROM users_permissions
WHERE user_id = (SELECT id FROM users WHERE email = 'jochn@example.com')
AND permission_id = (SELECT id FROM permissions WHERE code = 'books:write');

DELETE FROM users_permissions
WHERE user_id = (SELECT id FROM users WHERE email = 'jochn@example.com')
AND permission_id = (SELECT id FROM permissions WHERE code = 'reading_list:write');

DELETE FROM users_permissions
WHERE user_id = (SELECT id FROM users WHERE email = 'jochn@example.com')
AND permission_id = (SELECT id FROM permissions WHERE code = 'users:write');