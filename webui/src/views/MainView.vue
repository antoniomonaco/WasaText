<!-- MainView.vue -->
<template>
  <header>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
  </header>
  <div class="main-view">
    <div class="conversations-panel">
      <!-- Header della sidebar -->
      <div class="sidebar-header">
        <div class = "user-info">
          <div class="user-avatar">
            <img 
              :src="currentUser?.photoUrl || '/default-avatar.jpeg'" 
              :alt="currentUser?.username"
              class="avatar-img"
            />
          </div>
          <div class="user-name">
            {{currentUser?.username}}
          </div>
        </div>
        <div class="header-actions">
          <button class="icon-button" @click="handleNewChat">
            <i class="fas fa-comment-medical"></i>
          </button>
          <button class="icon-button" @click="showUserMenu = !showUserMenu">
            <i class="fas fa-ellipsis-v"></i>
          </button>
        </div>
        <!-- Menu utente -->
        <div v-if="showUserMenu" class="user-menu">
          <div class="menu-item" @click="handleProfileSettings">
            <i class="fas fa-user"></i>
            Profilo
          </div>
          <div class="menu-item" @click="handleLogout">
            <i class="fas fa-sign-out-alt"></i>
            Logout
          </div>
        </div>
      </div>

      <!-- Barra di ricerca -->
      <div class="search-container">
        <div class="search-box">
          <i class="fas fa-search search-icon"></i>
          <input 
            type="text" 
            placeholder="Cerca o inizia una nuova chat"
            v-model="searchQuery"
            class="search-input"
          />
        </div>
      </div>

      <ConversationsView 
        @select-conversation="selectConversation" 
        :searchQuery="searchQuery"
      />
    </div>

    <div class="chat-panel">
      <ChatView 
        v-if="selectedConversationId" 
        :conversationID="selectedConversationId" 
        @conversation-deleted="handleConversationDeleted"
      />
      <div v-else class="welcome-screen">
        <div class="welcome-content">
          <i class="fas fa-comments welcome-icon"></i>
          <h2>WASAText</h2>
          <p>Seleziona una chat per iniziare a messaggiare</p>
        </div>
      </div>
    </div>

    <!-- Modal per nuova chat/gruppo -->
    <NewChatModal 
      v-if="showNewChatModal"
      :users="users"
      @close="showNewChatModal = false"
      @create="createConversation"
    />
  </div>
</template>

<script>
import ConversationsView from './ConversationsView.vue'
import ChatView from './ChatView.vue'
import NewChatModal from '@/components/NewChatModal.vue'

export default {
  components: {
    ConversationsView,
    ChatView,
    NewChatModal
  },
  data() {
    return {
      selectedConversationId: null,
      searchQuery: '',
      showUserMenu: false,
      showNewChatModal: false,
      users: [],
      currentUser: null,
      pollingInterval : null
    }
  },
  async created() {
    await this.refreshUserData();
    this.startPolling();
    window.addEventListener('userProfileUpdated', this.refreshUserData);
  },
  beforeUnmount() {
    this.stopPolling();
    window.removeEventListener('userProfileUpdated', this.refreshUserData);
  },
  methods: {
    async refreshUserData() {
      try {
        const response = await this.$axios.get('/users/me', {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          }
        });
        this.currentUser = response.data;
      } catch (error) {
        console.error('Errore nel recupero dei dati utente:', error);
      }
    },
    async selectConversation(conversationId) {
      if (conversationId) {
        this.selectedConversationId = conversationId;
        this.showUserMenu = false;
      }
    },
    handleConversationDeleted(deletedConversationId) {
      if (this.selectedConversationId === deletedConversationId) {
        this.selectedConversationId = null
      }
    },
    async handleNewChat() {
      this.showNewChatModal = true
      await this.fetchUsers()
    },
    async fetchUsers() {
      try {
        const response = await this.$axios.get('/users/', {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          }
        })
        this.users = response.data
      } catch (error) {
        console.error('Errore nel recupero degli utenti:', error)
      }
    },
    async createConversation(data) {
      console.log('Received conversation data:', data);
      try {
        if (!data.type || !data.participants || !data.participants.length) {
          console.error('Dati conversazione non validi:', data);
          return;
        }

        const payload = {
          type: data.type,
          participants: data.participants,
        };

        if (data.type === 'group') {
          payload.name = data.name;
          if (data.photoUrl) {
            payload.photoUrl = data.photoUrl;
          }
        }

        console.log('Sending payload:', payload);
        const response = await this.$axios.post('/conversations/', payload, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          }
        });
        
        console.log('Server response:', response.data);
        
        if (response.data && response.data.conversation_id) { 
          this.showNewChatModal = false;
          await this.$nextTick();
          this.selectedConversationId = response.data.conversation_id; 
        } else {
          console.error('Response missing conversation ID:', response.data);
        }
      } catch (error) {
        console.error('Errore nella creazione della conversazione:', error.response || error);
      }
    },
    handleProfileSettings() {
      this.$router.push('/profile')
      this.showUserMenu = false
    },
    handleLogout() {
      localStorage.removeItem('authToken')
      this.$router.push('/session')
    },
    startPolling() {
      this.pollingInterval = setInterval(this.refreshUserData, 5000);
    },
    stopPolling() {
      if (this.pollingInterval) {
        clearInterval(this.pollingInterval);
      }
    }
  }
}
</script>

<style scoped>
.main-view {
  display: flex;
  height: 100vh;
  background-color: #111b21;
}

.conversations-panel {
  width: 40%;
  max-width: 400px;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #2a3942;
}

.chat-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  background-color: #202c33;
  height: 60px;
  position: relative;
}

.user-avatar {
  width: 40px;
  height: 40px;
}

.avatar-img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.icon-button {
  background: none;
  border: none;
  color: #aebac1;
  font-size: 18px;
  cursor: pointer;
  padding: 8px;
  border-radius: 50%;
}

.icon-button:hover {
  background-color: #384147;
}

.user-menu {
  position: absolute;
  top: 100%;
  right: 10px;
  background-color: #233138;
  border-radius: 3px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.3);
  z-index: 100;
}

.menu-item {
  padding: 12px 24px;
  color: #d1d7db;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 12px;
}

.menu-item:hover {
  background-color: #182229;
}

.search-container {
  padding: 8px;
  background-color: #111b21;
}

.search-box {
  position: relative;
  background-color: #202c33;
  border-radius: 8px;
  padding: 8px;
}

.search-icon {
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: #aebac1;
}

.search-input {
  width: 100%;
  background: none;
  border: none;
  color: #d1d7db;
  padding: 6px 35px;
  outline: none;
}

.search-input::placeholder {
  color: #8696a0;
}

.welcome-screen {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #222e35;
}

.welcome-content {
  text-align: center;
  color: #8696a0;
}

.welcome-icon {
  font-size: 48px;
  margin-bottom: 20px;
}

.user-name {
  color: #e9edef;
  font-size: 16px;
  font-weight: 500;
  margin-left: 10px; 
}
.user-info { 
  display: flex;
  align-items: center; 
}
</style>