CREATE TABLE IF NOT EXISTS work_reports
(
    id          BIGSERIAL PRIMARY KEY,
    unit_id     int8        NULL,
    data_from   date        NOT NULL,
    data_to     date        NOT NULL,
    report_name text UNIQUE NOT NULL,
    content     text        NOT NULL,
    file_data   bytea       NULL
);

CREATE INDEX ON work_reports (unit_id);
-- CREATE INDEX ON work_reports (from, to)

