// services/ApiClient.js
import axios from 'axios';

// API Base URL aus den ENV-Variablen laden
const API_BASE_URL = 'http://localhost:9000'
// Axios-Instanz erstellen
const apiClient = axios.create({
    baseURL: API_BASE_URL,
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
    },
});

export default apiClient;