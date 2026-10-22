"""Escenario completo de la Práctica 7 contra el backend real + Postgres local.

Uso: python3 escenario.py http://127.0.0.1:18080/api [fase]
fase = "todo" (por defecto) o "persistencia" (sólo TC08 tras reiniciar).
"""
import base64, hashlib, hmac, json, sys, time, urllib.request, urllib.error

API = sys.argv[1]
FASE = sys.argv[2] if len(sys.argv) > 2 else "todo"
ESTADO = "/var/tmp/bwpg/e2e_estado.json"
resultados = []


def jwt_hs256(sub, secret="prueba"):
    b64 = lambda d: base64.urlsafe_b64encode(json.dumps(d, separators=(",", ":")).encode()).rstrip(b"=").decode()
    h, p = b64({"alg": "HS256", "typ": "JWT"}), b64({"sub": sub, "email": "x@bw.test", "exp": int(time.time()) + 3600})
    s = base64.urlsafe_b64encode(hmac.new(secret.encode(), f"{h}.{p}".encode(), hashlib.sha256).digest()).rstrip(b"=").decode()
    return f"{h}.{p}.{s}"


def req(metodo, ruta, cuerpo=None, token=None, jwt=None, url=None):
    u = url or (API + ruta)
    data = json.dumps(cuerpo).encode() if cuerpo is not None else None
    r = urllib.request.Request(u, data=data, method=metodo)
    r.add_header("Content-Type", "application/json")
    if token:
        r.add_header("X-Soporte-Token", token)
    if jwt:
        r.add_header("Authorization", "Bearer " + jwt)
    try:
        with urllib.request.urlopen(r) as resp:
            raw = resp.read()
            return resp.status, (json.loads(raw) if raw else None)
    except urllib.error.HTTPError as e:
        raw = e.read()
        try:
            return e.code, json.loads(raw)
        except Exception:
            return e.code, raw.decode()


def caso(cid, desc, esperado, obtenido, ok, nota=""):
    resultados.append((cid, desc, esperado, obtenido, "PASA" if ok else "FALLA", nota))


def imprimir():
    for r in resultados:
        print(" | ".join(str(x) for x in r))
    fallas = [r for r in resultados if r[4] != "PASA"]
    print(f"\n{len(resultados) - len(fallas)}/{len(resultados)} casos pasan")
    sys.exit(1 if fallas else 0)


if FASE == "persistencia":
    est = json.load(open(ESTADO))
    st, rep = req("GET", "/soporte/reportes/" + est["rep"], token=est["admin"])
    caso("TC08", "Reiniciar el servidor y consultar el reporte guardado", "200 y prioridad Alta",
         f"{st} {rep.get('prioridad') if isinstance(rep, dict) else rep}", st == 200 and rep["prioridad"] == "Alta")
    imprimir()

# ── Sesión del admin y URLs preasignadas de rol ──────────────────
st, r = req("POST", "/admin/soporte/sesion", {"dispositivo": "Pixel 9 (admin)"}, jwt=jwt_hs256("22222222-2222-2222-2222-222222222222"))
caso("P7-R01", "Un usuario del blog que no es admin intenta abrir la mesa", "403", st, st == 403)

st, r = req("POST", "/admin/soporte/sesion", {"dispositivo": "Pixel 9 (admin)"}, jwt=jwt_hs256("11111111-1111-1111-1111-111111111111"))
caso("P7-R02", "El admin del blog abre la mesa de ayuda", "200 + token", st, st == 200 and r.get("token"))
admin, admin_id = r["token"], r["participante"]["id"]

tokens = {}
for rol, nombre in [("agente", "Agente Luis"), ("usuario", "Usuario Gloria"), ("usuario", "Usuario Ana")]:
    st, inv = req("POST", "/soporte/invitaciones", {"rol": rol, "nombre": nombre}, token=admin)
    st2, info = req("GET", "", url=inv["url"])
    st3, j = req("POST", "", {"dispositivo": nombre.split()[1] + " phone"}, url=inv["url"])
    tokens[nombre] = (j["token"], j["participante"]["id"])
    caso("P7-R03", f"Unirse como {rol} con la URL preasignada", "201", f"{st}/{st2}/{st3}", st == 201 and st2 == 200 and st3 == 201 and info["rol"] == rol)
    if rol == "agente":
        st4, _ = req("POST", "", {"dispositivo": "otro"}, url=inv["url"])
        caso("P7-R04", "Reusar la misma URL de invitación en otro dispositivo", "409", st4, st4 == 409)
        st5, _ = req("POST", "", {}, url=inv["url"].replace("rol=agente", "rol=admin"))
        caso("P7-R05", "Cambiar rol=agente por rol=admin en la URL", "403", st5, st5 == 403)

agente, agente_id = tokens["Agente Luis"]
u1, u1_id = tokens["Usuario Gloria"]
u2, u2_id = tokens["Usuario Ana"]

st, _ = req("POST", "/soporte/invitaciones", {"rol": "usuario", "nombre": "Intruso"}, token=u1)
caso("P7-R06", "Un usuario intenta generar invitaciones", "403", st, st == 403)

st, lat = req("POST", "/soporte/latido", token=admin)
conectados = [p["nombre"] for p in lat["participantes"] if p["conectado"]]
caso("P7-R07", "Listar los 4 dispositivos conectados", "4 conectados", len(conectados), len(conectados) == 4, ", ".join(conectados))

# ── Matriz TC01–TC07 ─────────────────────────────────────────────
st, rep = req("POST", "/soporte/reportes", {"id": "REP001", "equipo": "Laptop aula 3", "descripcion": "No enciende", "prioridad": "Media"}, token=u1)
st2, consulta = req("GET", "/soporte/reportes/REP001", token=u1)
caso("TC01", "Alta válida con todos los campos", "201 y se puede consultar", f"{st}/{st2}", st == 201 and st2 == 200)

st, err = req("POST", "/soporte/reportes", {"equipo": "", "descripcion": "  "}, token=u1)
caso("TC02", "Guardar sin equipo ni descripción", "422 con campos", f"{st} {sorted(err.get('campos', {}).keys())}", st == 422 and {"equipo", "descripcion"} <= set(err["campos"]))

st, err = req("POST", "/soporte/reportes", {"id": "REP001", "equipo": "Otro", "descripcion": "Sobrescribir"}, token=u2)
_, sigue = req("GET", "/soporte/reportes/REP001", token=admin)
caso("TC03", "Guardar dos veces REP001", "409 y el original intacto", f"{st} / equipo={sigue['equipo']}", st == 409 and sigue["equipo"] == "Laptop aula 3")

st, _ = req("PATCH", "/soporte/reportes/REP001", {"prioridad": "Alta"}, token=u1)
caso("P7-R08", "El usuario intenta subir su prioridad", "403", st, st == 403)
st, _ = req("POST", "/soporte/reportes/REP001/asignar", {"agente_id": agente_id}, token=admin)
st2, _ = req("PATCH", "/soporte/reportes/REP001", {"prioridad": "Alta"}, token=agente)
st3, otra = req("GET", "/soporte/reportes/REP001", token=agente)
caso("TC04", "Cambiar prioridad Media → Alta y volver a consultar", "Alta", f"{st}/{st2} → {otra['prioridad']}", st == 200 and st2 == 200 and otra["prioridad"] == "Alta")

st, rep2 = req("POST", "/soporte/reportes", {"equipo": "Proyector", "descripcion": "Parpadea"}, token=u2)
caso("P7-R09", "Alta sin ID: el sistema numera REPnnn", "REP002", rep2.get("id"), st == 201 and rep2["id"] == "REP002")
st, _ = req("DELETE", "/soporte/reportes/REP002", token=agente)
caso("P7-R10", "El agente intenta eliminar", "403", st, st == 403)
st, _ = req("DELETE", "/soporte/reportes/REP002", token=admin)
st2, _ = req("GET", "/soporte/reportes/REP002", token=admin)
_, lista = req("GET", "/soporte/reportes", token=admin)
caso("TC05", "Eliminar REP002 y consultarlo de nuevo", "204, luego 404 y fuera de la lista", f"{st}/{st2}", st == 204 and st2 == 404 and all(x["id"] != "REP002" for x in lista))
st, err = req("POST", "/soporte/reportes", {"id": "REP002", "equipo": "x", "descripcion": "y"}, token=u2)
caso("P7-R11", "Reutilizar el ID de un reporte eliminado", "409", st, st == 409)

st, err = req("GET", "/soporte/reportes/REP999", token=admin)
caso("TC06", "Consultar REP999", "404 con mensaje", f"{st} {err.get('error')}", st == 404)

st, _ = req("GET", "/soporte/reportes/REP001", token=u2)
caso("P7-R12", "Un usuario consulta el reporte de otro usuario", "404 (no revela que existe)", st, st == 404)

st, _ = req("GET", "/soporte/monitor", token=u1)
st2, _ = req("GET", "/soporte/monitor", token=agente)
st3, _ = req("GET", "/soporte/monitor", token=admin)
caso("P7-R13", "Tiempo real (monitor) por rol", "usuario 403, agente 200, admin 200", f"{st}/{st2}/{st3}", (st, st2, st3) == (403, 200, 200))

st, rep3 = req("POST", "/soporte/reportes", {"equipo": "Cuenta del blog", "descripcion": "Perdí el acceso", "tipo": "acceso"}, token=u2)
_, vista_ag = req("GET", "/soporte/reportes", token=agente)
caso("P7-R14", "Acceso perdido va directo al admin (el agente no lo ve)", "directo_admin y oculto al agente",
     f"directo={rep3.get('directo_admin')}", rep3.get("directo_admin") and all(x["id"] != rep3["id"] for x in vista_ag))

# ── Mensajes y archivos ─────────────────────────────────────────
st, _ = req("POST", "/soporte/mensajes", {"texto": "   "}, token=u1)
caso("P7-R15", "Enviar mensaje vacío", "422", st, st == 422)
req("POST", "/soporte/mensajes", {"texto": "Hola a todos"}, token=admin)
req("POST", "/soporte/mensajes", {"texto": "Privado para el agente", "para_id": agente_id}, token=admin)
st, sub = req("POST", "/soporte/archivos/subida", {"nombre": "../foto falla.PNG", "tamano": 2048}, token=u1)
caso("P7-R16", "Pedir URL preasignada de subida", "200 y ruta en carpeta propia", f"{st} {sub.get('path')}", st == 200 and sub["path"].startswith(f"soporte/{u1_id}/") and sub["upload_url"].startswith("http"))
st, _ = req("POST", "/soporte/archivos/subida", {"nombre": "script.exe", "tamano": 10}, token=u1)
caso("P7-R17", "Subir un .exe", "422", st, st == 422)
st, m = req("POST", "/soporte/mensajes", {"texto": "Adjunto foto", "archivo_path": sub["path"], "archivo_nombre": sub["nombre"], "reporte_id": "rep001"}, token=u1)
caso("P7-R18", "Mensaje con adjunto ligado a REP001", "201", st, st == 201 and m["reporte_id"] == "REP001")
st, _ = req("POST", "/soporte/mensajes", {"texto": "x", "archivo_path": sub["path"]}, token=u2)
caso("P7-R19", "Adjuntar el archivo de otro participante", "422", st, st == 422)
st, d = req("POST", "/soporte/archivos/descarga", {"path": sub["path"]}, token=admin)
st2, _ = req("POST", "/soporte/archivos/descarga", {"path": sub["path"]}, token=u2)
caso("P7-R20", "Descarga firmada de un adjunto visible (admin y usuario)", "200 / 200", f"{st}/{st2}", st == 200 and st2 == 200, "el mensaje con adjunto es para todos, por eso u2 también lo ve")

_, ms_u2 = req("GET", "/soporte/mensajes", token=u2)
_, ms_ag = req("GET", "/soporte/mensajes", token=agente)
priv = lambda ms: any(x["texto"] == "Privado para el agente" for x in ms["mensajes"])
caso("P7-R21", "Mensaje privado: lo ve el agente, no otro usuario", "agente sí / usuario no", f"{priv(ms_ag)}/{priv(ms_u2)}", priv(ms_ag) and not priv(ms_u2))
_, ms_ad = req("GET", "/soporte/mensajes", token=admin)
conteo = {c["emisor_nombre"]: c["total"] for c in ms_ad["conteo"]}
caso("P7-R22", "Contador de mensajes por dispositivo", "abraham 2, Gloria 1", conteo, conteo.get("abraham") == 2 and conteo.get("Usuario Gloria") == 1)

# ── Intercambio de roles por desconexión ────────────────────────
req("POST", "/soporte/salir", token=admin)
st, lat = req("POST", "/soporte/latido", token=agente)
st2, inv = req("POST", "/soporte/invitaciones", {"rol": "usuario", "nombre": "Nuevo"}, token=agente)
caso("P7-R23", "Admin se desconecta: el agente cubre como admin temporal", "admin_temporal y puede invitar", f"{lat['rol_efectivo']} / {st2}", lat["rol_efectivo"] == "admin_temporal" and st2 == 201)
st, _ = req("DELETE", "/soporte/reportes/REP001", token=agente)
caso("P7-R24", "El admin temporal intenta eliminar", "403 (eliminar es sólo del admin real)", st, st == 403)
st, lat = req("POST", "/soporte/latido", token=admin)
_, lat_ag = req("POST", "/soporte/latido", token=agente)
caso("P7-R25", "El admin vuelve: el agente regresa a su rol", "agente", lat_ag["rol_efectivo"], lat_ag["rol_efectivo"] == "agente")
req("POST", "/soporte/salir", token=agente)
st, n = req("POST", "/soporte/reportes/asumir-desconectados", token=admin)
_, r1 = req("GET", "/soporte/reportes/REP001", token=admin)
caso("P7-R26", "El agente se desconecta: el admin asume sus reportes", "REP001 pasa al admin", f"{n} / agente_id={'admin' if r1['agente_id'] == admin_id else r1['agente_id']}", st == 200 and n["reasignados"] >= 1 and r1["agente_id"] == admin_id)

st, _ = req("DELETE", "/soporte/mensajes", token=u1)
st2, _ = req("DELETE", "/soporte/mensajes", token=admin)
_, vacio = req("GET", "/soporte/mensajes", token=admin)
caso("P7-R27", "Limpiar historial: usuario no, admin sí", "403 / 204 y vacío", f"{st}/{st2}", st == 403 and st2 == 204 and vacio["mensajes"] == [])

st, _ = req("DELETE", f"/soporte/participantes/{u2_id}", token=admin)
st2, _ = req("POST", "/soporte/latido", token=u2)
caso("P7-R28", "El admin expulsa a un usuario y su token deja de servir", "204 / 401", f"{st}/{st2}", st == 204 and st2 == 401)

st, _ = req("GET", "/soporte/reportes", token=admin[:-2] + "xx")
caso("P7-R29", "Token alterado", "401", st, st == 401)

json.dump({"admin": admin, "rep": "REP001"}, open(ESTADO, "w"))
imprimir()
