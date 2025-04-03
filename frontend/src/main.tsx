import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';
import App from './App.tsx';

// Créer un fichier .env avec l'URL de l'API
if (import.meta.env.DEV) {
  window.addEventListener('error', (event) => {
    // Ignorer les erreurs de Tailwind CSS qui peuvent apparaitre pendant le développement
    if (event.message.includes('@tailwind') || event.message.includes('@apply')) {
      event.preventDefault();
      return;
    }
  });
}

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
