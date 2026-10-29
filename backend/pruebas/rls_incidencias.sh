#!/usr/bin/env bash
# Práctica 10 — matriz SEG-03…SEG-10 contra las reglas RLS de 006_incidencias_rls.sql.
#
# Corre cada operación como lo hace Supabase: con el rol "authenticated" (o
# "anon" sin sesión) y el JWT del usuario en request.jwt.claims. Si una regla
# falla, la base lo rechaza aquí exactamente igual que en Supabase.
#
# Uso: PSQL="psql -h localhost -d bw10" ./rls_incidencias.sh
# (la base debe tener los roles anon/authenticated, auth.users, auth.uid(),
#  001_initial_schema.sql y 006_incidencias_rls.sql; ver pruebas/README.md)

set -u
PSQL=${PSQL:-psql}
A=aaaaaaaa-0000-0000-0000-00000000000a
B=bbbbbbbb-0000-0000-0000-00000000000b
pasan=0; fallan=0

# como <uid|anon> <sql>  → imprime la salida o el error de Postgres
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

caso() { # caso <id> <descripción> <esperado> <obtenido> <condición: 0 = pasa>
  if [ "$5" -eq 0 ]; then r=PASA; pasan=$((pasan+1)); else r=FALLA; fallan=$((fallan+1)); fi
  printf '%s | %s | %s | %s | %s\n' "$1" "$2" "$3" "$(echo "$4" | tr '\n' ' ' | cut -c1-90)" "$r"
}

$PSQL -X -q -c "truncate public.incidencias" >/dev/null

# Datos de partida: dos incidencias de A y una de B (Fase 2).
IA=$(como $A "insert into incidencias(titulo, creado_por) values ('Pantalla rota', '$A') returning id")
como $A "insert into incidencias(titulo, creado_por) values ('Sin red en aula 3', '$A')" >/dev/null
IB=$(como $B "insert into incidencias(titulo, creado_por) values ('Proyector apagado', '$B') returning id")

out=$(como $A "select creado_por from incidencias where id = '$IA'")
caso SEG-03 "A registra una incidencia con su UID" "creada con propietario A" "$out" "$([ "$out" = "$A" ]; echo $?)"

out=$(como $A "insert into incidencias(titulo, creado_por) values ('Suplantación', '$B')")
caso SEG-04 "A crea una incidencia con el UID de B" "rechazada por RLS" "$out" "$(echo "$out" | grep -q 'row-level security'; echo $?)"

out=$(como $A "select count(*) from incidencias where id = '$IB'")
caso SEG-05 "A lee directamente la incidencia de B" "0 filas visibles" "$out" "$([ "$out" = 0 ]; echo $?)"

out=$(como $A "with u as (update incidencias set estado = 'cerrada' where id = '$IB' returning 1) select count(*) from u")
sigue=$($PSQL -X -q -A -t -c "select estado from incidencias where id = '$IB'")
caso SEG-06 "A modifica la incidencia de B" "0 filas cambiadas, B intacta" "$out filas; estado de B=$sigue" "$([ "$out" = 0 ] && [ "$sigue" = abierta ]; echo $?)"

out=$(como $A "with d as (delete from incidencias where id = '$IB' returning 1) select count(*) from d")
existe=$($PSQL -X -q -A -t -c "select count(*) from incidencias where id = '$IB'")
caso SEG-07 "A elimina la incidencia de B" "0 filas borradas, B existe" "$out filas; B existe=$existe" "$([ "$out" = 0 ] && [ "$existe" = 1 ]; echo $?)"

out=$(como $A "update incidencias set estado = 'en_proceso' where id = '$IA' returning estado")
caso SEG-08 "A cambia el estado de su incidencia" "permitido" "$out" "$([ "$out" = en_proceso ]; echo $?)"

out=$(como $A "update incidencias set creado_por = '$B' where id = '$IA'")
dueno=$($PSQL -X -q -A -t -c "select creado_por from incidencias where id = '$IA'")
caso SEG-09 "A pasa su incidencia a B (cambio de propietario)" "rechazado, sigue siendo de A" "$out" "$(echo "$out" | grep -q 'row-level security' && [ "$dueno" = "$A" ]; echo $?)"

out=$(como anon "select count(*) from incidencias")
out2=$(como anon "insert into incidencias(titulo, creado_por) values ('Anónima', '$A')")
caso SEG-10 "Sin sesión: leer y escribir" "permission denied" "$out / $out2" "$(echo "$out" | grep -q 'permission denied' && echo "$out2" | grep -q 'permission denied'; echo $?)"

out=$(como $B "select count(*) from incidencias")
caso EXTRA "B sólo ve lo suyo (A tiene 2)" "1" "$out" "$([ "$out" = 1 ]; echo $?)"

echo
echo "$pasan pasan, $fallan fallan"
[ "$fallan" -eq 0 ]
