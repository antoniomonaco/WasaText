<template>
  <div class="conversations-container">
    <div class="conversations-list">
      <div 
        v-for="conversation in filteredConversations" 
        :key="conversation.id"
        class="conversation-item"
        :class="{ 'active': selectedConversationId === conversation.id }"
        @click="selectConversation(conversation.id)"
      >
        <!-- Avatar della conversazione -->
        <div class="conversation-avatar">
          <img 
            :src="getConversationPhoto(conversation)" 
            :alt="getConversationName(conversation)"
            class="avatar-img"
          />
        </div>

        <!-- Info della conversazione -->
        <div class="conversation-content">
          <div class="conversation-header">
            <span class="conversation-name">
              {{ getConversationName(conversation) }}
            </span>
            <span class="conversation-time" v-if="conversation.latestMessage">
              {{ formatConversationTime(conversation.latestMessage.timestamp) }}
            </span>
          </div>

          <div class="conversation-preview">
            <!-- Anteprima ultimo messaggio -->
            <div class="message-preview" :class="{'unread': hasUnreadMessages(conversation)}">
              <span v-if="conversation.latestMessage && conversation.latestMessage.sender.id === currentUserId">
                <i class="fas fa-check message-status" 
                   :class="{'read': conversation.latestMessage.status === 'read'}">
                </i>
              </span>
              {{ getMessagePreview(conversation) }}
            </div>
            
            <!-- Badge messaggi non letti -->
            <div v-if="getUnreadCount(conversation) > 0" class="unread-badge">
              {{ getUnreadCount(conversation) }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    searchQuery: {
      type: String,
      default: ''
    }
  },
  
  data() {
    return {
      conversations: [],
      selectedConversationId: null,
      currentUserId: parseInt(localStorage.getItem('authToken')),
      pollingInterval: null
    }
  },

  computed: {
    filteredConversations() {
      if (!this.searchQuery) return this.conversations;

      const query = this.searchQuery.toLowerCase();
      return this.conversations.filter(conversation => {
        const conversationName = this.getConversationName(conversation).toLowerCase();
        const lastMessage = this.getMessagePreview(conversation).toLowerCase();
        return conversationName.includes(query) || lastMessage.includes(query);
      });
    }
  },

  methods: {
    getConversationName(conversation) {
      if (conversation.type === 'group') {
        return conversation.name || 'Gruppo senza nome';
      }
      const otherParticipant = this.getOtherParticipant(conversation);
      return otherParticipant ? otherParticipant.username : 'Utente sconosciuto';
    },

    getConversationPhoto(conversation) {
      if (conversation.type === 'group') {
        return conversation.photoUrl || '/default-avatar.jpeg';
      }
      const otherParticipant = this.getOtherParticipant(conversation);
      return otherParticipant?.photoUrl || '/default-avatar.jpeg';
    },

    getOtherParticipant(conversation) {
      return conversation.participants?.find(p => p.id !== this.currentUserId);
    },

    getMessagePreview(conversation) {
      if (!conversation.latestMessage) return 'Nessun messaggio';
      
      const message = conversation.latestMessage;
      if (message.type === 'text') {
        return message.content;
      } else if (message.type === 'media') {
        return 'Foto';
      }
      return 'Messaggio';
    },

    formatConversationTime(timestamp) {
      if (!timestamp) return '';
      
      const date = new Date(timestamp);
      
      // Verifica se la data è valida
      if (isNaN(date.getTime())) {
        console.error('Invalid date:', timestamp);
        return '';
      }

      const now = new Date();
      const yesterday = new Date(now);
      yesterday.setDate(yesterday.getDate() - 1);

      // Se è oggi
      if (date.toDateString() === now.toDateString()) {
        return date.toLocaleTimeString([], { 
          hour: '2-digit', 
          minute: '2-digit' 
        });
      }
      
      // Se è ieri
      if (date.toDateString() === yesterday.toDateString()) {
        return 'Ieri';
      }
      
      // Se è quest'anno
      if (date.getFullYear() === now.getFullYear()) {
        return date.toLocaleDateString([], {
          day: '2-digit',
          month: '2-digit'
        });
      }
      
      // Se è un anno diverso
      return date.toLocaleDateString([], {
        day: '2-digit',
        month: '2-digit',
        year: '2-digit'
      });
    },

    hasUnreadMessages(conversation) {
      return conversation.messages?.some(message => 
        message.sender.id !== this.currentUserId && 
        message.status !== 'read'
      );
    },

    getUnreadCount(conversation) {
      return conversation.messages?.filter(message => 
        message.sender.id !== this.currentUserId && 
        message.status !== 'read'
      ).length || 0;
    },

    async fetchConversations() {
      console.log("Auth code: ",localStorage.getItem('authToken'))
      try {
        const response = await this.$axios.get('/conversations/', {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          }
        });
        this.conversations = response.data;
      } catch (error) {
        console.error('Errore nel recupero delle conversazioni:', error);
      }
    },

    selectConversation(id) {
      this.selectedConversationId = id;
      this.$emit('select-conversation', id);
    },

    startPolling() {
      this.pollingInterval = setInterval(this.fetchConversations, 5000);
    },

    stopPolling() {
      if (this.pollingInterval) {
        clearInterval(this.pollingInterval);
        this.pollingInterval = null;
      }
    }
  },

  created() {
    this.fetchConversations();
    this.startPolling();
    window.addEventListener('userProfileUpdated', this.fetchConversations);
  },

  beforeUnmount() {
    this.stopPolling();
    window.removeEventListener('userProfileUpdated', this.fetchConversations);
  }
}
</script>

<style scoped>
.conversations-container {
  flex: 1;
  overflow-y: auto;
  background-color: #111b21;
}

.conversations-list {
  display: flex;
  flex-direction: column;
}

.conversation-item {
  display: flex;
  padding: 12px 16px;
  cursor: pointer;
  transition: background-color 0.2s;
  border-bottom: 1px solid #222d34;
}

.conversation-item:hover {
  background-color: #202c33;
}

.conversation-item.active {
  background-color: #2a3942;
}

.conversation-avatar {
  width: 49px;
  height: 49px;
  margin-right: 15px;
  flex-shrink: 0;
}

.avatar-img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.conversation-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.conversation-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 4px;
}

.conversation-name {
  color: #e9edef;
  font-size: 17px;
  font-weight: 400;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.conversation-time {
  color: #8696a0;
  font-size: 12px;
  white-space: nowrap;
  margin-left: 6px;
}

.conversation-preview {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.message-preview {
  color: #8696a0;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 85%;
  display: flex;
  align-items: center;
  gap: 4px;
}

.message-preview.unread {
  color: #e9edef;
  font-weight: 500;
}

.message-status {
  font-size: 14px;
  color: #8696a0;
}

.message-status.read {
  color: #53bdeb;
}

.unread-badge {
  background-color: #00a884;
  color: #111b21;
  font-size: 12px;
  font-weight: 500;
  min-width: 20px;
  height: 20px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 6px;
}
</style>