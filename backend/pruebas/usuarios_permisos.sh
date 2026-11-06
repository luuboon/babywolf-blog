#!/usr/bin/env bash
# Práctica 10 — permisos de public.users, audit_log y funciones de 004
# (reglas de 008_permisos_usuarios.sql), vistos desde la base.
#
# Mismo método que rls_incidencias.sh: cada operación corre con el rol de
# Supabase ("authenticated" con el JWT del usuario, o "anon").
#
# Uso: PSQL="psql -d bwapp" ./usuarios_permisos.sh
# (base con 001, 004, las reglas de users del setup y los usuarios A, B y ADMIN)

set -u
PSQL=${PSQL:-psql}
A=aaaaaaaa-0000-0000-0000-00000000000a
B=bbbbbbbb-0000-0000-0000-00000000000b
ADM=cccccccc-0000-0000-0000-00000000000c
pasan=0; fallan=0

como() {
  local uid=$1 sql=$2 rol=authenticated claims
  if [ "$uid" = anon ]; then rol=anon; claims='{"role":"anon"}'; else claims="{\"sub\":\"$uid\",\"role\":\"authenticated\"}"; fi
  $PSQL -X -q -A -t -v ON_ERROR_STOP=1 2>&1 <<SQL
begin;
set local role $rol;
select set_config('request.jwt.claims', '$claims', true) \g /dev/null
$sql;
commit;
SQL
}
admin_sql() { $PSQL -X -q -A -t -c "$1"; }
rechazado() { echo "$1" | grep -qiE "permission denied|42501|no autorizado|row-level security"; }

caso() {
  if [ "$5" -eq 0 ]; then r=PASA; pasan=$((pasan+1)); else r=FALLA; fallan=$((fallan+1)); fi
  printf '%s | %s | %s | %s | %s\n' "$1" "$2" "$3" "$(echo "$4" | tr '\n' ' ' | cut -c1-90)" "$r"
}

admin_sql "update users set role='user', active=true, failed_login_attempts=0, locked_until=null, must_change_password=false where id in ('$A','$B');
           update users set role='admin' where id='$ADM'; delete from audit_log" >/dev/null
EA=$(admin_sql "select email from users where id='$A'")

# ── Escalada de privilegios ─────────────────────────────────
out=$(como $A "update users set role='admin' where id='$A' returning role")
caso U-01 "A se pone role=admin a sí mismo" "rechazado" "$out" \
  "$(rechazado "$out" && [ "$(admin_sql "select role from users where id='$A'")" = user ]; echo $?)"

out=$(como $A "update users set active=false, failed_login_attempts=0, locked_until=null where id='$A'")
caso U-02 "A cambia active / contador de bloqueo" "rechazado" "$out" "$(rechazado "$out"; echo $?)"

out=$(como $A "update users set email='otro@bw.test' where id='$A'")
caso U-03 "A cambia su correo en users (no en Auth)" "rechazado" "$out" "$(rechazado "$out"; echo $?)"

out=$(como $A "update users set username='lector_a2', avatar_url='https://x/a.png' where id='$A' returning username")
caso U-04 "A edita su nombre y avatar (perfil de la web)" "lector_a2" "$out" "$([ "$out" = lector_a2 ]; echo $?)"
admin_sql "update users set username='lector_a', avatar_url=null where id='$A'" >/dev/null

out=$(como $A "update users set username='pirata' where id='$B' returning id")
caso U-05 "A edita el perfil de B" "0 filas" "${out:-0 filas}" "$([ -z "$out" ]; echo $?)"

out=$(como $A "insert into users(id, supabase_uid, email, username, role) values (gen_random_uuid(), gen_random_uuid(), 'x@bw.test', 'colado', 'admin')")
caso U-06 "A inserta un usuario admin a mano" "rechazado" "$out" "$(rechazado "$out"; echo $?)"

out=$(admin_sql "update users set role='admin' where id='$B' returning role"); admin_sql "update users set role='user' where id='$B'" >/dev/null
caso U-07 "El backend (postgres) cambia un rol" "admin" "$out" "$([ "$out" = admin ]; echo $?)"

# ── Lectura ─────────────────────────────────────────────────
out=$(como anon "select role, failed_login_attempts, locked_until from users limit 1")
caso U-08 "Sin sesión: leer rol y bloqueos" "rechazado" "$out" "$(rechazado "$out"; echo $?)"

out=$(como anon "select count(*) from users where username is not null and email is not null and avatar_url is null or true")
caso U-09 "Sin sesión: id, nombre, correo, avatar (autor de un post)" "3" "$out" "$([ "$out" = 3 ]; echo $?)"

out=$(como $A "select role||' '||active from users where id='$A'")
caso U-10 "A lee su propio perfil completo" "user true" "$out" "$([ "$out" = "user true" ]; echo $?)"

out=$(como $ADM "select count(*) from users where role is not null")
caso U-11 "El admin lista todos los usuarios" "3" "$out" "$([ "$out" = 3 ]; echo $?)"

# ── Funciones de bloqueo por intentos (004) ─────────────────
out=$(como anon "select check_account_lock('$EA')->>'locked'")
caso U-12 "Sin sesión: consultar si una cuenta está bloqueada" "false" "$out" "$([ "$out" = false ]; echo $?)"

for _ in 1 2 3 4 5; do como anon "select record_failed_login('$EA')" >/dev/null; done
out=$(como anon "select check_account_lock('$EA')->>'locked'")
caso U-13 "Sin sesión: 5 intentos fallidos bloquean la cuenta" "true" "$out" "$([ "$out" = true ]; echo $?)"

out=$(como anon "select record_successful_login('$A')")
caso U-14 "Sin sesión: desbloquear la cuenta de A" "rechazado" "$out" \
  "$(rechazado "$out" && [ "$(admin_sql "select failed_login_attempts from users where id='$A'")" = 5 ]; echo $?)"

out=$(como $B "select record_successful_login('$A')")
caso U-15 "B desbloquea la cuenta de A" "rechazado" "$out" "$(rechazado "$out"; echo $?)"

out=$(como $A "select record_successful_login('$A')")
caso U-16 "A entra y su contador vuelve a 0" "0" "$(admin_sql "select failed_login_attempts from users where id='$A'")" \
  "$([ "$(admin_sql "select failed_login_attempts from users where id='$A'")" = 0 ]; echo $?)"

out=$(como $B "select record_logout('$A')")
caso U-17 "B registra un cierre de sesión a nombre de A" "rechazado" "$out" "$(rechazado "$out"; echo $?)"

admin_sql "update users set must_change_password=true where id='$A'" >/dev/null
out=$(como $B "select record_password_change('$A')")
caso U-18 "B quita el 'debe cambiar contraseña' de A" "rechazado" "$out" \
  "$(rechazado "$out" && [ "$(admin_sql "select must_change_password from users where id='$A'")" = t ]; echo $?)"

out=$(como $A "select record_password_change('$A'); select must_change_password from users where id='$A'")
caso U-19 "A cambia su contraseña" "false" "$out" "$([ "$(echo "$out" | tail -1)" = f ]; echo $?)"

# ── Bitácora ────────────────────────────────────────────────
out=$(como anon "insert into audit_log(email, action, details) values ('$EA','role_changed','falso')")
caso U-20 "Sin sesión: escribir en la bitácora" "rechazado" "$out" "$(rechazado "$out"; echo $?)"

out=$(como $A "select record_admin_action('$ADM','$EA','user_deactivated','falso')")
caso U-21 "A escribe una acción de admin falsa" "rechazado" "$out" "$(rechazado "$out"; echo $?)"

out=$(como $A "select count(*) from audit_log")
caso U-22 "A (lector) lee la bitácora" "0 filas" "$out" "$([ "$out" = 0 ]; echo $?)"

out=$(como $ADM "select count(*) > 0 from audit_log")
caso U-23 "El admin lee la bitácora (eventos reales)" "t" "$out" "$([ "$out" = t ]; echo $?)"

# Segunda barrera: aunque alguien vuelva a dar UPDATE de más por error, el
# trigger users_proteger_columnas sigue frenando el cambio de rol.
admin_sql "grant update (role) on users to authenticated" >/dev/null
out=$(como $A "update users set role='admin' where id='$A'")
admin_sql "revoke update (role) on users from authenticated" >/dev/null
caso U-25 "Con permiso de más por error, A intenta ser admin" "rechazado por el trigger" "$out" \
  "$(echo "$out" | grep -q "No autorizado" && [ "$(admin_sql "select role from users where id='$A'")" = user ]; echo $?)"

# ── Registro de cuentas (trigger handle_new_user) ───────────
NUEVO=dddddddd-0000-0000-0000-00000000000d
admin_sql "delete from users where id='$NUEVO'; delete from auth.users where id='$NUEVO'" >/dev/null
admin_sql "insert into auth.users(id, email, raw_user_meta_data) values ('$NUEVO','nuevo@bw.test','{\"username\":\"nuevo\"}')" >/dev/null
out=$(admin_sql "select username||' '||role from users where id='$NUEVO'")
caso U-24 "Crear cuenta: el trigger crea su perfil de lector" "nuevo user" "$out" "$([ "$out" = "nuevo user" ]; echo $?)"
admin_sql "delete from users where id='$NUEVO'; delete from auth.users where id='$NUEVO'" >/dev/null

echo
echo "Resultado: $pasan pasan, $fallan fallan"
[ "$fallan" -eq 0 ]
