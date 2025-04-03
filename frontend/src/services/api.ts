import axios from 'axios';

// Créer une instance axios avec l'URL de base de l'API
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
});

// Intercepteur pour ajouter le token JWT à chaque requête
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  console.log('Token récupéré du localStorage:', token);
  
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
    console.log('En-tête Authorization ajouté:', `Bearer ${token}`);
  } else {
    console.warn('Aucun token trouvé dans le localStorage');
  }
  
  console.log('Configuration de la requête:', {
    url: config.url,
    method: config.method,
    headers: config.headers,
    data: config.data
  });
  
  return config;
});

// Intercepteur pour gérer les réponses et les erreurs
api.interceptors.response.use(
  (response) => {
    console.log('Réponse API réussie:', {
      url: response.config.url,
      status: response.status,
      data: response.data
    });
    return response;
  },
  (error) => {
    console.error('Erreur API:', {
      url: error.config?.url,
      status: error.response?.status,
      data: error.response?.data,
      message: error.message
    });
    return Promise.reject(error);
  }
);

export default api;
