CREATE TABLE roles_permissions (
    role_name       TEXT NOT NULL,
    permission_name TEXT NOT NULL,

    CONSTRAINT pk_roles_permissions 
        PRIMARY KEY (role_name, permission_name),

    CONSTRAINT fk_roles_permissions_role_name
        FOREIGN KEY (role_name) 
        REFERENCES roles(name) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,

    CONSTRAINT fk_roles_permissions_permission_name
        FOREIGN KEY (permission_name) 
        REFERENCES permissions(name) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE
);
