insert into Users(username, email, password, balance, role, createat)
values ('admin',
        'adminEmail@gmail.com',
        '$2a$10$P6UYTa6WD12GY88y6ZRo3.qLOvSHFgKT7JHaOMl0Qd8Ai2qslpbom',
        100,
        '{"user", "admin"}',
        clock_timestamp()) --password: admin