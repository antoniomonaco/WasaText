<template>
    <div class="conversations-view">
      <h2>Conversazioni</h2>
      <ul v-if="conversations.length > 0">
        <li v-for="conversation in conversations" :key="conversation.id" class="conversation">
          <router-link :to="`/conversations/${conversation.id}`" class="conversation-link">
            <img v-if="conversation.photoUrl" :src="conversation.photoUrl" alt="Foto profilo" class="conversation-photo" />
            <div class="conversation-info">
              <div class="conversation-name">
                <span v-if="conversation.type === 'group'">
                  {{ conversation.name || "Gruppo senza nome" }}
                </span>
                <span v-else>
                  {{ getOtherParticipant(conversation).username }}
                </span>
              </div>
              <div class="conversation-preview">
                {{ conversation.latestMessage ? conversation.latestMessage.content : "Nessun messaggio" }}
              </div>
            </div>
          </router-link>
        </li>
      </ul>
      <div v-else>Nessuna conversazione trovata.</div>
    </div>
</template>
  
<script>
  export default {
    data() {
      return {
        conversations: []
      };
    },
    async created() {
      await this.fetchConversations();
    },
    methods: {
      async fetchConversations() {
        console.log("fetch conversations")
        try {
          const token = localStorage.getItem('authToken');
          let authorization = `Bearer ${token}`
          console.log(authorization)
          const response = await this.$axios.get('/conversations/', {
            headers: {
              Authorization: `Bearer ${token}`
            }
          });
          this.conversations = response.data;
        } catch (error) {
          console.error("Errore nel recupero delle conversazioni:", error);
        }
      },
      getOtherParticipant(conversation) {
        const token = localStorage.getItem('authToken');
        const currentUserID = parseInt(token);

        return conversation.participants.find(p => p.id !== currentUserID);
    }
    }
  };
</script>

<style scoped>
/*
.conversations-view {
  // Stile per il container principale 
}

*/
.conversation {
  /* Stile per ogni conversazione */
  display: flex;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #ccc;
}

.conversation-link {
  /* Stile per il link alla conversazione */
  text-decoration: none;
  color: inherit;
  display: flex;
  width: 100%;
}

.conversation-photo {
  /* Stile per la foto profilo */
  width: 50px;
  height: 50px;
  border-radius: 50%;
  margin-right: 10px;
}

.conversation-info {
  /* Stile per il nome e l'anteprima */
  display: flex;
  flex-direction: column;
}

.conversation-name {
  /* Stile per il nome della conversazione */
  font-weight: bold;
}

.conversation-preview {
  /* Stile per l'anteprima del messaggio */
  font-size: 0.9em;
  color: #666;
}
</style>