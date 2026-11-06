#!/usr/bin/env bash
# Prácticas 6, 9 y 10 — reglas de 007_app_lectores.sql vistas desde la base.
#
# Mismo método que rls_incidencias.sh: cada operación corre con el rol de
# Supabase ("authenticated" con el JWT del usuario, o "anon") y la base
# acepta o rechaza igual que en Supabase.
#
# Uso: PSQL="psql -d bwapp" ./app_lectores.sh
# (base con 001, 002, 006 y 007, y los usuarios A, B y ADMIN; ver README.md)

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

caso() {
  if [ "$5" -eq 0 ]; then r=PASA; pasan=$((pasan+1)); else r=FALLA; fallan=$((fallan+1)); fi
  printf '%s | %s | %s | %s | %s\n' "$1" "$2" "$3" "$(echo "$4" | tr '\n' ' ' | cut -c1-90)" "$r"
}

admin_sql "truncate public.incidencias, public.comments, public.posts cascade" >/dev/null
POST=$(admin_sql "insert into posts(author_id,title,slug,content,published) values ('$ADM','Zelda retro','zelda-retro','...',true) returning id")

# ── Práctica 9: base compartida entre dispositivos ──────────
RA=$(como $A "insert into incidencias(titulo, descripcion, post_slug, post_titulo) values ('PRUEBA_A_001','Imagen rota','zelda-retro','Zelda retro') returning id")
caso DB-02 "A crea un reporte (sin mandar propietario)" "queda en la base con dueño A, versión 1" \
  "$(admin_sql "select creado_por||' v'||version||' '||prioridad from incidencias where id='$RA'")" \
  "$([ "$(admin_sql "select creado_por||' v'||version from incidencias where id='$RA'")" = "$A v1" ]; echo $?)"

out=$(como $ADM "select titulo from incidencias where id = '$RA'")
caso DB-03 "El admin (otro dispositivo) lee el reporte de A" "lo ve" "$out" "$([ "$out" = PRUEBA_A_001 ]; echo $?)"

RB=$(como $B "insert into incidencias(titulo) values ('PRUEBA_B_001') returning id")
out=$(como $ADM "select count(*) from incidencias")
caso ADM-01 "El admin ve los reportes de todos" "2" "$out" "$([ "$out" = 2 ]; echo $?)"

out=$(como $B "select count(*) from incidencias")
caso SEG-05b "B (lector) no ve el reporte de A" "1 (sólo el suyo)" "$out" "$([ "$out" = 1 ]; echo $?)"

out=$(como $ADM "update incidencias set estado='en_proceso', prioridad='critica' where id='$RA' and version=1 returning estado||' '||prioridad||' v'||version")
visto=$(como $A "select estado||' '||prioridad||' v'||version from incidencias where id='$RA'")
caso DB-06 "El admin atiende y A ve el cambio" "en_proceso critica v2" "admin: $out · A ve: $visto" \
  "$([ "$visto" = "en_proceso critica v2" ]; echo $?)"

quien=$(admin_sql "select actualizado_por from incidencias where id='$RA'")
caso DB-06b "Queda registrado quién hizo el último cambio" "UID del admin" "$quien" "$([ "$quien" = "$ADM" ]; echo $?)"

# DB-07: A tenía abierta la versión 2. B no puede tocarlo, así que el
# segundo editor es el admin otra vez: guarda primero, luego A con la v2.
como $ADM "update incidencias set prioridad='alta' where id='$RA' and version=2" >/dev/null
out=$(como $A "with u as (update incidencias set descripcion='Imagen rota en la portada' where id='$RA' and version=2 returning 1) select count(*) from u")
final=$(admin_sql "select prioridad||' v'||version||' '||descripcion from incidencias where id='$RA'")
caso DB-07 "Guardar con una versión vieja (otro guardó antes)" "0 filas: no se pisa el cambio ajeno" "$out filas; queda: $final" \
  "$([ "$out" = 0 ] && [ "$final" = "alta v3 Imagen rota" ]; echo $?)"

out=$(como $A "update incidencias set descripcion='Imagen rota en la portada' where id='$RA' and version=3 returning descripcion||' v'||version")
caso DB-07b "Al recargar (v3) A sí puede guardar" "v4 con los dos cambios" "$out / $(admin_sql "select prioridad from incidencias where id='$RA'")" \
  "$([ "$out" = "Imagen rota en la portada v4" ]; echo $?)"

out=$(como $B "with u as (update incidencias set estado='cerrada' where id='$RA' returning 1) select count(*) from u")
caso SEG-06b "B (lector) intenta cambiar el reporte de A" "0 filas" "$out" "$([ "$out" = 0 ]; echo $?)"

out=$(como $A "update incidencias set prioridad='urgente' where id='$RA'")
caso DB-RN "Prioridad fuera del catálogo" "rechazada por la base" "$out" "$(echo "$out" | grep -q 'incidencias_prioridad_check'; echo $?)"

out=$(como $ADM "update incidencias set creado_por='$A' where id='$RB'")
dueno=$(admin_sql "select creado_por from incidencias where id='$RB'")
caso ADM-02 "El admin intenta pasar el reporte de B a A" "rechazado, sigue siendo de B" "$out" \
  "$(echo "$out" | grep -q 'no se puede cambiar' && [ "$dueno" = "$B" ]; echo $?)"

out=$(como $ADM "with d as (delete from incidencias where id='$RA' returning 1) select count(*) from d")
caso ADM-03 "El admin intenta borrar un reporte ajeno" "0 filas (sólo lo cierra)" "$out" "$([ "$out" = 0 ]; echo $?)"

# ── Práctica 6: comentarios ─────────────────────────────────
out=$(como $A "insert into comments(post_id, content) values ('$POST', 'Gran nota') returning author_id")
caso COM-01 "A comenta (autor desde el JWT)" "author_id = A" "$out" "$([ "$out" = "$A" ]; echo $?)"

out=$(como anon "select count(*) from comments where post_id='$POST'")
caso COM-02 "Sin sesión se pueden leer los comentarios" "1" "$out" "$([ "$out" = 1 ]; echo $?)"

out=$(como anon "insert into comments(post_id, author_id, content) values ('$POST', '$A', 'anónimo')")
caso COM-03 "Sin sesión no se puede comentar" "rechazado" "$out" "$(echo "$out" | grep -qE 'row-level security|permission denied'; echo $?)"

out=$(como $B "insert into comments(post_id, author_id, content) values ('$POST', '$A', 'suplantado')")
caso COM-04 "B comenta a nombre de A" "rechazado por RLS" "$out" "$(echo "$out" | grep -q 'row-level security'; echo $?)"

out=$(como $B "insert into comments(post_id, content) values ('$POST', '   ')")
caso COM-05 "Comentario vacío" "rechazado por la base" "$out" "$(echo "$out" | grep -q 'comments_content_check'; echo $?)"

largo=$(printf 'x%.0s' $(seq 1 281))
out=$(como $B "insert into comments(post_id, content) values ('$POST', '$largo')")
caso COM-06 "Comentario de 281 caracteres" "rechazado por la base" "$out" "$(echo "$out" | grep -q 'comments_content_check'; echo $?)"

CA=$(admin_sql "select id from comments where author_id='$A' limit 1")
out=$(como $B "with d as (delete from comments where id='$CA' returning 1) select count(*) from d")
caso COM-07 "B borra el comentario de A" "0 filas" "$out" "$([ "$out" = 0 ]; echo $?)"

out=$(como $ADM "with d as (delete from comments where id='$CA' returning 1) select count(*) from d")
caso COM-08 "El admin modera (borra) el comentario de A" "1 fila" "$out" "$([ "$out" = 1 ]; echo $?)"

echo
echo "$pasan pasan, $fallan fallan"
[ "$fallan" -eq 0 ]
