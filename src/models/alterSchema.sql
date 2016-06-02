ALTER TABLE vsfootball.game ADD COLUMN `Whatsontheline` VARCHAR(255) NOT NULL DEFAULT '0';
CREATE TABLE IF NOT EXISTS vsfootball.preset (
	Id BIGINT NOT NULL UNIQUE AUTO_INCREMENT, 
	Value CHAR(32) NOT NULL,
	PRIMARY KEY(Id));
INSERT INTO vsfootball.preset(ID,Value) VALUES (1,'2');
INSERT INTO vsfootball.preset(ID,Value) VALUES (2,'4');
INSERT INTO vsfootball.preset(ID,Value) VALUES (3,'5');
INSERT INTO vsfootball.preset(ID,Value) VALUES (4,'6');
