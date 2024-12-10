<template>
  <div class="chat-view">
    <div class="chat-header" v-if="conversation">
      <img v-if="conversation.photoUrl" :src="conversation.photoUrl" alt="Foto profilo" class="chat-photo" />
      <div class="chat-name">
        <span v-if="conversation.type === 'chat'">
          {{ conversation.name || (conversation.participants && getOtherParticipant(conversation).username) }} 
        </span>
        <span v-else-if="conversation.type === 'group'">
          {{ conversation.name }}
        </span>
      </div>
    </div>
    <div class="chat-messages" ref="chatMessages">
      <div v-for="message in messages" :key="message.id" class="message" :class="{ 'my-message': message.sender.id === currentUserID }">
        <div class="message-header">
          <span class="message-sender">{{ message.sender.username }}</span>
          <span class="message-timestamp">{{ formatTimestamp(message.timestamp) }}</span>
        </div>
        <div class="message-content">{{ message.content }}</div>
      </div>
    </div>
    <div class="chat-input">
      <input type="text" v-model="newMessage" @keyup.enter="sendMessage" placeholder="Scrivi un messaggio..." />
      <button @click="sendMessage">Invia</button>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      conversation: {},
      messages: [],
      newMessage: '',
      currentUserID: this.getCurrentUserID()
    };
  },
  async created() {
    await this.fetchConversation();
  },
  methods: {
    getCurrentUserID() {
      const token = localStorage.getItem('authToken');
      return parseInt(token);
    },
    async fetchConversation() {
      try {
        const conversationID = this.$route.params.conversationID
        const token = localStorage.getItem('authToken');
        const response = await this.$axios.get(`/conversations/${conversationID}`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        this.conversation = response.data.conversation; 
        this.messages = response.data.messages;
      } catch (error) {
        console.error("Errore nel recupero della conversazione:", error);
      }
    },
    getOtherParticipant(conversation) {
      console.log(this.conversation)
      const currentUserID = this.currentUserID;
      return conversation.participants.find(p => p.id !== currentUserID);
    },
    formatTimestamp(timestamp) {
      return new Date(timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    },
    async sendMessage() {
      try {
        const conversationID = this.$route.params.conversationID;
        const token = localStorage.getItem('authToken');
        await this.$axios.post(`/conversations/${conversationID}`, {
          sender: this.currentUserID,
          type: 'text',
          content: this.newMessage
        }, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        this.newMessage = '';
        await this.fetchMessages(); // Aggiorno i messaggi dopo l'invio
      } catch (error) {
        console.error("Errore nell'invio del messaggio:", error);
      }
    }
  }
};
</script>

<style scoped>
.chat-view {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.chat-header {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #eee;
}

.chat-photo {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  margin-right: 10px;
}

.chat-name {
  font-weight: bold;
}

.chat-messages {
  flex-grow: 1;
  overflow-y: auto;
  padding: 10px;
}

.message {
  display: flex;
  flex-direction: column;
  margin-bottom: 10px;
}

.my-message {
  align-self: flex-end;
}

.message-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 5px;
}

.message-sender {
  font-weight: bold;
}

.message-timestamp {
  font-size: 0.8em;
  color: #666;
}

.message-content {
  background-color: #eee;
  padding: 10px;
  border-radius: 5px;
}

.chat-input {
  display: flex;
  padding: 10px;
  background-color: #eee;
}

.chat-input input {
  flex-grow: 1;
  margin-right: 10px;
  padding: 5px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.chat-input button {
  padding: 5px 10px;
  background-color: #0084ff;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}
</style>