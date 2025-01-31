<!-- ForwardModal.vue -->
<template>
  <div class="forward-modal">
    <div class="search-container">
      <div class="search-wrapper">
        <i class="fas fa-search search-icon"></i>
        <input 
          type="text" 
          v-model="searchQuery"
          placeholder="Cerca una chat o un utente"
          class="search-input"
        />
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <span>Caricamento...</span>
    </div>

    <div v-else-if="filteredTargets.length === 0" class="no-targets">
      Nessun risultato trovato
    </div>

    <div v-else class="forward-targets-list">
      <div 
        v-for="target in filteredTargets" 
        :key="target.isNewChat ? 'new-' + target.participants[0].id : target.id"
        class="target-item"
        @click="$emit('select', target.isNewChat ? target : target.id)"
      >
        <img 
          :src="getTargetPhoto(target)" 
          :alt="getTargetName(target)"
          class="target-avatar"
        />
        <div class="target-info">
          <span class="target-name">{{ getTargetName(target) }}</span>
          <span class="target-type">
            {{ target.isNewChat ? 'Nuova chat' : (target.type === 'group' ? 'Gruppo' : 'Chat') }}
          </span>
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
    filteredTargets() {
      // Prima filtra le conversazioni escludendo quella corrente
      let targets = this.conversations.filter(target => 
        target.id !== this.excludeConversationId
      );

      // Se c'è una query di ricerca, filtra ulteriormente
      if (this.searchQuery.trim()) {
        const query = this.searchQuery.toLowerCase();
        targets = targets.filter(target => {
          const name = this.getTargetName(target).toLowerCase();
          return name.includes(query);
        });
      }

      return targets;
    }
  },
  methods: {
    getTargetName(target) {
      if (target.isNewChat) {
        return target.participants[0].username;
      }
      
      if (target.type === 'group') {
        return target.name || 'Gruppo senza nome';
      }
      
      return this.getOtherParticipantName(target);
    },
    
    getTargetPhoto(target) {
      if (target.isNewChat) {
        return target.participants[0].photoUrl || '/default-avatar.jpeg';
      }
      
      if (target.type === 'group') {
        return target.photoUrl || '/default-avatar.jpeg';
      }
      
      const otherParticipant = target.participants?.find(
        p => p.id !== this.currentUserId
      );
      return otherParticipant?.photoUrl || '/default-avatar.jpeg';
    },
    
    getOtherParticipantName(conversation) {
      const otherParticipant = conversation.participants?.find(
        p => p.id !== this.currentUserId
      );
      return otherParticipant?.username || 'Chat';
    }
  }
}
</script>

<style scoped>
.forward-modal {
  display: flex;
  flex-direction: column;
  gap: 16px;
  height: 100%;
  max-height: 80vh;
}

.search-container {
  padding: 8px 16px;
}

.search-wrapper {
  position: relative;
  background-color: #2a3942;
  border-radius: 8px;
  padding: 8px 12px;
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
  padding: 6px 35px;
  outline: none;
  font-size: 15px;
}

.search-input::placeholder {
  color: #8696a0;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px;
  color: #8696a0;
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #2a3942;
  border-top-color: #00a884;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 12px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.no-targets {
  text-align: center;
  color: #8696a0;
  padding: 32px;
}

.forward-targets-list {
  overflow-y: auto;
  padding: 0 16px;
}

.target-item {
  display: flex;
  align-items: center;
  padding: 12px;
  cursor: pointer;
  border-radius: 8px;
  transition: background-color 0.2s;
}

.target-item:hover {
  background-color: #2a3942;
}

.target-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  margin-right: 12px;
  object-fit: cover;
}

.target-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.target-name {
  color: #e9edef;
  font-size: 16px;
}

.target-type {
  color: #8696a0;
  font-size: 13px;
}

/* Scrollbar styling */
.forward-targets-list::-webkit-scrollbar {
  width: 6px;
}

.forward-targets-list::-webkit-scrollbar-track {
  background: transparent;
}

.forward-targets-list::-webkit-scrollbar-thumb {
  background: rgba(134, 150, 160, 0.2);
  border-radius: 3px;
}

.forward-targets-list::-webkit-scrollbar-thumb:hover {
  background: rgba(134, 150, 160, 0.3);
}
</style>