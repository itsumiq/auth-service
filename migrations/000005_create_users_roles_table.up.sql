CREATE TABLE users_roles (
    user_id     INTEGER NOT NULL,
    role_name   TEXT NOT NULL,

    CONSTRAINT pk_users_roles 
        PRIMARY KEY (user_id, role_name),

    CONSTRAINT fk_users_roles_user_id
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,

    CONSTRAINT fk_users_roles_role_name
        FOREIGN KEY (role_name) 
        REFERENCES roles(name) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
);
