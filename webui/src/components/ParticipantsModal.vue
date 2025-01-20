<!-- ParticipantsModal.vue -->
<template>
  <BaseModal title="Partecipanti del gruppo" @close="$emit('close')">
    <div class="modal-content">
      <!-- Barra di ricerca (solo per gruppi) -->
      <div v-if="isGroup" class="search-container">
        <div class="search-wrapper">
          <i class="fas fa-search search-icon"></i>
          <input 
            type="text" 
            v-model="searchQuery"
            placeholder="Cerca partecipanti"
            class="search-input"
          />
        </div>
      </div>

      <!-- Lista partecipanti -->
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <span>Caricamento partecipanti...</span>
      </div>

      <div v-else-if="filteredParticipants.length === 0" class="no-participants">
        Nessun partecipante trovato.
      </div>

      <div v-else class="participants-list">
        <div 
          v-for="participant in filteredParticipants" 
          :key="participant.id" 
          class="participant-item"
        >
          <img 
            :src="participant.photoUrl || '/api/placeholder/40/40'" 
            :alt="participant.username"
            class="participant-avatar"
          />
    
          <div class="participant-info">
            <span class="participant-name">
              {{ participant.username }}
              <span v-if="participant.id === currentUserId" class="participant-you">
                (Tu)
              </span>
            </span>
          </div>
    
          <button 
            v-if="canRemoveParticipant(participant)"
            class="remove-button"
            @click="$emit('remove-participant', participant)"
          >
            <i class="fas fa-times"></i>
          </button>
        </div>
      </div>

      <!-- Pulsante Aggiungi partecipanti (solo per gruppi) -->
      <button 
        v-if="isGroup" 
        class="add-participants-button"
        @click="$emit('add-participants')"
      >
        <i class="fas fa-user-plus"></i>
        Aggiungi partecipanti
      </button>
    </div>
  </BaseModal>
</template>
  
  <script>
  import BaseModal from './BaseModal.vue'
  export default {
    components: {
      BaseModal
    },
    name: 'ParticipantsModal',
    props: {
      participants: {
        type: Array,
        required: true
      },
      currentUserId: {
        type: Number,
        required: true
      },
      isGroup: {
        type: Boolean,
        default: false
      },
      loading: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        searchQuery: ''
      }
    },
    computed: {
      filteredParticipants() {
        // Se non c'è una query di ricerca, restituisci tutti i partecipanti
        if (!this.searchQuery.trim()) {
          return this.participants
        }
  
        // Filtra in base alla query di ricerca
        const query = this.searchQuery.toLowerCase()
        return this.participants.filter(participant => 
          participant.username.toLowerCase().includes(query)
        )
      }
    },
    methods: {
      canRemoveParticipant(participant) {
        // Può rimuovere partecipanti solo se:
        // 1. È un gruppo
        // 2. Non è l'utente corrente
        return this.isGroup && participant.id !== this.currentUserId
      }
    }
  }
  </script>
  
  <style scoped>
  .participants-modal {
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
  .no-participants {
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
  
  .participants-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 16px;
  }
  
  .participant-item {
    display: flex;
    align-items: center;
    padding: 12px;
    background-color: #202c33;
    border-radius: 8px;
  }
  
  .participant-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    margin-right: 12px;
    object-fit: cover;
  }
  
  .participant-info {
    flex: 1;
  }
  
  .participant-name {
    color: #e9edef;
    font-size: 16px;
    display: flex;
    align-items: center;
    gap: 8px;
  }
  
  .participant-you {
    color: #8696a0;
    font-size: 14px;
  }
  
  .remove-button {
    background: none;
    border: none;
    color: #f15c6d;
    padding: 8px;
    cursor: pointer;
    border-radius: 50%;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .remove-button:hover {
    background-color: rgba(241, 92, 109, 0.1);
  }
  
  .add-participants-button {
    background-color: #00a884;
    color: #fff;
    border: none;
    border-radius: 8px;
    padding: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    cursor: pointer;
    transition: background-color 0.2s;
  }
  
  .add-participants-button:hover {
    background-color: #008069;
  }
  </style>