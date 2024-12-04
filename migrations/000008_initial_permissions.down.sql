DELETE FROM permissions 
WHERE code
IN
    ('reviews:read', 'reviews:write');