INSERT INTO products (name, type, flavor, size, price, stock)
VALUES
('Keripik Pangsit', 'Snack', 'Jagung Bakar', 'Small', 10000, 10),
('Keripik Pangsit', 'Snack', 'Rumput Laut', 'Small', 10000, 10),
('Keripik Pangsit', 'Snack', 'Original', 'Small', 10000, 10),
('Keripik Pangsit', 'Snack', 'Jagung Manis', 'Small', 10000, 10),
('Keripik Pangsit', 'Snack', 'Keju Asin', 'Small', 10000, 10),
('Keripik Pangsit', 'Snack', 'Keju Manis', 'Small', 10000, 10),
('Keripik Pangsit', 'Snack', 'Pedas', 'Small', 10000, 10),

('Keripik Pangsit', 'Snack', 'Jagung Bakar', 'Medium', 25000, 8),
('Keripik Pangsit', 'Snack', 'Rumput Laut', 'Medium', 25000, 8),
('Keripik Pangsit', 'Snack', 'Original', 'Medium', 25000, 8),
('Keripik Pangsit', 'Snack', 'Jagung Manis', 'Medium', 25000, 8),
('Keripik Pangsit', 'Snack', 'Keju Asin', 'Medium', 25000, 8),
('Keripik Pangsit', 'Snack', 'Keju Manis', 'Medium', 25000, 8),
('Keripik Pangsit', 'Snack', 'Pedas', 'Medium', 25000, 8),

('Keripik Pangsit', 'Snack', 'Jagung Bakar', 'Large', 35000, 5),
('Keripik Pangsit', 'Snack', 'Rumput Laut', 'Large', 35000, 5),
('Keripik Pangsit', 'Snack', 'Original', 'Large', 35000, 5),
('Keripik Pangsit', 'Snack', 'Jagung Manis', 'Large', 35000, 5),
('Keripik Pangsit', 'Snack', 'Keju Asin', 'Large', 35000, 5),
('Keripik Pangsit', 'Snack', 'Keju Manis', 'Large', 35000, 5),
('Keripik Pangsit', 'Snack', 'Pedas', 'Large', 35000, 5);

SELECT COUNT(*) FROM products;

SELECT id, size, flavor, price, stock
FROM products
ORDER BY size, flavor;

SELECT id, name, size, flavor, price, stock
FROM products
ORDER BY id
LIMIT 5;

SELECT stock FROM products WHERE id = 1;

SELECT id, customer_id, total_price, created_at
FROM transactions
ORDER BY created_at DESC
LIMIT 1;


SELECT product_id, quantity, price
FROM transaction_items
WHERE transaction_id = 'ed6ce6a9-a5ea-4e87-8507-657cc65c9b64';
