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

-- Finally insert machines
INSERT INTO machines (name, buyer_id, owner_id, ram, cpu, gpu, memory, key, host, ssh_user) VALUES
('Linux-Dev-1', NULL, 2, 8, 2, 0, 256, 'ssh-rsa AAAA...', '192.168.1.100', 'admin'),
('Linux-Pro-2', NULL, 2, 16, 4, 1, 512, 'ssh-rsa BBBB...', '192.168.1.101', 'root'),
('Linux-AI-3', NULL, 1, 32, 8, 2, 1024, 'ssh-rsa CCCC...', '192.168.1.102', 'ubuntu'),
('Linux-Gaming-4', NULL, 4, 64, 16, 4, 2048, 'ssh-rsa DDDD...', '192.168.1.103', 'admin'),
('Linux-Basic-5', NULL, 5, 4, 1, 0, 128, 'ssh-rsa EEEE...', '192.168.1.104', 'root'),
('Linux-Server-6', NULL, 6, 16, 4, 1, 512, 'ssh-rsa FFFF...', '192.168.1.105', 'ubuntu'),
('Linux-ML-7', NULL, 7, 32, 8, 2, 1024, 'ssh-rsa GGGG...', '192.168.1.106', 'admin'),
('Linux-Web-8', NULL, 8, 8, 2, 0, 256, 'ssh-rsa HHHH...', '192.168.1.107', 'root'),
('Linux-Data-9', NULL, 9, 16, 4, 1, 512, 'ssh-rsa IIII...', '192.168.1.108', 'ubuntu'),
('Linux-HPC-10', NULL, 10, 64, 16, 4, 2048, 'ssh-rsa JJJJ...', '192.168.1.109', 'admin');
