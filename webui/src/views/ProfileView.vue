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
                <button class="change-photo-button" @click="$refs.photoInput.click()">
                  <i class="fas fa-camera"></i>
                </button>
              </div>
            </div>
            <input 
              ref="photoInput"
              type="url"
              v-model="photoUrl"
              placeholder="Inserisci l'URL della foto"
              class="photo-input"
            />
            <div class="photo-actions" v-if="photoUrl !== currentUser.photoUrl">
              <button 
                class="action-button save"
                @click="updatePhoto"
                :disabled="!isValidUrl(photoUrl)"
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
  
        <!-- Notifications for actions -->
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
        usernameError: '',
        notification: null,
        currentUserId: Number(localStorage.getItem('authToken')),
        accountCreationDate: new Date()
    }
  },
  methods: {
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
      } catch (error) {
        console.error('Errore nel recupero del profilo:', error);
      }
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
                this.username,
                {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem('authToken')}`
                    }
                }
            );
            
            // Aggiorno il nome utente nel localStorage
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

    isValidUrl(url) {
      if (!url) return false;
      try {
        new URL(url);
        return true;
      } catch {
        return false;
      }
    },

    async updatePhoto() {
        if (!this.isValidUrl(this.photoUrl)) {
            this.showNotification('URL della foto non valido', 'error');
            return;
        }

        try {
            await this.$axios.put(
                '/users/me/photo',
                { photoUrl: this.photoUrl },
                {
                    headers: {
                        Authorization: `Bearer ${localStorage.getItem('authToken')}`
                    }
                }
            );
            
            localStorage.setItem('userPhoto', this.photoUrl);
            this.currentUser.photoUrl = this.photoUrl;

            // Notifico gli altri componenti del cambiamento
            window.dispatchEvent(new Event('userProfileUpdated'));

            this.showNotification('Foto profilo aggiornata con successo', 'success');
        } catch (error) {
            this.showNotification('Errore nell\'aggiornamento della foto profilo', 'error');
        }
    },

    cancelPhotoUpdate() {
      this.photoUrl = this.currentUser.photoUrl;
      this.previewPhotoUrl = '';
    },

    showNotification(message, type = 'info') {
      this.notification = { message, type };
      setTimeout(() => {
        this.notification = null;
      }, 3000);
    },

    formatDate(date) {
      return new Date(date).toLocaleDateString('it-IT', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      });
    }
  },
  created() {
    this.fetchUserProfile();
  },
  watch: {
    photoUrl(newUrl) {
      if (this.isValidUrl(newUrl)) {
        this.previewPhotoUrl = newUrl;
      } else {
        this.previewPhotoUrl = '';
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