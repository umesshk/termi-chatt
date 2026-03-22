CREATE TABLE users (
  id           SERIAL PRIMARY KEY
  email        TEXT UNIQUE NOT NULL 
  username     TEXT NOT NULL ,
  created_at   TIMESTAMP DEFAULT NOW()
);


CREATE TABLE rooms (
  
  id            SERIAL PRIMARY KEY ,
  created_at    TIMESTAMP DEFAULT NOW()

);

CREATE TABLE room_users (

  user_id       INT REFERENCES users(id) ON DELETE CASCADE,
  room_id       INT REFERENCES rooms(id) ON DELETE CASCADE,
  joined_at     TIMESTAMP DEFAULT NOW(),

  PRIMARY KEY(user_id,room_id)
);


CREATE TABLE message (

  id            SERIAL PRIMARY KEY, 
  user_id       INT REFERENCES users(id)  ON DELETE CASCADE,
  room_id       INT REFERENCES rooms(id)  ON DELETE CASCADE, 
  content       TEXT NOT NULL , 
  created_at    TIMESTAMP DEFAULT NOW(), 

); 
