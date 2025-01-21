<!-- GroupEditModal.vue -->
<template>
  <div class="group-edit-modal">
    <div class="modal-header">
      <h3>Modifica gruppo</h3>
      <button 
        class="close-button" 
        @click="$emit('close')"
      >
        <i class="fas fa-times"></i>
      </button>
    </div>
    
    <div class="modal-body">
      <!-- Form per il nome del gruppo -->
      <div class="form-group">
        <label for="groupName">Nome del gruppo</label>
        <input 
          id="groupName"
          v-model="editedName"
          type="text"
          :placeholder="placeholder"
          class="group-input"
          @input="validateName"
        />
        <span 
          v-if="nameError" 
          class="error-message"
        >
          {{ nameError }}
        </span>
      </div>

      <!-- Form per la foto del gruppo -->
      <div class="form-group">
        <label>Foto del gruppo (opzionale)</label>
        
        <div class="photo-section">
          <!-- Anteprima foto -->
          <div v-if="previewPhotoUrl || editedPhotoUrl" class="photo-preview-container">
            <img 
              :src="previewPhotoUrl || editedPhotoUrl" 
              alt="Anteprima foto gruppo"
              class="photo-preview"
              @error="handleImageError"
            />
            <button class="remove-image-button" @click="removeImage">
              <i class="fas fa-times"></i>
            </button>
          </div>

          <!-- Opzioni di upload -->
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

          <!-- Input URL -->
          <input 
            v-if="uploadMode === 'url'"
            v-model="editedPhotoUrl"
            type="url"
            placeholder="Inserisci l'URL dell'immagine"
            class="group-input"
            @input="validatePhotoUrl"
          />

          <!-- Upload file -->
          <div v-else class="file-upload-container">
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

          <!-- Messaggi di errore -->
          <span 
            v-if="photoError" 
            class="error-message"
          >
            {{ photoError }}
          </span>
          <div v-if="showSizeWarning" class="error-message">
            Immagine troppo grande! Dimensione massima: 5MB
          </div>
        </div>
      </div>

      <!-- Pulsante Salva modifiche -->
      <button 
        class="save-button"
        :disabled="!isFormValid"
        @click="saveChanges"
      >
        Salva modifiche
      </button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'GroupEditModal',
  props: {
    initialName: {
      type: String,
      default: ''
    },
    initialPhotoUrl: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      editedName: this.initialName,
      editedPhotoUrl: this.initialPhotoUrl,
      nameError: '',
      photoError: '',
      placeholder: 'Es. Famiglia, Amici, Lavoro...',
      uploadMode: 'file',
      previewPhotoUrl: null,
      photoBase64: null,
      showSizeWarning: false,
      isProcessingImage: false
    }
  },
  computed: {
    isFormValid() {
      // Il nome deve essere valido e non deve esserci un errore
      const hasValidChanges = this.editedName !== this.initialName || 
                            this.hasPhotoChanges;
      
      return this.editedName.trim().length >= 3 && 
             !this.nameError && 
             !this.photoError &&
             hasValidChanges;
    },
    hasPhotoChanges() {
      return (this.uploadMode === 'url' && this.editedPhotoUrl !== this.initialPhotoUrl) ||
             (this.uploadMode === 'file' && this.photoBase64 !== null);
    }
  },
  methods: {
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
        this.editedPhotoUrl = '';
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

    validateName() {
      this.nameError = '';
      
      const trimmedName = this.editedName.trim();
      
      if (!trimmedName) {
        this.nameError = 'Il nome del gruppo è obbligatorio';
        return;
      }

      if (trimmedName.length < 3) {
        this.nameError = 'Il nome deve contenere almeno 3 caratteri';
        return;
      }

      if (trimmedName.length > 30) {
        this.nameError = 'Il nome non può superare i 30 caratteri';
      }
    },

    validatePhotoUrl() {
      this.photoError = '';
      this.photoBase64 = null;
      this.previewPhotoUrl = null;

      if (!this.editedPhotoUrl) return;

      if (!this.isValidUrl(this.editedPhotoUrl)) {
        this.photoError = 'Inserisci un URL valido';
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
      this.editedPhotoUrl = '';
      this.previewPhotoUrl = null;
      this.photoBase64 = null;
      this.photoError = '';
      this.showSizeWarning = false;
    },

    saveChanges() {
      if (!this.isFormValid) return;

      const groupData = {
        name: this.editedName.trim()
      };

      // Aggiungi la foto solo se è stata modificata
      if (this.uploadMode === 'url' && this.editedPhotoUrl) {
        groupData.photoUrl = this.editedPhotoUrl.trim();
      } else if (this.uploadMode === 'file' && this.photoBase64) {
        groupData.photoUrl = this.photoBase64;
      }

      this.$emit('save', groupData);
    }
  }
}
</script>

<style scoped>
.group-edit-modal {
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
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

label {
  color: #e9edef;
  font-size: 14px;
}

.group-input {
  background-color: #2a3942;
  border: none;
  border-radius: 8px;
  color: #e9edef;
  padding: 12px;
  font-size: 15px;
  outline: none;
  transition: box-shadow 0.2s;
}

.group-input:focus {
  box-shadow: 0 0 0 2px #00a884;
}

.group-input::placeholder {
  color: #8696a0;
}

.error-message {
  color: #f15c6d;
  font-size: 13px;
}

.photo-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
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

.file-input {
  display: none;
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

.save-button {
  background-color: #00a884;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 12px;
  font-size: 15px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.save-button:hover:not(:disabled) {
  background-color: #008069;
}

.save-button:disabled {
  background-color: #2a3942;
  color: #8696a0;
  cursor: not-allowed;
}
</style>