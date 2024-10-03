CREATE TABLE IF NOT EXISTS work_reports (
	work_report_id int8 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE) NOT NULL,
	unit_id int8 NULL,
	work_report_docname text NOT NULL,
	work_report_from date NOT NULL,
	work_report_to date NOT NULL,
	work_report_text text NOT NULL,
	work_report_data bytea NULL,
	CONSTRAINT work_reports_pkey PRIMARY KEY (work_report_id),
	CONSTRAINT work_reports_work_report_docname_key UNIQUE (work_report_docname)
);