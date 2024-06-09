-- name: CreateCeoUser :one
insert into users (id, create_at,update_at,name,user_name,password,role)
values ($1,$2,$3,$4,$5,$6,$7)
returning *;

-- name: LoginUser :one
select id,password from users where user_name = $1;

-- name: GenApiKey :one
update users set api_key = $1 where id = $2
returning *;

-- name: GetUserByApi :one
select * from users where api_key =$1;

-- name: LogoutUser :exec
update users set api_key = null where id = $1 and api_key = $2
returning *;