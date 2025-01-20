<!-- ForwardModal.vue -->
<template>
    <div class="forward-modal">
      <div class="modal-header">
        <h3>Inoltra messaggio</h3>
        <button 
          class="close-button" 
          @click="$emit('close')"
        >
          <i class="fas fa-times"></i>
        </button>
      </div>
      
      <div class="modal-body">
        <!-- Ricerca conversazioni -->
        <div class="search-container">
          <div class="search-wrapper">
            <i class="fas fa-search search-icon"></i>
            <input 
              type="text" 
              v-model="searchQuery"
              placeholder="Cerca conversazione"
              class="search-input"
            />
          </div>
        </div>
  
        <!-- Lista conversazioni -->
        <div 
          v-if="loading" 
          class="loading-state"
        >
          <div class="loading-spinner"></div>
          <span>Caricamento conversazioni...</span>
        </div>
  
        <div 
          v-else-if="filteredConversations.length === 0" 
          class="no-conversations"
        >
          Nessuna conversazione disponibile per l'inoltro.
        </div>
  
        <div 
          v-else 
          class="conversations-list"
        >
          <div 
            v-for="conversation in filteredConversations" 
            :key="conversation.id"
            class="conversation-item"
            @click="selectConversation(conversation)"
          >
            <img 
              :src="getConversationPhoto(conversation)" 
              :alt="getConversationName(conversation)"
              class="conversation-avatar"
            />
            
            <div class="conversation-info">
              <span class="conversation-name">
                {{ getConversationName(conversation) }}
              </span>
              <span class="conversation-type">
                {{ conversation.type === 'group' ? 'Gruppo' : 'Chat privata' }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    name: 'ForwardModal',
    props: {
      conversations: {
        type: Array,
        required: true
      },
      currentUserId: {
        type: Number,
        required: true
      },
      loading: {
        type: Boolean,
        default: false
      },
      excludeConversationId: {
        type: Number,
        default: null
      }
    },
    data() {
      return {
        searchQuery: ''
      }
    },
    computed: {
      filteredConversations() {
        // Filtra le conversazioni escludendo quella corrente
        const filteredByCurrentConversation = this.conversations.filter(
          conv => conv.id !== this.excludeConversationId
        )
  
        // Se non c'è una query di ricerca, restituisci tutte le conversazioni filtrate
        if (!this.searchQuery.trim()) {
          return filteredByCurrentConversation
        }
  
        // Filtra in base alla query di ricerca
        const query = this.searchQuery.toLowerCase()
        return filteredByCurrentConversation.filter(conv => 
          this.getConversationName(conv).toLowerCase().includes(query)
        )
      }
    },
    methods: {
      getConversationName(conversation) {
        // Per i gruppi, usa il nome del gruppo
        if (conversation.type === 'group') {
          return conversation.name || 'Gruppo senza nome'
        }
        
        // Per le chat private, trova l'altro partecipante
        const otherParticipant = conversation.participants.find(
          p => p.id !== this.currentUserId
        )
        
        return otherParticipant ? otherParticipant.username : 'Chat'
      },
      getConversationPhoto(conversation) {
        // Per i gruppi, usa la foto del gruppo
        if (conversation.type === 'group') {
          return conversation.photoUrl || '/default-avatar.jpeg'
        }
        
        // Per le chat private, usa la foto dell'altro partecipante
        const otherParticipant = conversation.participants.find(
          p => p.id !== this.currentUserId
        )
        
        return otherParticipant?.photoUrl || '/default-avatar.jpeg'
      },
      selectConversation(conversation) {
        this.$emit('select', conversation.id)
      }
    }
  }
  </script>
  
  <style scoped>
  .forward-modal {
    background-color: #1f2c34;
    border-radius: 12px;
    width: 100%;
    max-width: 450px;
    display: flex;
    flex-direction: column;
    max-height: 80vh;
  }
  
  .modal-header {
    padding: 16px 20px;
    border-bottom: 1px solid #2a3942;
    display: flex;
    justify-content: space-between;
    align-items: center;
    background-color: #202c33;
  }
  
  .modal-header h3 {
    color: #e9edef;
    margin: 0;
    font-size: 16px;
  }
  
  .close-button {
    background: none;
    border: none;
    color: #aebac1;
    font-size: 24px;
    cursor: pointer;
    padding: 4px;
    line-height: 1;
    border-radius: 50%;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .close-button:hover {
    background-color: #384147;
  }
  
  .modal-body {
    padding: 20px;
    display: flex;
    flex-direction: column;
    flex: 1;
    overflow-y: auto;
  }
  
  .search-container {
    margin-bottom: 16px;
  }
  
  .search-wrapper {
    position: relative;
    background-color: #2a3942;
    border-radius: 8px;
  }
  
  .search-icon {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    color: #8696a0;
  }
  
  .search-input {
    width: 100%;
    background: none;
    border: none;
    color: #e9edef;
    padding: 12px 12px 12px 40px;
    outline: none;
    font-size: 15px;
  }
  
  .search-input::placeholder {
    color: #8696a0;
  }
  
  .loading-state,
  .no-conversations {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    flex: 1;
    color: #8696a0;
    text-align: center;
    padding: 20px;
  }
  
  .loading-spinner {
    width: 30px;
    height: 30px;
    border: 3px solid #2a3942;
    border-top-color: #00a884;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 10px;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  
  .conversations-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  
  .conversation-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    cursor: pointer;
    transition: background-color 0.2s;
    border-radius: 8px;
  }
  
  .conversation-item:hover {
    background-color: #2a3942;
  }
  
  .conversation-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    object-fit: cover;
  }
  
  .conversation-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  
  .conversation-name {
    color: #e9edef;
    font-size: 16px;
  }
  
  .conversation-type {
    color: #8696a0;
    font-size: 13px;
  }
  </style>