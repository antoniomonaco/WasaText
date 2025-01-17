<template>
  <div class="login-view">
     <div class="login-container">
         <div class="whatsapp-logo">
             <img src="/logo_whatsapp.svg" alt="WhatsApp Logo">
         </div>
         <h1>Accedi a WASAText</h1>
         <form @submit.prevent="login">
             <input 
                 type="text" 
                 v-model="username" 
                 placeholder="Inserisci il tuo nome" 
                 required 
                 class="login-input"
             >
             <button type="submit" class="login-button">Inizia a chattare</button>
         </form>
         <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
         <div class="disclaimer">WASAText è un progetto universitario sviluppato da Antonio Monaco.</div>
     </div>
  </div>
 </template>
 
 <script>
 export default {
     data() {
         return {
             username: '',
             errorMessage: null,
         };
     },
     methods: {
         async login() {
             this.errorMessage = null;
             try {
                 const response = await this.$axios.post('/session', { name: this.username });
                 const token = response.data.identifier;
                 localStorage.setItem('authToken', token);
                 this.$router.push({ name: 'Main' });
             } catch (error) {
                 if (error.response) {
                     this.errorMessage = error.response.data;
                 } else if (error.request) {
                     this.errorMessage = "Errore di rete";
                 } else {
                     this.errorMessage = "Errore nella richiesta";
                 }
             }
         },
     },
 };
 </script>
 
 <style scoped>
 .login-view {
     display: flex;
     justify-content: center;
     align-items: center;
     height: 100vh;
     background-color: #111b21;
 }
 
 .login-container {
     background-color: #202c33;
     padding: 40px;
     border-radius: 8px;
     text-align: center;
     width: 400px;
     color: #d1d7db;
 }
 
 .whatsapp-logo {
     width: 80px;
     height: 80px;
     margin: 0 auto 20px;
 }
 
 .whatsapp-logo img {
     width: 100%;
     height: 100%;
     object-fit: contain;
 }
 
 .login-container h1 {
     font-size: 20px;
     color: #e9edef;
     margin-bottom: 20px;
 }
 
 .login-container form {
     display: flex;
     flex-direction: column;
 }
 
 .login-input {
     background-color: #2a3942;
     border: none;
     color: #d1d7db;
     padding: 12px;
     border-radius: 8px;
     margin-bottom: 15px;
     font-size: 15px;
     outline: none;
 }
 
 .login-input::placeholder {
     color: #8696a0;
 }
 
 .login-button {
     background-color: #00a884;
     color: #fff;
     padding: 12px;
     border: none;
     border-radius: 8px;
     font-size: 16px;
     cursor: pointer;
     transition: background-color 0.3s ease;
 }
 
 .login-button:hover {
     background-color: #008069;
 }
 
 .error-message {
     color: #f15c6d;
     margin-top: 10px;
     font-size: 14px;
 }
 
 .disclaimer {
     margin-top: 20px;
     font-size: 12px;
     color: #8696a0;
 }
 </style>