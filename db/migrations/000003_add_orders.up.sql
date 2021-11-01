create table ORDERS
(
    ORDERS_NR bigint primary key,
    USERS_ID  bigint not null references USERS (USERS_ID)
);