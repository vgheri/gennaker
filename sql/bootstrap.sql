BEGIN;

CREATE TABLE IF NOT EXISTS deployment (id SERIAL PRIMARY KEY, name TEXT UNIQUE NOT NULL, chart TEXT NOT NULL, repository_url TEXT NOT NULL, creation_date TIMESTAMP WITH TIME ZONE DEFAULT NOW(), last_update TIMESTAMP WITH TIME ZONE DEFAULT NOW());
CREATE TABLE IF NOT EXISTS pipeline_step (id SERIAL PRIMARY KEY, step_number INT NOT NULL, parent_step_number int, deployment_id INT NOT NULL, target_namespace TEXT NOT NULL, auto_deploy BOOLEAN DEFAULT FALSE);
CREATE TABLE IF NOT EXISTS release (id SERIAL PRIMARY KEY, deployment_id INT NOT NULL, timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(), namespace TEXT NOT NULL, image TEXT NOT NULL, tag TEXT NOT NULL);

CREATE INDEX ON pipeline_step (deployment_id);
CREATE INDEX ON pipeline_step (id, parent_step_number);
ALTER TABLE pipeline_step ADD CONSTRAINT FK_PIPELINE_STEP_DEPLOYMENT_ID FOREIGN KEY (deployment_id) REFERENCES deployment (id);
ALTER TABLE pipeline_step ADD CONSTRAINT PIPELINE_STEP_UNIQUE_STEP_NUMBER_DEPLOYMENT_ID UNIQUE (step_number, deployment_id);
ALTER TABLE pipeline_step ADD CONSTRAINT FK_PIPELINE_STEP_PARENT_STEP_NUMBER FOREIGN KEY (parent_step_number, deployment_id) REFERENCES pipeline_step (step_number, deployment_id);
CREATE INDEX on release (deployment_id);
ALTER TABLE release ADD CONSTRAINT FK_RELEASE_DEPLOYMENT_ID FOREIGN KEY (deployment_id) REFERENCES deployment (id);

COMMIT;
