const CACHE_NAME = 'boardgame-engine-v1';
const ASSETS = [
  '/',
  '/index.html',
  '/main.wasm',
  '/wasm_exec.js',
  '/output.css',
  '/manifest.json',
  'https://unpkg.com/htmx.org@2.0.0'
];

importScripts('wasm_exec.js');

const go = new Go();
let wasmInstance;

async function loadWasm() {
    if (wasmInstance) return;
    const response = await fetch('main.wasm');
    const buffer = await response.arrayBuffer();
    const result = await WebAssembly.instantiate(buffer, go.importObject);
    wasmInstance = result.instance;
    go.run(wasmInstance);
}

self.addEventListener('install', (event) => {
    event.waitUntil(
        caches.open(CACHE_NAME).then((cache) => cache.addAll(ASSETS))
    );
});

self.addEventListener('activate', (event) => {
    event.waitUntil(self.clients.claim());
});

self.addEventListener('fetch', (event) => {
    const url = new URL(event.request.url);
    const isHX = event.request.headers.get('HX-Request') === 'true';

    // Intercept HTMX requests or root navigation
    if (isHX || url.pathname === '/' || url.pathname === '/slots') {
        event.respondWith((async () => {
            await loadWasm();
            const html = self.handleRequest(event.request.method, event.request.url, isHX);
            return new Response(html, {
                headers: { 'Content-Type': 'text/html' }
            });
        })());
    } else {
        // Cache-first strategy for static assets
        event.respondWith(
            caches.match(event.request).then((response) => {
                return response || fetch(event.request);
            })
        );
    }
});
