import { createClient } from '@supabase/supabase-js';
import pkg from 'pg';
const { Client } = pkg;
import * as fs from 'fs';
import { fileURLToPath } from 'url';
import * as path from 'path';

const SUPABASE_URL = 'https://uwzynkvnqcupthvrtkry.supabase.co';
const SUPABASE_SERVICE_ROLE_KEY = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InV3enlua3ZucWN1cHRodnJ0a3J5Iiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTc3NDMzMTY0NywiZXhwIjoyMDg5OTA3NjQ3fQ.HYt9y9KLu1n24H60urawFYhWZu0KQX2Ujik4QHCTUD8';
const DATABASE_URL = 'postgresql://postgres:Yukiylu03!!@db.uwzynkvnqcupthvrtkry.supabase.co:5432/postgres';

const supabase = createClient(SUPABASE_URL, SUPABASE_SERVICE_ROLE_KEY);

async function runAll() {
    console.log("🚀 Iniciando automatización completa para Supabase...");

    let userId = null;

    // 1. Crear usuario admin si no existe usando Auth Admin API (Solo disponible con Service Role)
    const email = 'abrajam285@gmail.com';
    const password = 'AdminPassword123!';
    
    console.log(`👤 Creando/Verificando cuenta de Auth: ${email}...`);
    // buscar si el usuario existe
    const { data: usersData, error: listErr } = await supabase.auth.admin.listUsers();
    if (!listErr && usersData && usersData.users) {
        const existing = usersData.users.find(u => u.email === email);
        if (existing) {
            console.log(`✓ El usuario ya existe en Auth (ID: ${existing.id})`);
            userId = existing.id;
        }
    }

    if (!userId) {
        console.log(`Registrando nuevo usuario admin...`);
        const { data: authData, error: authErr } = await supabase.auth.admin.createUser({
            email,
            password,
            email_confirm: true,
            user_metadata: { role: 'admin' }
        });
        
        if (authErr) {
            console.error("❌ Error al crear usuario en Auth:", authErr.message);
            return;
        }
        userId = authData.user.id;
        console.log(`✓ Usuario creado con éxito (ID: ${userId})`);
    }

    // 2. Ejecutar SQL para tablas, RLS e insertar Seed usando el ID real
    console.log("📦 Conectando a PostgreSQL para procesar el esquema y RLS...");
    const client = new Client({
        connectionString: DATABASE_URL,
    });
    
    try {
        await client.connect();
        
        // Leemos el archivo sql que generamos antes
        const __filename = fileURLToPath(import.meta.url);
        const __dirname = path.dirname(__filename);
        let sqlContent = fs.readFileSync(path.join(__dirname, '../../supabase_setup_and_seed.sql'), 'utf-8');
        
        // Reemplazamos el ID por defecto por el ID real
        sqlContent = sqlContent.replace(/uuid-del-admin-aqui/g, userId);
        
        console.log("📜 Ejecutando Script SQL (Tablas + RLS + Datos)...");
        await client.query(sqlContent);
        
        console.log("✅ Script SQL ejecutado exitosamente.");
        
    } catch (dbErr) {
        console.error("❌ Error al ejecutar el SQL:", dbErr);
    } finally {
        await client.end();
    }

    console.log("🎉 ¡Proceso completado al 100%!");
}

runAll().catch(console.error);
