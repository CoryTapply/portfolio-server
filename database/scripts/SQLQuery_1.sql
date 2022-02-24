SELECT datname FROM pg_database;

SELECT * FROM information_schema.tables 
WHERE table_schema = 'videoschema';


SELECT * FROM pg_catalog.pg_tables
WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema';


CREATE SCHEMA VideoSchema;

CREATE TABLE VideoSchema.tblVideo (
    VideoID    SERIAL CONSTRAINT PK_tblVideo_VideoID PRIMARY KEY,
    UrlSlug    VARCHAR(40) NOT NULL,
    Title      VARCHAR(100) NOT NULL,
    Tags       VARCHAR(100) NOT NULL,
    Game       VARCHAR(100) NOT NULL,
    HasVoice   BOOLEAN NOT NULL,
    ViewCount  INTEGER NOT NULL DEFAULT 0,
    Duration   NUMERIC NOT NULL,
    DateCreated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    DateUpdated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    DateDeleted TIMESTAMP WITH TIME ZONE
);




-- CREATE TABLE tblVideo (VideoID    VARCHAR(40) CONSTRAINT PK_tblVideo_VideoID PRIMARY KEY, Title      VARCHAR(100) NOT NULL, Tags       VARCHAR(100) NOT NULL, Game       VARCHAR(100),HasVoice   BOOLEAN);

INSERT INTO VideoSchema.tblVideo (UrlSlug, Title, Tags, Game, HasVoice, Duration) VALUES ('asdasdasd', 'Video Title', 'cod,snipe', 'Warzone', false, 90)
-- INSERT INTO tblVideo (VideoID, Title, Tags, Game, HasVoice) VALUES (12345, 'Video Title', 'cod,snipe', 'Warzone', false);



SELECT * FROM VideoSchema.tblVideo


DELETE FROM VideoSchema.tblVideo

DROP TABLE VideoSchema.tblVideo
