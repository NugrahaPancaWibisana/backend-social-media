-- Insert Accounts
INSERT INTO accounts (id, email, password, lastlogin_at, deleted_at, created_at, updated_at) VALUES
(1, 'john.doe@email.com', '$argon2id$v=19$m=65536,t=2,p=32$mw$OvWJD6JEeUmzAZ4CTkl0CA', '2026-02-07 08:30:00', NULL, '2025-01-15 10:00:00', '2026-02-07 08:30:00'),
(2, 'jane.smith@email.com', '$argon2id$v=19$m=65536,t=2,p=32$mw$OvWJD6JEeUmzAZ4CTkl0CA', '2026-02-06 14:20:00', NULL, '2025-01-20 11:30:00', '2026-02-06 14:20:00'),
(3, 'mike.johnson@email.com', '$argon2id$v=19$m=65536,t=2,p=32$mw$OvWJD6JEeUmzAZ4CTkl0CA', '2026-02-07 07:15:00', NULL, '2025-02-01 09:45:00', '2026-02-07 07:15:00'),
(4, 'sarah.williams@email.com', '$argon2id$v=19$m=65536,t=2,p=32$mw$OvWJD6JEeUmzAZ4CTkl0CA', '2026-02-05 16:40:00', NULL, '2025-02-10 14:20:00', '2026-02-05 16:40:00'),
(5, 'david.brown@email.com', '$argon2id$v=19$m=65536,t=2,p=32$mw$OvWJD6JEeUmzAZ4CTkl0CA', '2026-02-07 09:50:00', NULL, '2025-02-15 08:10:00', '2026-02-07 09:50:00'),
(6, 'emma.davis@email.com', '$argon2id$v=19$m=65536,t=2,p=32$mw$OvWJD6JEeUmzAZ4CTkl0CA', '2026-02-06 12:30:00', NULL, '2025-03-01 10:30:00', '2026-02-06 12:30:00'),
(7, 'alex.wilson@email.com', '$argon2id$v=19$m=65536,t=2,p=32$mw$OvWJD6JEeUmzAZ4CTkl0CA', '2026-02-04 18:20:00', NULL, '2025-03-10 15:45:00', '2026-02-04 18:20:00'),
(8, 'olivia.taylor@email.com', '$argon2id$v=19$m=65536,t=2,p=32$mw$OvWJD6JEeUmzAZ4CTkl0CA', '2026-02-07 06:45:00', NULL, '2025-03-20 09:00:00', '2026-02-07 06:45:00'),
(9, 'daniel.anderson@email.com', '$argon2id$v=19$m=65536,t=2,p=32$mw$OvWJD6JEeUmzAZ4CTkl0CA', NULL, '2025-12-01 10:00:00', '2025-04-01 11:15:00', '2025-12-01 10:00:00'),
(10, 'sophia.martinez@email.com', '$argon2id$v=19$m=65536,t=2,p=32$mw$OvWJD6JEeUmzAZ4CTkl0CA', '2026-02-06 20:10:00', NULL, '2025-04-15 13:30:00', '2026-02-06 20:10:00');

SELECT pg_catalog.setval('public.accounts_id_seq', 10, true);