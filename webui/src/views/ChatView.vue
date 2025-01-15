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
      <div v-for="message in messages" :key="message.id" class="message" :class="{ 'my-message': message.sender.id === currentUserID }" @click="toggleMessageMenu(message.id)">
        <div class="message-header">
          <span class="message-sender">{{ message.sender.username }}</span>
          <span class="message-timestamp">{{ formatTimestamp(message.timestamp) }}</span>
        </div>
        <div class="message-content">
          <span v-if="message.replyToMessage" class="reply-indicator">
            In risposta a: {{ getReplyToMessageContent(message.replyToMessage) }}
          </span>
          {{ message.content }}
          <div v-if="showCommentsForMessage === message.id">
            <div v-for="comment in message.comments" :key="comment.id" class="comment">
              <span class="comment-sender">{{ comment.sender.username }}: </span>
              <span class="comment-content">{{ comment.content }}</span>
            </div>
            <div class="add-comment-section">
              <input type="text" v-model="newComment" placeholder="Aggiungi un commento..." @keyup.enter="addComment(message.id)" />
              <button @click="addComment(message.id)">Invia</button>
            </div>
          </div>
        </div>
        <div class="message-menu" v-if="showMenuForMessage === message.id">
          <button @click="deleteMessage(message.id)">Elimina</button>
          <button @click="replyToMessage(message)">Rispondi</button>
          <button @click="forwardMessage(message)">Inoltra</button>
          <button @click="commentMessage(message)">Commenta</button>
        </div>
      </div>
    </div>
    <div class="chat-input">
      <input type="text" v-model="newMessage" @keyup.enter="sendMessage" :placeholder="replyTo ? 'Rispondi al messaggio...' : 'Scrivi un messaggio...'" />
      <button @click="sendMessage">Invia</button>
      <span v-if="replyTo" class="replying-to">
        Rispondendo a: {{ replyTo.sender.username }} - {{ replyTo.content }}
        <button @click="cancelReply">X</button>
      </span>
    </div>
    <div v-if="showForwardingModal" class="forward-modal">
      <div class="forward-modal-content">
        <h3>Scegli una conversazione:</h3>
        <ul>
          <li v-for="conversation in otherConversations" :key="conversation.id" @click="confirmForward(conversation.id)">
            <div class="forward-chat-photo">
              <img v-if="conversation.photoUrl" :src="conversation.photoUrl" alt="Foto profilo" class="chat-photo" />
            </div>
            <div class="forward-chat-name">
              {{ conversation.type === 'group' ? conversation.name : getOtherParticipant(conversation).username }}
            </div>
          </li>
        </ul>
        <button @click="showForwardingModal = false">Annulla</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    conversationID: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      conversation: {},
      messages: [],
      newMessage: '',
      currentUserID: this.getCurrentUserID(),
      showMenuForMessage: null,
      replyTo: null, 
      forwardingMessage: null,
      showForwardingModal: false,
      otherConversations: [],
      newComment: '',
      showCommentsForMessage : null
    };
  },
  watch: {
    conversationID: {
        handler : function(newConversationId, oldConversationId) {
            if (newConversationId) {
                this.fetchConversation();
            }
        },
        immediate: true
    }
  },
  methods: {
    getCurrentUserID() {
      const token = localStorage.getItem('authToken');
      return parseInt(token);
    },
    async fetchConversation() {
      try {
        const conversationID = this.conversationID; 
        const token = localStorage.getItem('authToken');
        const response = await this.$axios.get(`/conversations/${conversationID}`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        this.conversation = response.data.conversation;
        this.messages = response.data.messages;
        this.scrollToBottom();
      } catch (error) {
        console.error("Errore nel recupero della conversazione:", error);
      }
    },
    getOtherParticipant(conversation) {
      const currentUserID = this.getCurrentUserID();
      return conversation.participants.find(p => p.id !== currentUserID);
    },
    formatTimestamp(timestamp) {
      return new Date(timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    },
    getReplyToMessageContent(replyToMessageId){
      let replyToMessage = this.messages.find(m => m.id === replyToMessageId);
      return replyToMessage ? replyToMessage.content : 'Messaggio non trovato';
    },
    toggleMessageMenu(messageId) {
      this.showMenuForMessage = this.showMenuForMessage === messageId ? null : messageId;
      this.showCommentsForMessage = null;
    },
    async deleteMessage(messageId) {
      try {
        const conversationID = this.conversationID;
        const token = localStorage.getItem('authToken');
        await this.$axios.delete(`/conversations/${conversationID}/messages/${messageId}`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        this.showMenuForMessage = null;
        await this.fetchConversation(); // Aggiorna la conversazione dopo l'eliminazione
      } catch (error) {
        console.error("Errore nell'eliminazione del messaggio:", error);
      }
    },
    replyToMessage(message) {
      this.replyTo = message;
      this.newMessage = ""; 
      this.showMenuForMessage = null;

    },
    cancelReply() {
      this.replyTo = null;
      this.newMessage = '';
    },
    async forwardMessage(message) {
      this.forwardingMessage = message;
      this.showForwardingModal = true;
      this.showMenuForMessage = null;
      // Carica le conversazioni disponibili per l'inoltro
      await this.fetchOtherConversations();
    },
    async fetchOtherConversations() {
      try {
        const token = localStorage.getItem('authToken');
        const response = await this.$axios.get('/conversations/', {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        // Escludo la conversazione corrente
        this.otherConversations = response.data.filter(c => c.id !== this.conversation.id);
      } catch (error) {
        console.error("Errore nel recupero delle conversazioni:", error);
      }
    },
    async confirmForward(targetConversationId) {
      try {
        const token = localStorage.getItem('authToken');
        const currentConversationId = this.conversationID;
        console.log("Inoltro messaggio", this.forwardingMessage.id, "dalla conversazione", currentConversationId, "alla conversazione", targetConversationId, "con bearer token",token); 

        await this.$axios.post(
          `/conversations/${currentConversationId}/messages/${this.forwardingMessage.id}`,
          {
            id: targetConversationId 
          },
          {
            headers: {
              Authorization: `Bearer ${token}`
            }
          }
        );
        this.showForwardingModal = false;
        this.forwardingMessage = null;
        alert('Messaggio inoltrato con successo!');
      } catch (error) {
        console.error("Errore nell'inoltro del messaggio:", error);
        if (error.response && error.response.status === 403) {
          alert("Non hai i permessi per inoltrare messaggi a questa conversazione.");
        }
      }
    },
    commentMessage(message) {
      this.showCommentsForMessage = this.showCommentsForMessage === message.id ? null : message.id;
      this.showMenuForMessage = null; 
      if (this.showCommentsForMessage === message.id) {
        this.fetchComments(message.id);
      }
    },

    async addComment(messageId) {
      try {
        const token = localStorage.getItem('authToken');
        const response = await this.$axios.post(
          `/conversations/<span class="math-inline">\{this\.conversationID\}/messages/</span>{messageId}/comments`,
          {
            content: this.newComment
          },
          {
            headers: {
              Authorization: `Bearer ${token}`
            }
          }
        );

        const message = this.messages.find(m => m.id === messageId);
        if (message) {
          if (!message.comments) {
            message.comments = [];
          }
          message.comments.push(response.data);
        }

        this.newComment = ''; 
        this.showCommentsForMessage = messageId;
      } catch (error) {
        console.error("Errore nell'aggiunta del commento:", error);
      }
    },

    async fetchComments(messageId) {
      try {
        const token = localStorage.getItem('authToken');
        const response = await this.$axios.get(`/conversations/<span class="math-inline">\{this\.conversationID\}/messages/</span>{messageId}/comments`, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });

        // Trova il messaggio corrispondente e aggiorna i suoi commenti
        const message = this.messages.find(m => m.id === messageId);
        if (message) {
          message.comments = response.data;
        }
      } catch (error) {
        console.error("Errore nel recupero dei commenti:", error);
      }
    },
    getMessageById(messageId) {
      return this.messages.find(m => m.id === messageId);
    },
    async sendMessage() {
      try {
        const token = localStorage.getItem('authToken');
        let messageToSend = {
          sender: this.currentUserID,
          type: 'text',
          content: this.newMessage
        };
        if (this.replyTo) {
          messageToSend.replyToMessage = this.replyTo.id;
        }

        await this.$axios.post(`/conversations/${this.conversationID}`, messageToSend, {
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        this.newMessage = '';
        this.replyTo = null; // Resetta lo stato di risposta
        await this.fetchConversation();
        this.scrollToBottom();
      } catch (error) {
        console.error("Errore nell'invio del messaggio:", error);
      }
    },
    scrollToBottom() {
      if (this.$refs.chatMessages) {
        this.$refs.chatMessages.scrollTop = this.$refs.chatMessages.scrollHeight;
      }
    },
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
  position : relative;
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
  position: sticky;
  bottom: 0; /* Fissa la barra in basso */
  width: 100%; /* Assicura che occupi l'intera larghezza */
  padding: 10px;
  background-color: #eee;
  border-top: 1px solid #ccc;
  z-index: 10;
}

.chat-input input {
  flex-grow: 1;
  width:90%;
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
.message-menu {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  background-color: white;
  border: 1px solid #ccc;
  border-radius: 5px;
  z-index: 10;
}

.message-menu button {
  padding: 5px 10px;
  margin: 5px;
  background-color: #f8f8f8;
  border: 1px solid #ccc;
  border-radius: 5px;
  cursor: pointer;
}

.reply-indicator {
  font-size: 0.9em;
  color: #666;
  margin-bottom: 5px;
}
.comment {
  background-color: #f0f0f0;
  border: 1px solid #ccc;
  border-radius: 5px;
  padding: 5px;
  margin-top: 5px;
}

.comment-sender {
  font-weight: bold;
}

.comment-content {
  font-size: 0.9em;
}
.add-comment-section {
  margin-top: 10px;
}

.add-comment-section input {
  width: 70%;
  padding: 5px;
  border: 1px solid #ccc;
  border-radius: 5px;
  margin-right: 5px;
}

.add-comment-section button {
  padding: 5px 10px;
  background-color: #0084ff;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}
.replying-to {
  display: block;
  margin-top: 5px;
  font-size: 0.9em;
  color: #666;
}

.replying-to button {
  background: none;
  border: none;
  padding: 0;
  font-size: 0.9em;
  color: #0084ff;
  cursor: pointer;
  margin-left: 5px;
}
.forward-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 100; 
}

.forward-modal-content {
  background-color: white;
  border-radius: 10px;
  padding: 20px;
  width: 300px;
}

.forward-modal-content h3 {
  margin-top: 0;
}

.forward-modal-content ul {
  list-style: none;
  padding: 0;
}

.forward-modal-content li {
  padding: 10px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  display : flex;
}

.forward-modal-content li:hover {
  background-color: #f0f0f0;
}

.forward-modal-content button {
  margin-top: 10px;
  padding: 5px 10px;
  background-color: #ccc;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}
.forward-chat-photo {
  width: 50px;
  height: 50px;
  display : flex;
  align-items: center;
  margin-right: 10px;
}
.chat-photo {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  margin-right: 0px;
}
.forward-chat-name{
  display : flex;
  align-items: center;
}
</style>