<script setup lang="ts">
import CounterCard from './modules/counter/components/CounterCard.vue'
</script>

<template>
  <div class="page">
    <div class="glow glow-1" />
    <div class="glow glow-2" />

    <header class="header">
      <div class="logo">
        <span class="logo-icon">⬡</span>
        <span class="logo-text">CounterApp</span>
      </div>
    </header>

    <main class="main">
      <div class="hero-text">
        <h1 class="title">Dynamic Realtime Rooms</h1>
        <p class="subtitle">
          Type any name below to fetch its state. If it doesn't exist, Postgres will magically mint a new atomic counter row the moment you fire a transaction!
        </p>
      </div>

      <!-- Pluggable Single Responsibility Module injected safely! -->
      <CounterCard />
    </main>

    <footer class="footer">
      <p>Running on localhost</p>
    </footer>
  </div>
</template>

<style scoped>
/* Abstracted App Macro Layout Styles */
.page { min-height: 100vh; display: flex; flex-direction: column; align-items: center; position: relative; overflow: hidden; }

.glow { position: fixed; border-radius: 50%; filter: blur(120px); pointer-events: none; z-index: 0; }
.glow-1 { width: 600px; height: 600px; background: rgba(109, 40, 217, 0.15); top: -200px; left: -200px; }
.glow-2 { width: 500px; height: 500px; background: rgba(34, 211, 238, 0.08); bottom: -150px; right: -150px; }

.header { width: 100%; max-width: 900px; display: flex; justify-content: space-between; align-items: center; padding: 24px 32px; z-index: 1; }
.logo { display: flex; align-items: center; gap: 10px; }
.logo-icon { font-size: 1.6rem; color: #7c3aed; line-height: 1; }
.logo-text { font-size: 1.1rem; font-weight: 700; letter-spacing: -0.5px; color: #fafafa; }
.stack-badges { display: flex; gap: 8px; }
.badge { font-size: 0.72rem; font-weight: 600; padding: 3px 10px; border-radius: 999px; letter-spacing: 0.3px; }
.badge-vue { background: rgba(65,184,131,0.15); color: #41b883; border: 1px solid rgba(65,184,131,0.3); }
.badge-go  { background: rgba(0,173,216,0.15);  color: #00add8; border: 1px solid rgba(0,173,216,0.3); }
.badge-pg  { background: rgba(51,103,145,0.2);  color: #7ab4dc; border: 1px solid rgba(51,103,145,0.4); }

.main { flex: 1; width: 100%; max-width: 900px; padding: 24px 32px 48px; display: flex; flex-direction: column; align-items: center; gap: 40px; z-index: 1; }

.hero-text { text-align: center; }
.title { font-size: clamp(1.8rem, 4vw, 2.6rem); font-weight: 800; letter-spacing: -1.5px; background: linear-gradient(135deg, #fafafa 30%, #a78bfa 100%); -webkit-background-clip: text; -webkit-text-fill-color: transparent; background-clip: text; margin-bottom: 12px; }
.subtitle { font-size: 1.05rem; color: #a1a1aa; max-width: 560px; line-height: 1.6; }

.info-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); gap: 16px; width: 100%; }
.info-card { background: rgba(24, 24, 27, 0.6); border: 1px solid #27272a; border-radius: 16px; padding: 24px; display: flex; flex-direction: column; gap: 10px; transition: border-color 0.2s, transform 0.2s; backdrop-filter: blur(10px); }
.info-card:hover { border-color: #52525b; transform: translateY(-3px); }
.info-icon { font-size: 1.5rem; }
.icon-vue { color: #41b883; }
.icon-go  { color: #00add8; }
.icon-pg  { color: #7ab4dc; }
.info-card h3 { font-size: 0.95rem; font-weight: 700; color: #e4e4e7; }
.info-card p { font-size: 0.85rem; color: #71717a; line-height: 1.6; }

.footer { padding: 20px; color: #3f3f46; font-size: 0.8rem; z-index: 1; }
</style>
