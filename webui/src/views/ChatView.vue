<template>
  <div class="chat-view">
    <div class="chat-header" v-if="conversation">
      <img
        v-if="conversation.photoUrl"
        :src="conversation.photoUrl"
        alt="Foto profilo"
        class="chat-photo"
      />
      <div class="chat-name-container">
        <div class="chat-name">
          <span v-if="conversation.type === 'chat'">
            {{
              conversation.name ||
                (conversation.participants &&
                  getOtherParticipant(conversation).username)
            }}
          </span>
          <span v-else-if="conversation.type === 'group'">
            {{ conversation.name }}
          </span>
        </div>
        <div class="chat-status">
          <span v-if="conversation.type === 'group'">
            {{ conversation.participants.length }} partecipanti
          </span>
        </div>
      </div>
      <div class="chat-header-icons">
          <i class="fas fa-search" @click="toggleSearch"></i>
          <i class="fas fa-ellipsis-v" @click="toggleMenu"></i>
      </div>
      <div class="search-bar" v-if="showSearch">
        <i class="fas fa-arrow-left search-back-icon" @click="toggleSearch"></i>
        <input type="text" placeholder="Cerca..." v-model="searchQuery" />
      </div>
    </div>
    <div class="chat-messages" ref="chatMessages">
      <div
        v-for="message in filteredMessages"
        :key="message.id"
        class="message"
        :class="{ 'my-message': message.sender.id === currentUserID }"
      >
        <div class="message-content" @click="toggleMessageMenu(message.id)">
          <div v-if="message.replyToMessage" class="reply-message" :class="{'my-reply': getMessageById(message.replyToMessage).sender.id === currentUserID}">
            <span class="reply-sender">
              {{ getMessageById(message.replyToMessage) ? getMessageById(message.replyToMessage).sender.username : 'Messaggio non trovato' }}
            </span>
            <span class="reply-content">
                {{ getReplyToMessageContent(message.replyToMessage) }}
            </span>
          </div>
          <span class="message-sender">
            {{ message.sender.id === currentUserID ? "" : message.sender.username }}
          </span>
          {{ message.content }}
          <span class="message-timestamp">
            {{ formatTimestamp(message.timestamp) }}
            <i v-if="message.sender.id === currentUserID" :class="[message.readStatus ? 'fas fa-check-double' : 'fas fa-check', {'read': message.readStatus}]"></i>
          </span>
          <div v-if="showCommentsForMessage === message.id" class="comments-section">
            <div v-for="comment in message.comments" :key="comment.id" class="comment">
              <span class="comment-sender">{{ comment.sender.username }}: </span>
              <span class="comment-content">{{ comment.content }}</span>
            </div>
            <div class="add-comment-section">
              <input
                type="text"
                v-model="newComment"
                placeholder="Aggiungi un commento..."
                @keyup.enter="addComment(message.id)"
              />
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
      <div class="reply-preview" v-if="replyTo">
        <div class="reply-to-content">
            <span class="replying-to-label">Rispondendo a: </span>
            <span class="reply-to-user">{{ replyTo.sender.username }}</span>
            <span class="reply-to-message">{{ replyTo.content }}</span>
        </div>
        <button class="close-reply" @click="cancelReply">X</button>
      </div>
      <i class="fas fa-smile emoji-icon"></i>
      <input
        type="text"
        v-model="newMessage"
        @keyup.enter="sendMessage"
        :placeholder="replyTo ? 'Rispondi al messaggio...' : 'Scrivi un messaggio...'"
      />
      <i class="fas fa-paperclip attachment-icon"></i>
      <button class="send-button" @click="sendMessage">
        <i class="fas fa-paper-plane"></i>
      </button>
    </div>
    <div v-if="showForwardingModal" class="forward-modal">
      <div class="forward-modal-content">
        <h3>Scegli una conversazione:</h3>
        <ul>
          <li
            v-for="conversation in otherConversations"
            :key="conversation.id"
            @click="confirmForward(conversation.id)"
          >
            <div class="forward-chat-photo">
              <img
                v-if="conversation.photoUrl"
                :src="conversation.photoUrl"
                alt="Foto profilo"
                class="chat-photo"
              />
            </div>
            <div class="forward-chat-name">
              {{
                conversation.type === "group"
                  ? conversation.name
                  : getOtherParticipant(conversation).username
              }}
            </div>
          </li>
        </ul>
        <button @click="showForwardingModal = false">Annulla</button>
      </div>
    </div>
    <div class="conversation-menu" v-if="showConversationMenu">
        <button @click="deleteConversation">Elimina Conversazione</button>
        <button @click="toggleMenu">Chiudi</button>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    conversationID: {
      type: Number,
      required: true,
    },
  },
  data() {
    return {
      conversation: {},
      messages: [],
      newMessage: "",
      currentUserID: this.getCurrentUserID(),
      showMenuForMessage: null,
      replyTo: null, // Aggiunto per gestire le risposte
      forwardingMessage: null,
      showForwardingModal: false,
      otherConversations: [],
      newComment: "",
      showCommentsForMessage: null,
      searchQuery: "",
      showSearch: false,
      showConversationMenu: false,
    };
  },
  computed: {
    filteredMessages() {
      if (!this.searchQuery) {
        return this.messages;
      }
      const query = this.searchQuery.toLowerCase();
      return this.messages.filter(message => {
        const messageContent = (message.content || "").toLowerCase();
        const senderUsername = (message.sender.username || "").toLowerCase();
        // Cerca la corrispondenza nel contenuto del messaggio o nel nome del mittente
        return messageContent.includes(query) || senderUsername.includes(query);
      });
    }
  },
  watch: {
    conversationID: {
      handler: function (newConversationId, oldConversationId) {
        if (newConversationId) {
          this.fetchConversation();
        }
      },
      immediate: true,
    },
  },
  methods: {
    getCurrentUserID() {
      const token = localStorage.getItem("authToken");
      return parseInt(token);
    },
    async fetchConversation() {
      try {
        // const conversationID = this.$route.params.conversationID // Non serve più
        const conversationID = this.conversationID; // Usa la prop
        const token = localStorage.getItem("authToken");
        const response = await this.$axios.get(
          `/conversations/${conversationID}`,
          {
            headers: {
              Authorization: `Bearer ${token}`,
            },
          }
        );
        this.conversation = response.data.conversation;
        this.messages = response.data.messages;
        this.scrollToBottom();
        this.markMessagesAsRead(conversationID);
      } catch (error) {
        console.error("Errore nel recupero della conversazione:", error);
      }
    },
    getOtherParticipant(conversation) {
      const currentUserID = this.getCurrentUserID();
      return conversation.participants.find((p) => p.id !== currentUserID);
    },
    formatTimestamp(timestamp) {
      return new Date(timestamp).toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });
    },
    getReplyToMessageContent(replyToMessageId) {
      let replyToMessage = this.messages.find((m) => m.id === replyToMessageId);
      return replyToMessage ? replyToMessage.content : "Messaggio non trovato";
    },
    toggleMessageMenu(messageId) {
      this.showMenuForMessage = this.showMenuForMessage === messageId ? null : messageId;
      this.showCommentsForMessage = null; // Chiudi i commenti se aperti
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
        // Escludi la conversazione corrente
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

        // Aggiungi il nuovo commento alla lista dei commenti del messaggio
        const message = this.messages.find(m => m.id === messageId);
        if (message) {
          if (!message.comments) {
            message.comments = [];
          }
          message.comments.push(response.data);
        }

        this.newComment = ''; // Resetta l'input del commento
        this.showCommentsForMessage = messageId; // Mantieni aperta la sezione commenti
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
      this.$nextTick(() => {
        if (this.$refs.chatMessages) {
          this.$refs.chatMessages.scrollTop = this.$refs.chatMessages.scrollHeight;
        }
      });
    },
    toggleSearch() {
        this.showSearch = !this.showSearch;
        if (!this.showSearch) {
            this.searchQuery = ''; // Resetta la query di ricerca quando la barra di ricerca viene chiusa
        }
    },
    toggleMenu() {
        this.showConversationMenu = !this.showConversationMenu;
    },
    async deleteConversation() {
        try {
            const token = localStorage.getItem('authToken');
            await this.$axios.delete(`/conversations/${this.conversationID}`, {
            headers: {
                Authorization: `Bearer ${token}`
            }
            });
            // Resetta la conversazione e i messaggi
            this.conversation = {};
            this.messages = [];
            this.showConversationMenu = false;
            // Emetti un evento per notificare che la conversazione è stata eliminata
            this.$emit('conversation-deleted', this.conversationID);
        } catch (error) {
            console.error("Errore nell'eliminazione della conversazione:", error);
        }
    },
    async markMessagesAsRead(conversationId) {
      try {
        const token = localStorage.getItem('authToken');
        await this.$axios.put(`/conversations/${conversationId}`,
          {},
          {
            headers: {
              Authorization: `Bearer ${token}`
            }
          });
        // Aggiorna localmente lo stato dei messaggi come letti
        this.messages.forEach(message => {
            if (message.sender.id !== this.currentUserID) {
                message.readStatus = true;
            }
        });
      } catch (error) {
        console.error("Errore nell'aggiornamento dello stato di lettura dei messaggi:", error);
      }
    },
  }
};
</script>

<style scoped>
@import url("https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.2/css/all.min.css");
.chat-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #efeae2;
}
.chat-header {
  position: sticky;
  top: 0;
  display: flex;
  align-items: center;
  padding: 10px 15px;
  background-color: #f0f2f5;
  border-bottom: 1px solid #d1d7db;
  z-index: 10;
}
.chat-photo {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 15px;
}
.chat-name-container {
  flex-grow: 1;
}
.chat-name {
  font-weight: 600;
  color: #111b21;
}
.chat-status {
  font-size: 12px;
  color: #667781;
}
.chat-header-icons {
  display: flex;
  align-items: center;
}
.chat-header-icons i {
  margin-left: 20px;
  color: #54656f;
  cursor: pointer;
}
.search-bar {
  display: flex;
  align-items: center;
  background-color: #fff;
  padding: 8px;
  border-radius: 20px;
  flex-grow: 1;
  margin-left: 15px;
}
.search-back-icon {
  margin-right: 10px;
  color: #919191;
  cursor: pointer;
}
.search-bar input {
  border: none;
  outline: none;
  flex-grow: 1;
  font-size: 14px;
}
.chat-messages {
  flex-grow: 1;
  overflow-y: auto;
  padding: 10px;
}
.message {
  display: flex;
  margin-bottom: 5px;
  flex-direction: column;
  align-items: flex-start;
}
.my-message {
  align-items: flex-end;
}
.message-content {
  background-color: #fff;
  padding: 8px 12px;
  border-radius: 7.5px;
  box-shadow: 0 1px 0.5px rgba(0, 0, 0, 0.13);
  max-width: 80%;
  position: relative;
}
.my-message .message-content {
  background-color: #d9fdd3;
}
.message-sender {
  font-weight: bold;
  color: #008069;
  font-size: 12px;
  margin-bottom: 2px;
  display: inline-block;
}
.message-timestamp {
  font-size: 11px;
  color: #888;
  margin-left: 5px;
  display: flex;
  align-items: center;
}
.message-timestamp .fas {
  margin-left: 3px;
}
.message-timestamp .fa-check {
  color: #888; /* Colore per il segno di spunta non letto */
}

.message-timestamp .fa-check-double {
  color: #888; /* Colore per il segno di spunta non letto */
}

.message-timestamp .read {
  color: #34b7f1; /* Colore per il segno di spunta blu (letto) */
}
.reply-message{
    background-color: #f0f0f0;
    border: 1px solid #ccc;
    border-radius: 5px;
    padding: 5px;
    margin-bottom: 5px;
    font-size: 14px;
}
.my-reply{
    background-color: #c3e8bd;
}
.reply-sender{
    font-weight: bold;
    color: #008069;
    font-size: 12px;
    margin-bottom: 2px;
    display: block;
}
.reply-content{
    color: #333;
    margin-bottom: 2px;
    display: block;
}
.comments-section {
  margin-top: 10px;
  border-top: 1px solid #ccc;
  padding-top: 5px;
}
.comment {
  background-color: #f8f9fa;
  border: 1px solid #ccc;
  border-radius: 5px;
  padding: 5px;
  margin-bottom: 5px;
  font-size: 13px;
}
.comment-sender {
  font-weight: bold;
  color: #008069;
}
.comment-content {
  color: #333;
}
.add-comment-section {
  margin-top: 5px;
  display: flex;
  align-items: center;
}
.add-comment-section input {
  flex-grow: 1;
  padding: 5px;
  border: 1px solid #ccc;
  border-radius: 5px;
  margin-right: 5px;
  font-size: 13px;
}
.add-comment-section button {
  padding: 5px 10px;
  background-color: #00a884;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 13px;
}
.chat-input {
  position: sticky;
  bottom: 0;
  background-color: #f0f2f5;
  padding: 10px;
  display: flex;
  align-items: center;
  border-top: 1px solid #d1d7db;
}
.emoji-icon, .attachment-icon {
  margin-right: 15px;
  color: #54656f;
  cursor: pointer;
  font-size: 20px;
}
.chat-input input {
  flex-grow: 1;
  border: none;
  background-color: #fff;
  padding: 10px;
  border-radius: 20px;
  margin-right: 10px;
  font-size: 14px;
  outline: none;
}
.send-button {
  background-color: #00a884;
  border: none;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}
.send-button i {
  color: white;
  font-size: 18px;
}
.message-menu {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -100%);
  background-color: white;
  border: 1px solid #ccc;
  border-radius: 5px;
  z-index: 10;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
  padding: 5px 0;
}

.message-menu button {
  display: block;
  width: 100%;
  padding: 5px 10px;
  background-color: transparent;
  border: none;
  text-align: left;
  font-size: 14px;
  cursor: pointer;
}

.message-menu button:hover {
  background-color: #f0f0f0;
}
.reply-preview {
  display: flex;
  align-items: center;
  background-color: #e9e9eb;
  padding: 5px;
  margin-bottom: 5px;
  border-radius: 5px;
}

.reply-to-content {
  flex-grow: 1;
}

.replying-to-label {
  font-size: 12px;
  color: #667781;
}

.reply-to-user {
  font-weight: bold;
  color: #111b21;
}

.reply-to-message {
  font-size: 14px;
  color: #333;
}

.close-reply {
  background: none;
  border: none;
  padding: 0;
  font-size: 14px;
  color: #667781;
  cursor: pointer;
  margin-left: 10px;
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
.conversation-menu {
  position: absolute;
  top: 50px;
  right: 10px;
  background-color: white;
  border: 1px solid #ccc;
  border-radius: 5px;
  z-index: 10;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
  padding: 5px 0;
}

.conversation-menu button {
  display: block;
  width: 100%;
  padding: 5px 10px;
  background-color: transparent;
  border: none;
  text-align: left;
  font-size: 14px;
  cursor: pointer;
}

.conversation-menu button:hover {
  background-color: #f0f0f0;
}
</style>