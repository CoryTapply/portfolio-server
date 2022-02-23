SELECT datname FROM pg_database;

SELECT * FROM information_schema.tables 
WHERE table_schema = 'videoschema';


SELECT * FROM pg_catalog.pg_tables
WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';


CREATE SCHEMA VideoSchema;

CREATE TABLE VideoSchema.tblVideo (
    VideoID    VARCHAR(40) CONSTRAINT PK_tblVideo_VideoID PRIMARY KEY,
    Title      VARCHAR(150) NOT NULL,
    Tags       VARCHAR(150) NOT NULL,
    Game       VARCHAR(100) NOT NULL,
    HasVoice   BOOLEAN NOT NULL,
    ViewCount  INTEGER DEFAULT 0 NOT NULL,
    Duration   NUMERIC
);


-- CREATE TABLE tblVideo (VideoID    VARCHAR(40) CONSTRAINT PK_tblVideo_VideoID PRIMARY KEY, Title      VARCHAR(100) NOT NULL, Tags       VARCHAR(100) NOT NULL, Game       VARCHAR(100),HasVoice   BOOLEAN);

INSERT INTO VideoSchema.tblVideo (VideoID, Title, Tags, Game, HasVoice, Duration) VALUES ('123456', 'Video Title', 'cod,snipe', 'Warzone', false, '18.5')
-- INSERT INTO tblVideo (VideoID, Title, Tags, Game, HasVoice) VALUES (12345, 'Video Title', 'cod,snipe', 'Warzone', false);



SELECT * FROM VideoSchema.tblVideo