CREATE TABLE IF NOT EXISTS work_report_topics (
	work_report_topic_id int8 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE) NOT NULL,
	work_report_id int8 NULL,
	work_report_topic_title text NOT NULL,
	work_report_topic_text text NOT NULL,
	ts tsvector GENERATED ALWAYS AS (setweight(to_tsvector('portuguese'::regconfig, work_report_topic_title), 'A'::"char") || setweight(to_tsvector('portuguese'::regconfig, work_report_topic_text), 'B'::"char")) STORED NULL,
	work_report_topic_text_hash int4 NULL,
	CONSTRAINT work_report_topics_pkey PRIMARY KEY (work_report_topic_id)
);


