CREATE TABLE users (
  id           SERIAL PRIMARY KEY,
  username     TEXT   UNIQUE  NOT NULL,
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


CREATE TABLE messages (

  id            SERIAL PRIMARY KEY, 
  user_id       INT REFERENCES users(id)  ON DELETE CASCADE,
  room_id       INT REFERENCES rooms(id)  ON DELETE CASCADE, 
  content       TEXT NOT NULL , 
  created_at    TIMESTAMP DEFAULT NOW() 

); 


CREATE INDEX  idx_messages_to_room_id   ON  messages(room_id); 

CREATE INDEX  idx_messages_to_user_id   ON  messages(user_Id);  

CREATE INDEX  idx_room_users_to_room_id ON  room_users(room_id);

CREATE INDEX  idx_room_users_to_user_id ON  room_users(user_id);







