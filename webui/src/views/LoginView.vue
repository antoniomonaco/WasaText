<template>
    <div>
      <h1>Login</h1>
      <form @submit.prevent="login">
        <input type="text" v-model="username" placeholder="Username">
        <button type="submit">Login</button>
      </form>
      <div v-if="errorMessage" class="error">{{ errorMessage }}</div>
    </div>
</template>
  
<script>
  export default {
    data() {
      return {
        username: '',
        errorMessage: null
      };
    },
    methods: {
      async login() {
        this.errorMessage = null;
        try {
          const response = await this.$axios.post('/session', { name: this.username });
          const token = response.data.identifier;
          localStorage.setItem('authToken', token);
  
          this.$router.push('/conversations/');
        } catch (error) {
          if (error.response) {
            this.errorMessage = error.response.data;
          } else if (error.request) {
            this.errorMessage = "Errore di rete";
          } else {
            this.errorMessage = "Errore nella richiesta";
          }
        }
      }
    }
  };
</script>