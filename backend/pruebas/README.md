# Pruebas de la mesa de ayuda (Práctica 7)

`escenario_soporte.py` recorre la matriz TC01–TC08 y las reglas de roles
(P7-R01…R29) contra el backend real conectado a un Postgres local. No usa
Supabase: `fake_storage.py` simula los dos endpoints de Storage que firman URLs.

## Cómo correrlo

```bash
# 1. Postgres local con el esquema del blog + 005_soporte.sql
createdb bw
psql -d bw -c "create schema storage; create table storage.buckets(id text primary key, name text, public bool, file_size_limit bigint, allowed_mime_types text[]);"
psql -d bw -f sql/001_initial_schema.sql -f sql/002_add_category.sql
psql -d bw -c "ALTER TABLE users ADD COLUMN IF NOT EXISTS active BOOLEAN NOT NULL DEFAULT true"
psql -d bw -f sql/005_soporte.sql
psql -d bw -c "insert into users(supabase_uid,email,username,role) values
  ('11111111-1111-1111-1111-111111111111','admin@bw.test','abraham','admin'),
  ('22222222-2222-2222-2222-222222222222','lector@bw.test','lector','user')"

# 2. Storage simulado y backend
python3 pruebas/fake_storage.py 18099 &
SUPABASE_URL=http://127.0.0.1:18099 SUPABASE_SERVICE_ROLE_KEY=service-key \
JWT_SECRET=prueba PORT=18080 DATABASE_URL="host=localhost dbname=bw sslmode=disable" \
go run ./cmd/api &

# 3. Escenario completo y, tras reiniciar el backend, la persistencia (TC08)
python3 pruebas/escenario_soporte.py http://127.0.0.1:18080/api
python3 pruebas/escenario_soporte.py http://127.0.0.1:18080/api persistencia
```

El script guarda el token del admin en `/var/tmp/bwpg/e2e_estado.json` para
la fase de persistencia; cambia la ruta `ESTADO` si tu sistema no tiene esa
carpeta.

Resultado obtenido el 22/10/2026: **37/37** casos del escenario y **1/1** de
persistencia.

# Práctica 10 — reglas RLS de `incidencias`

`rls_incidencias.sh` ejecuta SEG-03…SEG-10 con los roles de Supabase
(`authenticated` con el JWT del usuario, y `anon` sin sesión).

```bash
createdb bw10
psql -d bw10 <<'SQL'
create role anon nologin; create role authenticated nologin;
create schema auth;
create table auth.users (id uuid primary key, email text);
create function auth.uid() returns uuid language sql stable as $$
  select coalesce(nullif(current_setting('request.jwt.claim.sub', true), ''),
                  (nullif(current_setting('request.jwt.claims', true), '')::jsonb ->> 'sub'))::uuid $$;
grant usage on schema public, auth to anon, authenticated;
alter default privileges in schema public grant all on tables to anon, authenticated;
insert into auth.users values ('aaaaaaaa-0000-0000-0000-00000000000a','a@bw.test'),
                              ('bbbbbbbb-0000-0000-0000-00000000000b','b@bw.test');
SQL
psql -d bw10 -f sql/001_initial_schema.sql -f sql/006_incidencias_rls.sql
PSQL="psql -d bw10" pruebas/rls_incidencias.sh
```

Resultado obtenido el 29/10/2026: **9 pasan, 0 fallan**. Control negativo:
con RLS desactivado y `anon` con permisos, fallan 7 de 9 (A lee, modifica y
borra lo de B; sin sesión se lee todo). Al aplicar de nuevo 006, 9/9.
