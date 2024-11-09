create table if not exists FileMD
(
    Id        serial primary key,

    Cid       varchar(255) not null,
    ownerAddr varchar(255) not null,
    fileName  varchar(255) not null,

    fileSize  smallint     not null check (fileSize >= 0),

    isDelete  boolean      not null default true, --true-delete, false - no
    isSecured boolean      not null default false  --есть пароль или нет
);

create table if not exists Blockchain
(
    Index         serial       not null,
    FileMD_id     serial       not null,

    Hash          varchar(255) not null,
    PrevBlockHash varchar(255) not null,
    Timestamp     timestamp    not null default clock_timestamp(),

    foreign key (FileMD_id) references FileMD (Id)
);



create index if not exists BcIndex_idx on Blockchain using hash (Index);
create index if not exists Bc_FileMd_id_idx on Blockchain using hash (FileMD_id);


-- drop table Blockchain;
-- drop table FileMD;