CREATE TABLE IF NOT EXISTS units_v2 (
	unit_id int8 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE) NOT NULL,
	unit_name text NOT NULL,
	storage_path text DEFAULT ''::text NOT NULL,
	prometheus_server_address text DEFAULT ''::text NULL,
	CONSTRAINT units_v2_pkey PRIMARY KEY (unit_id),
	CONSTRAINT units_v2_unit_name_check CHECK ((unit_name <> ''::text)),
	CONSTRAINT units_v2_unit_name_key UNIQUE (unit_name)
);

-- INSERT INTO units_v2 (unit_id, unit_name, storage_path, prometheus_server_address) OVERRIDING SYSTEM VALUE VALUES (9, 'EDISA', '/srv/sbsautfs01/02_Repositorio', 'http://sbsautlx08:9090');
-- INSERT INTO units_v2 (unit_id, unit_name, storage_path, prometheus_server_address) OVERRIDING SYSTEM VALUE VALUES (8, 'PMXL', '/srv/sbsautfs01/05_Replica_Repositorio_PMXL-1', 'http://smxlautlx02:9090');
-- INSERT INTO units_v2 (unit_id, unit_name, storage_path, prometheus_server_address) OVERRIDING SYSTEM VALUE VALUES (6, 'P71', '/srv/sbsautfs01/05_Replica_Repositorio_P-71', 'https://sp71autlx02:9090');
-- INSERT INTO units_v2 (unit_id, unit_name, storage_path, prometheus_server_address) OVERRIDING SYSTEM VALUE VALUES (1, 'P66', '/srv/sbsautfs01/05_Replica_Repositorio_P-66', 'https://sp66autlx02:9090');
-- INSERT INTO units_v2 (unit_id, unit_name, storage_path, prometheus_server_address) OVERRIDING SYSTEM VALUE VALUES (2, 'P67', '/srv/sbsautfs01/05_Replica_Repositorio_P-67', 'https://sp67autlx02:9090');
-- INSERT INTO units_v2 (unit_id, unit_name, storage_path, prometheus_server_address) OVERRIDING SYSTEM VALUE VALUES (3, 'P68', '/srv/sbsautfs01/05_Replica_Repositorio_P-68', 'https://sp68autlx02:9090');
-- INSERT INTO units_v2 (unit_id, unit_name, storage_path, prometheus_server_address) OVERRIDING SYSTEM VALUE VALUES (4, 'P69', '/srv/sbsautfs01/05_Replica_Repositorio_P-69', 'https://sp69autlx02:9090');
-- INSERT INTO units_v2 (unit_id, unit_name, storage_path, prometheus_server_address) OVERRIDING SYSTEM VALUE VALUES (5, 'P70', '/srv/sbsautfs01/05_Replica_Repositorio_P-70', 'https://sp70autlx02:9090');