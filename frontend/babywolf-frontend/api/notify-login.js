/**
 * Vercel Serverless Function — POST /api/notify-login
 * Sends a login notification email via Resend.
 */
export default async function handler(req, res) {
  // Only accept POST
  if (req.method !== 'POST') {
    return res.status(405).json({ error: 'Method not allowed' });
  }

  const { email, timestamp } = req.body || {};

  if (!email) {
    return res.status(400).json({ error: 'Email is required' });
  }

  const RESEND_API_KEY = process.env.RESEND_API_KEY;
  if (!RESEND_API_KEY) {
    console.error('RESEND_API_KEY not configured');
    return res.status(500).json({ error: 'Email service not configured' });
  }

  const loginTime = timestamp
    ? new Date(timestamp).toLocaleString('es-MX', { timeZone: 'America/Mexico_City' })
    : new Date().toLocaleString('es-MX', { timeZone: 'America/Mexico_City' });

  try {
    const response = await fetch('https://api.resend.com/emails', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${RESEND_API_KEY}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        from: 'BabyWolf Blog <onboarding@resend.dev>',
        to: [email],
        subject: '🐺 Nuevo inicio de sesión en BabyWolf Blog',
        html: `
          <div style="font-family: 'Courier New', monospace; max-width: 500px; margin: 0 auto; padding: 20px; border: 3px solid #222; background: #fafafa;">
            <h2 style="text-align: center; color: #e74c3c; border-bottom: 2px dashed #222; padding-bottom: 10px;">
              🐺 BabyWolf Blog — Alerta de Seguridad
            </h2>
            <p>Hola,</p>
            <p>Se detectó un <strong>nuevo inicio de sesión</strong> en tu cuenta:</p>
            <table style="width: 100%; border-collapse: collapse; margin: 15px 0;">
              <tr>
                <td style="padding: 8px; border: 1px solid #ddd; font-weight: bold;">📧 Email</td>
                <td style="padding: 8px; border: 1px solid #ddd;">${email}</td>
              </tr>
              <tr>
                <td style="padding: 8px; border: 1px solid #ddd; font-weight: bold;">🕐 Fecha y Hora</td>
                <td style="padding: 8px; border: 1px solid #ddd;">${loginTime}</td>
              </tr>
            </table>
            <p style="color: #666; font-size: 0.9rem;">
              Si no fuiste tú, cambia tu contraseña inmediatamente.
            </p>
            <div style="text-align: center; margin-top: 20px; padding: 10px; background: #222; color: #fff; font-size: 0.8rem;">
              🐺 BabyWolf Blog &copy; 2026
            </div>
          </div>
        `,
      }),
    });

    const data = await response.json();

    if (!response.ok) {
      console.error('Resend API error:', data);
      return res.status(500).json({ error: 'Failed to send email', details: data });
    }

    return res.status(200).json({ success: true, id: data.id });
  } catch (err) {
    console.error('Email send error:', err);
    return res.status(500).json({ error: 'Internal server error' });
  }
}
