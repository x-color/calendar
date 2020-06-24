CREATE SCHEMA auth;
CREATE SCHEMA calendar;

CREATE TABLE IF NOT EXISTS auth.users (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    password VARCHAR(72) NOT NULL
);

CREATE TABLE IF NOT EXISTS calendar.users (
    id CHAR(36) PRIMARY KEY,
    FOREIGN KEY (id) REFERENCES auth.users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS calendar.calendars (
    id CHAR(36) PRIMARY KEY,
    userid CHAR(36),
    name NAME NOT NULL,
    color VARCHAR(20) NOT NULL,
    FOREIGN KEY (userid) REFERENCES calendar.users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS calendar.calendar_shares (
    userid CHAR(36),
    calendarid CHAR(36),
    PRIMARY KEY(userid, calendarid),
    FOREIGN KEY (userid) REFERENCES calendar.users(id) ON DELETE CASCADE,
    FOREIGN KEY (calendarid) REFERENCES calendar.calendars(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS calendar.plans (
    id CHAR(36) PRIMARY KEY,
    userid CHAR(36),
    calendarid CHAR(36),
    name NAME NOT NULL,
    memo VARCHAR(400),
    color VARCHAR(20) NOT NULL,
    private BOOLEAN NOT NULL,
    isallday BOOLEAN NOT NULL,
    begintime BIGINT NOT NULL,
    endtime BIGINT NOT NULL,
    FOREIGN KEY (userid) REFERENCES calendar.users(id) ON DELETE CASCADE,
    FOREIGN KEY (calendarid) REFERENCES calendar.calendars(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS calendar.plan_shares (
    calendarid CHAR(36),
    planid CHAR(36),
    PRIMARY KEY(calendarid, planid),
    FOREIGN KEY (calendarid) REFERENCES calendar.calendars(id) ON DELETE CASCADE,
    FOREIGN KEY (planid) REFERENCES calendar.plans(id) ON DELETE CASCADE
);
