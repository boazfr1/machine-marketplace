-- Insert users first (they must exist before referencing)
INSERT INTO users (name, email, password) VALUES
('Alice Smith', 'alice@test.com', sha256('password1')),
('Bob Jones', 'bob@test.com', sha256('password2')),
('Carol Wilson', 'carol@test.com', sha256('password3')),
('David Brown', 'david@test.com', sha256('password4')),
('Eva Green', 'eva@test.com', sha256('password5')),
('Frank Miller', 'frank@test.com', sha256('password6')),
('Grace Lee', 'grace@test.com', sha256('password7')),
('Henry Ford', 'henry@test.com', sha256('password8')),
('Iris West', 'iris@test.com', sha256('password9')),
('Jack White', 'jack@test.com', sha256('password10'));

-- Then insert credit cards
INSERT INTO credit_cards (owner_id, number, expiration_date, secret) VALUES
(1, 1234567890, '12/25', 123),
(1, 1234567891, '01/26', 456),
(2, 1234567892, '03/26', 789),
(2, 1234567893, '04/26', 234),
(3, 1234567894, '05/26', 567),
(3, 1234567895, '06/26', 890),
(4, 1234567896, '07/26', 345),
(4, 1234567897, '08/26', 678),
(5, 1234567898, '09/26', 901),
(5, 1234567899, '10/26', 432);

-- Finally insert machines
INSERT INTO machines (name, buyer_id, owner_id, ram, cpu, memory, key, host, ssh_user) VALUES
('Server-1', 1, 2, 16384, 4, 512, 'ssh-rsa AAAA...', '192.168.1.100', 'admin'),
('Server-2', NULL, 2, 32768, 8, 1024, 'ssh-rsa BBBB...', '192.168.1.101', 'root'),
('Server-3', 3, 1, 8192, 2, 256, 'ssh-rsa CCCC...', '192.168.1.102', 'ubuntu'),
('Server-4', 3, 4, 65536, 16, 2048, 'ssh-rsa DDDD...', '192.168.1.103', 'admin'),
('Server-5', NULL, 5, 4096, 1, 128, 'ssh-rsa EEEE...', '192.168.1.104', 'root'),
('Server-6', 4, 6, 16384, 4, 512, 'ssh-rsa FFFF...', '192.168.1.105', 'ubuntu'),
('Server-7', 5, 7, 32768, 8, 1024, 'ssh-rsa GGGG...', '192.168.1.106', 'admin'),
('Server-8', NULL, 8, 8192, 2, 256, 'ssh-rsa HHHH...', '192.168.1.107', 'root'),
('Server-9', 6, 9, 16384, 4, 512, 'ssh-rsa IIII...', '192.168.1.108', 'ubuntu'),
('Server-10', 7, 10, 65536, 16, 2048, 'ssh-rsa JJJJ...', '192.168.1.109', 'admin');
