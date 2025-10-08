create table if not exists workouts (
   id         uuid primary key,
   name       text not null,
   exercises  jsonb, -- corresponde al []string de Go
   duration   int,
   user_id    uuid not null,
   route_id   uuid not null,
   created_at timestamp default now(),
   foreign key ( user_id )
      references users ( id ),
   foreign key ( route_id )
      references routes ( id )
);