<!-- ChatInput.vue -->
<template>
  <div class="input-area">
    <!-- Area risposta -->
    <div 
      v-if="replyingTo" 
      class="reply-container"
    >
      <div class="reply-content">
        <div class="reply-sender">
          Risposta a {{ getReplyingSenderName() }}
        </div>
        <div class="reply-text">
          {{ getReplyMessageContent() }}
        </div>
      </div>
      <button 
        class="icon-button" 
        @click="cancelReply"
      >
        <i class="fas fa-times"></i>
      </button>
    </div>

    <!-- Anteprima immagine -->
    <div v-if="imagePreview" class="image-preview-container">
      <img :src="imagePreview" alt="Preview" class="image-preview"/>
      <button class="remove-image-button" @click="removeImage">
        <i class="fas fa-times"></i>
      </button>
      <div class="image-size-warning" v-if="showSizeWarning">
        Immagine troppo grande! Dimensione massima: 5MB
      </div>
    </div>

    <!-- Input messaggi -->
    <div class="input-wrapper">
      <!-- Menu immagini -->
      <div class="image-menu">
        <button 
          class="attachment-button"
          @click="$refs.fileInput.click()"
        >
          <i class="fas fa-image"></i>
        </button>
        <input 
          type="file"
          ref="fileInput"
          @change="handleFileUpload"
          accept="image/*"
          class="file-input"
        />
      </div>
      
      <!-- Input testo -->
      <input 
        v-model="messageText" 
        type="text"
        placeholder="Scrivi un messaggio"
        @keyup.enter="sendMessage"
        class="message-input"
        :disabled="isProcessingImage"
      />

      <!-- Pulsante invio -->
      <button 
        class="send-button"
        :disabled="!canSend"
        @click="sendMessage"
      >
        <i class="fas fa-paper-plane"></i>
      </button>
    </div>
  </div>
</template>

<script>
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB in bytes

export default {
  name: 'ChatInput',
  props: {
    replyingTo: {
      type: Object,
      default: null
    }
  },
  data() {
    return {
      messageText: '',
      imagePreview: null,
      imageBase64: null,
      isProcessingImage: false,
      showSizeWarning: false,
      showImageInput: false
    }
  },
  computed: {
    canSend() {
      return (this.messageText.trim() || this.imageBase64) && !this.isProcessingImage;
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
        const base64 = await this.convertToBase64(file);
        this.imageBase64 = base64;
        this.imagePreview = base64;
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

    async sendMessage() {
      if (!this.canSend) return;

      const messagePayload = {
        type: this.imageBase64 ? 'media' : 'text',
        content: this.imageBase64 || this.messageText.trim(),
        text: this.messageText.trim(),
        replyTo: this.replyingTo || null
      };

      this.$emit('send-message', messagePayload);
      
      // Reset dello stato
      this.messageText = '';
      this.imageBase64 = null;
      this.imagePreview = null;
      this.showImageInput = false;
      this.$emit('cancel-reply');
    },

    removeImage() {
      this.imageBase64 = null;
      this.imagePreview = null;
      this.showSizeWarning = false;
    },

    cancelReply() {
      this.$emit('cancel-reply');
    },

    getReplyingSenderName() {
      return this.replyingTo?.sender?.username || 'Mittente sconosciuto';
    },

    getReplyMessageContent() {
      if (!this.replyingTo) return '';
      
      if (this.replyingTo.type === 'text') {
        return this.replyingTo.content;
      }
      return 'Media';
    }
  }
}
</script>

<style scoped>
.input-area {
  background-color: #202c33;
  padding: 10px 16px;
  position: relative;
  z-index: 2;
  flex-shrink: 0;
  border-left: 1px solid #2a3942;
}

.reply-container {
  background-color: #1f2c34;
  padding: 8px 12px;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-left: 4px solid #00a884;
  border-radius: 4px;
}

.image-preview-container {
  position: relative;
  margin-bottom: 8px;
  display: inline-block;
  max-width: 200px;
}

.image-preview {
  max-width: 100%;
  max-height: 150px;
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
}

.image-size-warning {
  position: absolute;
  bottom: -20px;
  left: 0;
  right: 0;
  color: #f15c6d;
  font-size: 12px;
  text-align: center;
}

.input-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: #2a3942;
  border-radius: 8px;
  padding: 9px 12px;
  position: relative;
}

.image-menu {
  position: relative;
}

.file-input {
  display: none;
}

.attachment-button {
  background: none;
  border: none;
  color: #8696a0;
  padding: 12px;
  cursor: pointer;
  border-radius: 50%;
  transition: all 0.2s;
  width: 45px;
  height: 45px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 4px;
}

.attachment-button i {
  font-size: 20px;
}

.attachment-button:hover {
  color: #00a884;
  background-color: rgba(0, 168, 132, 0.1);
  transform: scale(1.1);
}

.message-input {
  flex: 1;
  background: none;
  border: none;
  color: #e9edef;
  font-size: 15px;
  outline: none;
  min-height: 20px;
  line-height: 20px;
  padding: 0;
}

.message-input::placeholder {
  color: #8696a0;
}

.message-input:disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.send-button {
  background: none;
  border: none;
  color: #00a884;
  font-size: 20px;
  padding: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.send-button:disabled {
  color: #8696a0;
  cursor: default;
}

.send-button:not(:disabled):hover {
  color: #00c29d;
}
</style>