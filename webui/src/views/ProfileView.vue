<!-- ProfileView.vue -->
<template>
  <header>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css">
  </header>
  <div class="profile-container">
    <div class="profile-card">
      <div class="profile-header">
        <button class="back-button" @click="$router.push('/')">
          <i class="fas fa-arrow-left"></i>
        </button>
        <h2>Il tuo profilo</h2>
      </div>

      <div class="profile-content">
        <!-- Photo Section -->
        <div class="photo-section">
          <div class="photo-container">
            <img 
              :src="previewPhotoUrl || currentUser.photoUrl || '/default-avatar.jpeg'" 
              :alt="currentUser.username"
              class="profile-photo"
            />
            <div class="photo-overlay">
              <button class="change-photo-button" @click="$refs.fileInput.click()">
                <i class="fas fa-camera"></i>
              </button>
            </div>
          </div>

          <!-- Input nascosto per il file -->
          <input 
            ref="fileInput"
            type="file"
            @change="handleFileUpload"
            accept="image/*"
            class="file-input"
          />

          <!-- Input URL opzionale -->
          <div class="upload-options">
            <button 
              class="option-button"
              :class="{ 'active': uploadMode === 'url' }"
              @click="uploadMode = 'url'"
            >
              Usa URL
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
            v-model="photoUrl"
            placeholder="Inserisci l'URL della foto"
            class="photo-input"
          />

          <div v-if="showSizeWarning" class="error-message">
            Immagine troppo grande! Dimensione massima: 5MB
          </div>

          <div class="photo-actions" v-if="hasChanges">
            <button 
              class="action-button save"
              @click="updatePhoto"
              :disabled="!canSave"
            >
              Salva foto
            </button>
            <button 
              class="action-button cancel"
              @click="cancelPhotoUpdate"
            >
              Annulla
            </button>
          </div>
        </div>

        <!-- Username Section -->
        <div class="username-section">
          <label for="username">Nome utente</label>
          <div class="username-input-container">
            <input 
              id="username"
              type="text"
              v-model="username"
              :placeholder="currentUser.username"
              class="username-input"
              @input="validateUsername"
            />
            <button 
              v-if="username !== currentUser.username && !usernameError"
              class="save-username-button"
              @click="updateUsername"
            >
              <i class="fas fa-check"></i>
            </button>
          </div>
          <span class="error-message" v-if="usernameError">
            {{ usernameError }}
          </span>
        </div>

        <!-- Account Info -->
        <div class="account-info">
          <h3>Informazioni account</h3>
          <p class="info-item">
            <span class="info-label">ID Account:</span>
            <span class="info-value">{{ currentUserId }}</span>
          </p>
        </div>
      </div>

      <!-- Notifications -->
      <div 
        v-if="notification"
        class="notification"
        :class="notification.type"
      >
        {{ notification.message }}
      </div>
    </div>
  </div>
</template>

<script>
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB in bytes

export default {
  name: 'ProfileView',
  data() {
    return {
      currentUser: {
        id: Number(localStorage.getItem('authToken')),
        username: localStorage.getItem('username') || '',
        photoUrl: ''
      },
      username: localStorage.getItem('username') || '',
      photoUrl: '',
      previewPhotoUrl: '',
      photoBase64: null,
      usernameError: '',
      notification: null,
      currentUserId: Number(localStorage.getItem('authToken')),
      uploadMode: 'file',
      showSizeWarning: false,
      isProcessingImage: false
    }
  },
  computed: {
    hasChanges() {
      return (this.photoUrl && this.photoUrl !== this.currentUser.photoUrl) || 
             (this.photoBase64 !== null);
    },
    canSave() {
      return (this.uploadMode === 'url' && this.isValidUrl(this.photoUrl)) ||
             (this.uploadMode === 'file' && this.photoBase64 !== null);
    }
  },
  methods: {
    async handleFileUpload(event) {
      const file = event.target.files[0];
      if (!file) return;

      // Verifica dimensione
      if (file.size > MAX_FILE_SIZE) {
        this.showSizeWarning = true;
        setTimeout(() => {
          this.showSizeWarning = false;
        }, 3000);
        return;
      }

      this.isProcessingImage = true;
      this.showSizeWarning = false;

      try {
        // Converti in base64
        this.photoBase64 = await this.convertToBase64(file);
        this.previewPhotoUrl = this.photoBase64;
        this.uploadMode = 'file';
      } catch (error) {
        console.error('Errore nella conversione dell\'immagine:', error);
      } finally {
        this.isProcessingImage = false;
      }

      // Reset input file
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

    async updatePhoto() {
      if (!this.canSave) return;

      try {
        const photoUrlToSave = this.uploadMode === 'file' ? this.photoBase64 : this.photoUrl;
        
        await this.$axios.put(
          '/users/me/photo',
          { photoUrl: photoUrlToSave },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        );
        
        localStorage.setItem('userPhoto', photoUrlToSave);
        this.currentUser.photoUrl = photoUrlToSave;
        
        window.dispatchEvent(new Event('userProfileUpdated'));
        this.showNotification('Foto profilo aggiornata con successo', 'success');
        
        // Reset dello stato
        this.photoUrl = '';
        this.photoBase64 = null;
        this.uploadMode = 'file';
      } catch (error) {
        this.showNotification('Errore nell\'aggiornamento della foto profilo', 'error');
      }
    },

    cancelPhotoUpdate() {
      this.photoUrl = '';
      this.photoBase64 = null;
      this.previewPhotoUrl = this.currentUser.photoUrl;
      this.uploadMode = 'file';
      this.showSizeWarning = false;
    },

    validateUsername() {
      this.usernameError = '';
      
      if (this.username.length < 3) {
        this.usernameError = 'Il nome utente deve contenere almeno 3 caratteri';
        return false;
      }
      
      if (this.username.length > 16) {
        this.usernameError = 'Il nome utente non può superare i 16 caratteri';
        return false;
      }
      
      const validUsernameRegex = /^[a-zA-Z0-9]+$/;
      if (!validUsernameRegex.test(this.username)) {
        this.usernameError = 'Il nome utente può contenere solo lettere e numeri';
        return false;
      }
      
      return true;
    },

    async updateUsername() {
      if (!this.validateUsername()) return;

      try {
        await this.$axios.put(
          '/users/me/name',
          { name: this.username },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        );
        
        localStorage.setItem('username', this.username);
        this.currentUser.username = this.username;
        this.showNotification('Nome utente aggiornato con successo', 'success');
      } catch (error) {
        if (error.response?.status === 409) {
          this.usernameError = 'Nome utente già in uso';
        } else {
          this.showNotification('Errore nell\'aggiornamento del nome utente', 'error');
        }
      }
    },

    showNotification(message, type = 'info') {
      this.notification = { message, type };
      setTimeout(() => {
        this.notification = null;
      }, 3000);
    },

    isValidUrl(url) {
      if (!url) return false;
      try {
        new URL(url);
        return true;
      } catch {
        return false;
      }
    },

    async fetchUserProfile() {
      try {
        const response = await this.$axios.get('/users/me', {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          }
        });
        
        this.currentUser = response.data;
        this.username = this.currentUser.username;
        this.photoUrl = this.currentUser.photoUrl;
        this.previewPhotoUrl = this.currentUser.photoUrl;
      } catch (error) {
        console.error('Errore nel recupero del profilo:', error);
      }
    }
  },
  created() {
    this.fetchUserProfile();
  },
  watch: {
    photoUrl(newUrl) {
      if (this.uploadMode === 'url' && this.isValidUrl(newUrl)) {
        this.previewPhotoUrl = newUrl;
      }
    }
  }
}
</script>
  
<style scoped>
.profile-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #111b21;
  padding: 20px;
}

.profile-card {
  background-color: #202c33;
  border-radius: 12px;
  width: 100%;
  max-width: 500px;
  padding: 24px;
  position: relative;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.back-button {
  background: none;
  border: none;
  color: #8696a0;
  font-size: 20px;
  cursor: pointer;
  padding: 8px;
  border-radius: 50%;
  transition: all 0.2s;
}

.back-button:hover {
  background-color: #2a3942;
  color: #e9edef;
}

.profile-header h2 {
  color: #e9edef;
  margin: 0;
  font-size: 20px;
}

.profile-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.photo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.photo-container {
  position: relative;
  width: 150px;
  height: 150px;
  border-radius: 50%;
  overflow: hidden;
}

.profile-photo {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.photo-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}

.photo-container:hover .photo-overlay {
  opacity: 1;
}

.change-photo-button {
  background-color: #00a884;
  border: none;
  color: white;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.photo-input {
  width: 100%;
  background-color: #2a3942;
  border: none;
  border-radius: 8px;
  padding: 12px;
  color: #e9edef;
  outline: none;
}

.photo-input::placeholder {
  color: #8696a0;
}

.photo-actions {
  display: flex;
  gap: 8px;
}

.action-button {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.action-button.save {
  background-color: #00a884;
  color: white;
}

.action-button.save:disabled {
  background-color: #2a3942;
  color: #8696a0;
  cursor: not-allowed;
}

.action-button.cancel {
  background-color: #2a3942;
  color: #e9edef;
}

.username-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.username-section label {
  color: #8696a0;
  font-size: 14px;
}

.username-input-container {
  display: flex;
  gap: 8px;
}

.username-input {
  flex: 1;
  background-color: #2a3942;
  border: none;
  border-radius: 8px;
  padding: 12px;
  color: #e9edef;
  outline: none;
}

.save-username-button {
  background-color: #00a884;
  border: none;
  border-radius: 8px;
  width: 40px;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
}

.save-username-button:hover {
  background-color: #008069;
}

.error-message {
  color: #f15c6d;
  font-size: 13px;
}

.account-info {
  background-color: #1f2c34;
  border-radius: 8px;
  padding: 16px;
}

.account-info h3 {
  color: #e9edef;
  margin: 0 0 12px 0;
  font-size: 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin: 8px 0;
  color: #8696a0;
  font-size: 14px;
}

.info-value {
  color: #e9edef;
}

.notification {
  position: absolute;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 14px;
  animation: slideUp 0.3s ease-out;
}

.notification.success {
  background-color: #00a884;
  color: white;
}

.notification.error {
  background-color: #f15c6d;
  color: white;
}
.file-input {
  display: none;
}

.upload-options {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
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
@keyframes slideUp {
  from {
    transform: translate(-50%, 20px);
    opacity: 0;
  }
  to {
    transform: translate(-50%, 0);
    opacity: 1;
  }
}
</style>