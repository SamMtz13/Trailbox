// src/main.ts
import './app.css';         // 👈 IMPORTA LOS ESTILOS GLOBALES

import App from './App.svelte';
import { mount } from 'svelte';

const target = document.getElementById('app');

if (!target) {
  throw new Error('No se encontró el elemento #app');
}

mount(App, {
  target
});
