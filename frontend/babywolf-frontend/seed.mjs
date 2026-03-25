import { createClient } from '@supabase/supabase-js';

const SUPABASE_URL = 'https://uwzynkvnqcupthvrtkry.supabase.co';
const SUPABASE_ANON_KEY = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InV3enlua3ZucWN1cHRodnJ0a3J5Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NzQzMzE2NDcsImV4cCI6MjA4OTkwNzY0N30.H9cmX7cq5WFzOnpseB6Jay5-aSFh-nwc39uEQbCQwfg';

const supabase = createClient(SUPABASE_URL, SUPABASE_ANON_KEY);

async function runFaseA() {
    console.log("=== FASE A: SEEDING REAL ===");

    // 1. Registrar usuario admin (Fase A)
    const email = 'abrajam285@gmail.com';
    const password = 'AdminPassword123!';
    
    console.log(`Registrando usuario ${email}...`);
    let { data: authData, error: authErr } = await supabase.auth.signUp({
        email,
        password,
    });
    
    // Si ya existe, tratamos de hacer login
    if (authErr && authErr.message.includes('already registered')) {
        console.log("El usuario ya existe, iniciando sesión...");
        const loginRes = await supabase.auth.signInWithPassword({ email, password });
        authData = loginRes.data;
        authErr = loginRes.error;
    }

    if (authErr) {
        console.error("Error al autenticar:", authErr.message);
        return;
    }

    const userId = authData.user?.id;
    if (!userId) {
        console.error("No se pudo obtener el ID del usuario.");
        return;
    }

    console.log(`Usuario autenticado con ID: ${userId}`);

    // 2. Insertar registro en public.users con role=admin
    // Primero comprobamos si existe
    const { data: existingUser } = await supabase.from('users').select('*').eq('id', userId).single();
    if (!existingUser) {
        console.log("Insertando en public.users...");
        const { error: userErr } = await supabase.from('users').insert({
            id: userId,
            email: email,
            username: 'admin_abrajam',
            role: 'admin'
        });
        if (userErr) console.error("Error al insertar public.user:", userErr.message);
    } else {
        console.log("El usuario ya existe en public.users, asegurando role=admin...");
        await supabase.from('users').update({ role: 'admin' }).eq('id', userId);
    }

    // 3. Insertar 6+ posts de prueba con cover images reales
    console.log("Insertando posts de prueba...");
    const posts = [
        {
            title: "El Renacer de los 8-bits",
            slug: "renacer-8-bits",
            content: "<p>Los juegos de 8-bits están volviendo con fuerza. Desarrolladores indie están creando joyas que superan a muchos juegos AAA modernos en diversión pura.</p>",
            category: "retro",
            cover_image_url: "https://images.unsplash.com/photo-1550745165-9bc0b252726f?w=800&q=80",
            author_id: userId,
            published: true
        },
        {
            title: "Review: The Legend of Zelda Tears of the Kingdom",
            slug: "zelda-totk-review",
            content: "<p>Una obra maestra absoluta que redefine la física y la creatividad en los videojuegos.</p>",
            category: "gaming",
            cover_image_url: "https://images.unsplash.com/photo-1588636545233-1ec91a1a5b82?w=800&q=80",
            author_id: userId,
            published: true
        },
        {
            title: "Desarrollo con Angular 17",
            slug: "desarrollo-angular-17",
            content: "<p>Las nuevas Signals han cambiado para siempre la forma en que manejamos el estado en Angular. Es más fácil, limpio y rápido.</p>",
            category: "tech",
            cover_image_url: "https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=800&q=80",
            author_id: userId,
            published: true
        },
        {
            title: "El fin del E3: ¿Qué sigue para la industria?",
            slug: "fin-del-e3",
            content: "<p>Con la cancelación permanente del E3, el Summer Game Fest y los eventos digitales directos toman la corona.</p>",
            category: "gaming",
            cover_image_url: "https://images.unsplash.com/photo-1542751371-adc38448a05e?w=800&q=80",
            author_id: userId,
            published: true
        },
        {
            title: "Supabase vs Firebase: La batalla definitiva",
            slug: "supabase-vs-firebase",
            content: "<p>PostgreSQL bajo el capó hace que Supabase sea una alternativa increíblemente potente frente a las bases de datos NoSQL de Firebase.</p>",
            category: "tech",
            cover_image_url: "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=800&q=80",
            author_id: userId,
            published: true
        },
        {
            title: "Coleccionismo de Consolas Clásicas",
            slug: "coleccionismo-consolas",
            content: "<p>¿Vale la pena invertir en consolas originales o la emulación es suficiente? Analizamos el mercado de segundamano.</p>",
            category: "retro",
            cover_image_url: "https://images.unsplash.com/photo-1605333555230-01967269b8fa?w=800&q=80",
            author_id: userId,
            published: true
        }
    ];

    const { data: insertedPosts, error: postsErr } = await supabase.from('posts').upsert(posts, { onConflict: 'slug' }).select();
    
    if (postsErr) {
        console.error("Error al insertar posts:", postsErr.message);
    } else {
        console.log(`¡Se insertaron/actualizaron ${insertedPosts?.length} posts!`);
        
        // 4. Insertar comentarios de prueba
        if (insertedPosts && insertedPosts.length > 0) {
            console.log("Insertando comentarios de prueba...");
            const postId = insertedPosts[0].id;
            const comments = [
                { post_id: postId, author_id: userId, content: "¡Excelente artículo, totalmente de acuerdo!" },
                { post_id: postId, author_id: userId, content: "Yo sigo prefiriendo jugar en hardware original, la latencia de la emulación me molesta mucho." },
                { post_id: postId, author_id: userId, content: "Gran lectura. Saludos desde la comunidad." }
            ];
            
            const { error: commErr } = await supabase.from('comments').insert(comments);
            if (commErr) {
                console.error("Error al insertar comentarios:", commErr.message);
            } else {
                console.log("¡Comentarios insertados!");
            }
        }
    }

    console.log("=== FASE A COMPLETADA ===");
}

runFaseA();
