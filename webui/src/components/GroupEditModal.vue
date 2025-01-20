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
          <label for="groupPhoto">Foto del gruppo (opzionale)</label>
          <div class="photo-input-container">
            <input 
              id="groupPhoto"
              v-model="editedPhotoUrl"
              type="url"
              placeholder="Inserisci l'URL dell'immagine"
              class="group-input"
              @input="validatePhotoUrl"
            />
            <div 
              v-if="editedPhotoUrl" 
              class="photo-preview"
            >
              <img 
                :src="editedPhotoUrl" 
                :alt="'Anteprima foto gruppo'"
                @error="handleImageError"
              />
            </div>
          </div>
          <span 
            v-if="photoUrlError" 
            class="error-message"
          >
            {{ photoUrlError }}
          </span>
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
        photoUrlError: '',
        placeholder: 'Es. Famiglia, Amici, Lavoro...',
        isPhotoValid: true
      }
    },
    computed: {
      isFormValid() {
        // Il nome deve essere valido e non deve esserci un errore
        return this.editedName.trim().length >= 3 && 
               !this.nameError && 
               (!this.editedPhotoUrl || 
                (this.editedPhotoUrl && !this.photoUrlError && this.isPhotoValid)) &&
               (this.editedName !== this.initialName || 
                this.editedPhotoUrl !== this.initialPhotoUrl)
      }
    },
    methods: {
      validateName() {
        this.nameError = ''
        
        const trimmedName = this.editedName.trim()
        
        if (!trimmedName) {
          this.nameError = 'Il nome del gruppo è obbligatorio'
          return
        }
  
        if (trimmedName.length < 3) {
          this.nameError = 'Il nome deve contenere almeno 3 caratteri'
          return
        }
  
        if (trimmedName.length > 30) {
          this.nameError = 'Il nome non può superare i 30 caratteri'
        }
      },
      validatePhotoUrl() {
        this.photoUrlError = ''
        this.isPhotoValid = true
  
        if (!this.editedPhotoUrl) {
          // URL vuoto è consentito (opzionale)
          return
        }
  
        // Regex di base per validare gli URL
        const urlPattern = /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$/
        
        if (!urlPattern.test(this.editedPhotoUrl)) {
          this.photoUrlError = 'Inserisci un URL valido'
          this.isPhotoValid = false
        }
      },
      handleImageError() {
        // Gestisce l'errore di caricamento dell'immagine
        this.photoUrlError = 'Impossibile caricare l\'immagine'
        this.isPhotoValid = false
      },
      saveChanges() {
        if (!this.isFormValid) return
  
        const groupData = {
          name: this.editedName.trim(),
          photoUrl: this.editedPhotoUrl ? this.editedPhotoUrl.trim() : null
        }
  
        this.$emit('save', groupData)
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
  
  .photo-input-container {
    display: flex;
    align-items: center;
    gap: 12px;
  }
  
  .photo-preview {
    width: 60px;
    height: 60px;
    border-radius: 8px;
    overflow: hidden;
    border: 2px solid #2a3942;
  }
  
  .photo-preview img {
    width: 100%;
    height: 100%;
    object-fit: cover;
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
    margin-top: 8px;
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