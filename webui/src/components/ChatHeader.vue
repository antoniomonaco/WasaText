<!-- ChatHeader.vue -->
<template>
  <header class="chat-header">
    <div class="chat-header-info" @click="$emit('show-participants')">
      <img 
        :src="conversationPhoto" 
        :alt="conversationName"
        class="chat-avatar"
      />
      <div class="chat-details">
        <div class="chat-name">{{ conversationName }}</div>
        <div class="chat-status">
          {{ statusText }}
        </div>
      </div>
    </div>
    
    <!-- Azioni per i gruppi -->
    <div 
      v-if="isGroup" 
      class="chat-header-actions"
    >
      <button 
        class="icon-button" 
        @click="toggleGroupMenu"
      >
        <i class="fas fa-ellipsis-v"></i>
      </button>
      
      <!-- Menu contestuale del gruppo -->
      <div 
        v-if="showGroupMenu" 
        class="group-menu"
      >
        <div 
          class="menu-item" 
          @click="$emit('add-participants')"
        >
          <i class="fas fa-user-plus"></i>
          Aggiungi partecipanti
        </div>
        <div 
          class="menu-item" 
          @click="$emit('show-participants')"
        >
          <i class="fas fa-users"></i>
          Visualizza partecipanti
        </div>
        <div 
          class="menu-item" 
          @click="$emit('edit-group')"
        >
          <i class="fas fa-edit"></i>
          Modifica gruppo
        </div>
        <div 
          class="menu-item delete" 
          @click="$emit('leave-group')"
        >
          <i class="fas fa-sign-out-alt"></i>
          Esci dal gruppo
        </div>
      </div>
    </div>
  </header>
</template>

<script>
import { onMounted, onUnmounted } from 'vue'

export default {
  name: 'ChatHeader',
  props: {
    conversationName: {
      type: String,
      required: true
    },
    conversationPhoto: {
      type: String,
      default: '/default-avatar.jpeg'
    },
    isGroup: {
      type: Boolean,
      default: false
    },
    participantsCount: {
      type: Number,
      default: 0
    }
  },
  data() {
    return {
      showGroupMenu: false
    }
  },
  computed: {
    statusText() {
      if (this.isGroup) {
        return `${this.participantsCount} partecipanti`
      }
      return 'online'
    }
  },
  methods: {
    toggleGroupMenu() {
      this.showGroupMenu = !this.showGroupMenu
    },
    handleOutsideClick(event) {
      if (this.showGroupMenu && !this.$el.contains(event.target)) {
        this.showGroupMenu = false
      }
    }
  },
  mounted() {
    // Usa addEventListener invece di this.$once
    document.addEventListener('click', this.handleOutsideClick)
  },
  beforeUnmount() {
    // Rimuovi l'event listener quando il componente viene distrutto
    document.removeEventListener('click', this.handleOutsideClick)
  }
}
</script>

<style scoped>
/* Stili precedenti rimangono invariati */
.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  background-color: #202c33;
  height: 60px;
  border-left: 1px solid #2a3942;
  flex-shrink: 0;
  position: relative;
  z-index: 10;
}

.chat-header-info {
  display: flex;
  align-items: center;
  gap: 15px;
  cursor: pointer;
}

.chat-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.chat-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.chat-name {
  color: #e9edef;
  font-size: 16px;
  font-weight: 500;
}

.chat-status {
  color: #8696a0;
  font-size: 13px;
}

.chat-header-actions {
  position: relative;
}

.icon-button {
  background: none;
  border: none;
  color: #aebac1;
  padding: 8px;
  cursor: pointer;
  border-radius: 50%;
  transition: background-color 0.2s;
}

.icon-button:hover {
  background-color: #384147;
}

.group-menu {
  position: absolute;
  top: 100%;
  right: 0;
  background-color: #233138;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.2);
  z-index: 1000;
  min-width: 220px;
  padding: 8px 0;
  margin-top: 8px;
}

.group-menu .menu-item {
  padding: 10px 16px;
  color: #e9edef;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
  transition: background-color 0.2s;
}

.group-menu .menu-item:hover {
  background-color: #182229;
}

.group-menu .menu-item i {
  width: 20px;
  text-align: center;
  font-size: 16px;
}

.group-menu .menu-item.delete {
  color: #f15c6d;
}

.group-menu .menu-item.delete i {
  color: #f15c6d;
}
</style>