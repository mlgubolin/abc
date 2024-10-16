CREATE TABLE IF NOT EXISTS work_report_topics
(
    id             BIGSERIAL PRIMARY KEY,
    work_report_id int8 NOT NULL REFERENCES work_reports (id),
    title          text NOT NULL,
    content        text NOT NULL,
    ts             tsvector GENERATED ALWAYS AS (setweight(
                                                         to_tsvector('portuguese'::regconfig, work_report_topic_title),
                                                         'A'::"char") ||
                                                 setweight(to_tsvector('portuguese'::regconfig, work_report_topic_text),
                                                           'B'::"char")) STORED
);

CREATE INDEX ON work_report_topics (work_report_id);
-- CREATE INDEX ON work_report_topics (title)
-- CREATE INDEX ON work_report_topics (ts)