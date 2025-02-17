<!-- NewChatModal.vue -->
<template>
  <div class="modal">
    <div class="modal-content">
      <div class="modal-header">
        <h3>Nuova conversazione</h3>
        <button class="close-button" @click="$emit('close')">×</button>
      </div>

      <div class="modal-body">
        <!-- Tabs per scegliere tra chat e gruppo -->
        <div class="chat-type-tabs">
          <button 
            class="tab-button"
            :class="{ active: activeTab === 'chat' }"
            @click="activeTab = 'chat'"
          >
            <i class="fas fa-comments"></i>
            Nuova chat
          </button>
          <button 
            class="tab-button"
            :class="{ active: activeTab === 'group' }"
            @click="activeTab = 'group'"
          >
            <i class="fas fa-users"></i>
            Nuovo gruppo
          </button>
        </div>

        <!-- Form creazione gruppo -->
        <div v-if="activeTab === 'group'" class="group-form">
          <input 
            type="text" 
            v-model="groupName"
            placeholder="Nome del gruppo"
            class="group-input"
          />

          <!-- Sezione foto gruppo -->
          <div class="photo-section">
            <div class="photo-preview-container" v-if="previewPhotoUrl || groupPhotoUrl">
              <img 
                :src="previewPhotoUrl || groupPhotoUrl" 
                alt="Preview" 
                class="photo-preview"
                @error="handleImageError"
              />
              <button class="remove-image-button" @click="removeImage">
                <i class="fas fa-times"></i>
              </button>
            </div>

            <div class="upload-options">
              <button 
                class="option-button"
                :class="{ 'active': uploadMode === 'url' }"
                @click="uploadMode = 'url'"
              >
                URL
              </button>
              <button 
                class="option-button"
                :class="{ 'active': uploadMode === 'file' }"
                @click="uploadMode = 'file'"
              >
                Carica file
              </button>
            </div>

            <input 
              v-if="uploadMode === 'url'"
              type="url" 
              v-model="groupPhotoUrl"
              placeholder="URL foto del gruppo"
              class="group-input"
              @input="validatePhotoUrl"
            />

            <div class="file-upload-container" v-else>
              <button class="upload-button" @click="$refs.fileInput.click()">
                <i class="fas fa-camera"></i>
                Scegli una foto
              </button>
              <input 
                ref="fileInput"
                type="file"
                @change="handleFileUpload"
                accept="image/*"
                class="file-input"
              />
            </div>

            <div v-if="photoError" class="error-message">
              {{ photoError }}
            </div>
            <div v-if="showSizeWarning" class="error-message">
              Immagine troppo grande! Dimensione massima: 5MB
            </div>
          </div>
        </div>

        <!-- Ricerca utenti -->
        <input 
          type="text" 
          placeholder="Cerca un utente"
          v-model="searchQuery"
          class="search-input"
        />

        <!-- Lista utenti -->
        <div class="users-list">
          <div 
            v-for="user in filteredUsers" 
            :key="user.id"
            class="user-item"
            :class="{ selected: selectedUsers.includes(user.id) }"
            @click="toggleUserSelection(user)"
          >
            <img 
              :src="user.photoUrl || '/default-avatar.jpeg'" 
              alt="User avatar" 
              class="user-avatar"
            />
            <span class="user-name">{{ user.username }}</span>
            <span v-if="selectedUsers.includes(user.id)" class="check-icon">
              <i class="fas fa-check"></i>
            </span>
          </div>
        </div>

        <!-- Pulsante crea -->
        <button 
          class="create-button"
          :disabled="!canCreate"
          @click="createConversation"
        >
          {{ createButtonText }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'NewChatModal',
  props: {
    users: {
      type: Array,
      required: true
    }
  },
  data() {
    return {
      activeTab: 'chat',
      searchQuery: '',
      selectedUsers: [],
      groupName: '',
      groupPhotoUrl: '',
      uploadMode: 'file',
      photoError: '',
      showSizeWarning: false,
      previewPhotoUrl: null,
      photoBase64: null,
      isProcessingImage: false
    }
  },
  computed: {
    filteredUsers() {
      if (!this.searchQuery) return this.users
      const query = this.searchQuery.toLowerCase()
      return this.users.filter(user => 
        user.username.toLowerCase().includes(query)
      )
    },
    canCreate() {
      if (this.activeTab === 'chat') {
        return this.selectedUsers.length === 1
      }
      return this.selectedUsers.length >= 1 && this.groupName.trim().length >= 3
    },
    createButtonText() {
      if (this.activeTab === 'chat') {
        return 'Inizia chat'
      }
      return 'Crea gruppo'
    }
  },
  methods: {
    toggleUserSelection(user) {
      const index = this.selectedUsers.indexOf(user.id)
      if (this.activeTab === 'chat') {
        // Per le chat private, permette di selezionare solo un utente
        this.selectedUsers = [user.id]
      } else {
        // Per i gruppi, permette selezione multipla
        if (index === -1) { // Se l'utente NON è già selezionato
          this.selectedUsers.push(user.id)
        } else {
          this.selectedUsers.splice(index, 1)
        }
      }
    },
    async handleFileUpload(event) {
      const file = event.target.files[0];
      if (!file) return;

      const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB in bytes

      if (file.size > MAX_FILE_SIZE) {
        this.showSizeWarning = true;
        setTimeout(() => {
          this.showSizeWarning = false;
        }, 3000);
        return;
      }

      this.isProcessingImage = true;
      this.showSizeWarning = false;
      this.photoError = '';

      try {
        const base64 = await this.convertToBase64(file);
        this.photoBase64 = base64;
        this.previewPhotoUrl = base64;
        this.groupPhotoUrl = '';
      } catch (error) {
        console.error('Errore nella conversione dell\'immagine:', error);
        this.photoError = 'Errore nel caricamento dell\'immagine';
      } finally {
        this.isProcessingImage = false;
      }

      event.target.value = '';
    },

    convertToBase64(file) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = () => resolve(reader.result);
        reader.onerror = error => reject(error);
        reader.readAsDataURL(file);
      });
    },

    validatePhotoUrl() {
      this.photoError = '';
      if (this.groupPhotoUrl && !this.isValidUrl(this.groupPhotoUrl)) {
        this.photoError = 'URL non valido';
      }
    },

    isValidUrl(url) {
      try {
        new URL(url);
        return true;
      } catch {
        return false;
      }
    },

    handleImageError() {
      this.photoError = 'Impossibile caricare l\'immagine';
      this.previewPhotoUrl = null;
    },

    removeImage() {
      this.groupPhotoUrl = '';
      this.previewPhotoUrl = null;
      this.photoBase64 = null;
      this.photoError = '';
      this.showSizeWarning = false;
    },

    createConversation() {
      if (!this.canCreate) return

      // Creazione della struttura dati corretta per la richiesta
      const conversationData = {
        type: this.activeTab === 'group' ? 'group' : 'chat',
        participants: this.selectedUsers,
      }

      if (this.activeTab === 'group') {
        // Aggiungi name e photoUrl solo per i gruppi
        conversationData.name = this.groupName.trim()
        if (this.uploadMode === 'url' && this.groupPhotoUrl) {
          conversationData.photoUrl = this.groupPhotoUrl.trim()
        } else if (this.uploadMode === 'file' && this.photoBase64) {
          conversationData.photoUrl = this.photoBase64
        }
      }

      console.log('Creating conversation with data:', conversationData)
      this.$emit('create', conversationData)
    },
    resetForm() {
      this.activeTab = 'chat'
      this.searchQuery = ''
      this.selectedUsers = []
      this.groupName = ''
      this.groupPhotoUrl = ''
    }
  }
}
</script>

<style scoped>
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: #233138;
  width: 400px;
  border-radius: 12px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 16px;
  border-bottom: 1px solid #2a3942;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  color: #e9edef;
  margin: 0;
  font-size: 18px;
}

.close-button {
  background: none;
  border: none;
  color: #aebac1;
  font-size: 24px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 4px;
}

.close-button:hover {
  background-color: #384147;
}

.modal-body {
  padding: 16px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.chat-type-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 8px;
}

.tab-button {
  flex: 1;
  background-color: #2a3942;
  border: none;
  color: #8696a0;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.2s;
}

.tab-button.active {
  background-color: #00a884;
  color: #fff;
}

.tab-button i {
  font-size: 16px;
}

.group-form {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.group-input,
.search-input {
  width: 100%;
  background-color: #2a3942;
  border: none;
  border-radius: 8px;
  color: #e9edef;
  padding: 12px;
  font-size: 15px;
  outline: none;
}

.group-input::placeholder,
.search-input::placeholder {
  color: #8696a0;
}

.users-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: 300px;
  overflow-y: auto;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
  position: relative;
}

.user-item:hover {
  background-color: #2a3942;
}

.user-item.selected {
  background-color: rgba(0, 168, 132, 0.1);
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 12px;
}

.user-name {
  color: #e9edef;
  flex: 1;
}

.check-icon {
  color: #00a884;
  font-size: 18px;
}

.create-button {
  background-color: #00a884;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 12px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-top: 8px;
}

.create-button:disabled {
  background-color: #2a3942;
  color: #8696a0;
  cursor: not-allowed;
}

.create-button:not(:disabled):hover {
  background-color: #008069;
}

/* Stili per la sezione foto */
.photo-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 8px;
}

.photo-preview-container {
  position: relative;
  width: 100%;
  max-width: 200px;
  margin: 0 auto;
}

.photo-preview {
  width: 100%;
  height: auto;
  border-radius: 8px;
  display: block;
}

.remove-image-button {
  position: absolute;
  top: -8px;
  right: -8px;
  background-color: #202c33;
  border: none;
  color: #e9edef;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 1px 3px rgba(0,0,0,0.3);
}

.upload-options {
  display: flex;
  gap: 8px;
}

.option-button {
  flex: 1;
  background-color: #2a3942;
  border: none;
  color: #8696a0;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.option-button.active {
  background-color: #00a884;
  color: #fff;
}

.file-upload-container {
  width: 100%;
}

.upload-button {
  width: 100%;
  background-color: #2a3942;
  border: none;
  color: #e9edef;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: background-color 0.2s;
}

.upload-button:hover {
  background-color: #374045;
}

.file-input {
  display: none;
}

.error-message {
  color: #f15c6d;
  font-size: 13px;
  text-align: center;
}
</style>