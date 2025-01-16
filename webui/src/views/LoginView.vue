<template>
  <div class="login-view">
    <div class="login-container">
      <div class="whatsapp-logo">
        <img src="/logo_whatsapp.png" alt="WhatsApp Logo"> 
      </div>
      <h1>Accedi a WhatsApp</h1>
      <form @submit.prevent="login">
        <input type="text" v-model="username" placeholder="Inserisci il tuo nome" required>
        <button type="submit">Inizia a chattare</button>
      </form>
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      <div class="disclaimer">Questo non è il vero WhatsApp, è una demo a scopo didattico.</div>
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
  background-color: #f0f2f5;
}

.login-container {
  background-color: #fff;
  padding: 40px;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  text-align: center;
  width: 400px;
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
  color: #4a4a4a;
  margin-bottom: 20px;
}

.login-container form {
  display: flex;
  flex-direction: column;
}

.login-container input {
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  margin-bottom: 15px;
  font-size: 16px;
}

.login-container button {
  background-color: #00a884;
  color: #fff;
  padding: 10px;
  border: none;
  border-radius: 5px;
  font-size: 16px;
  cursor: pointer;
}

.login-container button:hover {
  background-color: #008069;
}

.error-message {
  color: red;
  margin-top: 10px;
  font-size: 14px;
}
.disclaimer{
  margin-top: 20px;
  font-size: 12px;
  color: #888;
}
</style>