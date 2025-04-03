import api from './api';

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface User {
  id: number;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

const AuthService = {
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await api.post<AuthResponse>('/api/register', data);
    console.log('Réponse d\'inscription:', response.data);
    // Stocker le token et les infos utilisateur dans le localStorage
    if (response.data && response.data.token) {
      localStorage.setItem('token', response.data.token);
      console.log('Token stocké:', response.data.token);
    } else {
      console.error('Pas de token dans la réponse d\'inscription');
    }
    if (response.data && response.data.user) {
      localStorage.setItem('user', JSON.stringify(response.data.user));
    }
    return response.data;
  },

  login: async (data: LoginRequest): Promise<AuthResponse> => {
    // Vérifier les données envoyées
    console.log('Données de connexion envoyées:', data);
    
    // Formater les données pour s'assurer qu'elles correspondent au format attendu par le backend
    const loginData = {
      email: data.email,
      password: data.password
    };
    
    try {
      const response = await api.post<AuthResponse>('/api/login', loginData);
      console.log('Réponse de connexion:', response.data);
      
      // Stocker le token et les infos utilisateur dans le localStorage
      if (response.data && response.data.token) {
        localStorage.setItem('token', response.data.token);
        console.log('Token stocké:', response.data.token);
      } else {
        console.error('Pas de token dans la réponse de connexion');
      }
      
      if (response.data && response.data.user) {
        localStorage.setItem('user', JSON.stringify(response.data.user));
      }
      
      return response.data;
    } catch (error: any) {
      console.error('Erreur de connexion détaillée:', {
        message: error.message,
        response: error.response?.data,
        status: error.response?.status
      });
      throw error;
    }
  },

  logout: (): void => {
    // Supprimer le token et les infos utilisateur du localStorage
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  },

  getCurrentUser: (): User | null => {
    const userStr = localStorage.getItem('user');
    if (userStr) {
      return JSON.parse(userStr) as User;
    }
    return null;
  },

  isAuthenticated: (): boolean => {
    return localStorage.getItem('token') !== null;
  }
};

export default AuthService;
