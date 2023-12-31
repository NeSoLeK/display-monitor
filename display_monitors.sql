CREATE TABLE IF NOT EXISTS Displays (
    ID_Display INTEGER PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    Display_Diagonal REAL NOT NULL,
    Display_Resolution TEXT NOT NULL,
    Display_Type TEXT NOT NULL,
    Display_Gsync BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS Monitors (
    ID_Monitor INTEGER PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    Display_ID INTEGER NOT NULL references Displays(ID_Display),
    Monitor_Gsync_Premium BOOLEAN NOT NULL,
    Monitor_Curved BOOLEAN NOT NULL

);

CREATE TABLE IF NOT EXISTS Users (
    ID_User INTEGER PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    Username_User VARCHAR(20) NOT NULL,
    Password_User TEXT NOT NULL,
    Email_User TEXT NOT NULL,
    Is_Admin_User BOOLEAN
)
