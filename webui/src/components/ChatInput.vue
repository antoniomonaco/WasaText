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
  
      <!-- Input messaggi -->
      <div class="input-wrapper">
        <input 
          v-model="messageText" 
          type="text"
          placeholder="Scrivi un messaggio"
          @keyup.enter="sendMessage"
          class="message-input"
        />
        <button 
          class="send-button"
          :disabled="!messageText.trim()"
          @click="sendMessage"
        >
          <i class="fas fa-paper-plane"></i>
        </button>
      </div>
    </div>
  </template>
  
  <script>
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
        messageText: ''
      }
    },
    methods: {
      sendMessage() {
        if (!this.messageText.trim()) return
  
        const messagePayload = {
          text: this.messageText.trim(),
          replyTo: this.replyingTo || null
        }
  
        this.$emit('send-message', messagePayload)
        
        // Resetta l'input e la risposta
        this.messageText = ''
        this.$emit('cancel-reply')
      },
      cancelReply() {
        this.$emit('cancel-reply')
      },
      getReplyingSenderName() {
        return this.replyingTo?.sender?.username || 'Mittente sconosciuto'
      },
      getReplyMessageContent() {
        if (!this.replyingTo) return ''
        
        if (this.replyingTo.type === 'text') {
          return this.replyingTo.content
        }
        return 'Media'
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
  
  .reply-content {
    flex: 1;
    min-width: 0;
  }
  
  .reply-sender {
    color: #00a884;
    font-size: 13px;
    font-weight: 500;
    margin-bottom: 2px;
  }
  
  .reply-text {
    color: #d1d7db;
    font-size: 13px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .icon-button {
    background: none;
    border: none;
    color: #8696a0;
    padding: 4px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    width: 24px;
    height: 24px;
    transition: all 0.2s;
  }
  
  .icon-button:hover {
    background-color: #2a3942;
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
  
  .input-wrapper:focus-within {
    box-shadow: 0 0 0 2px #00a884;
  }
  
  .message-input {
    flex: 1;
    background: none;
    border: none;
    color: #d1d7db;
    font-size: 15px;
    outline: none;
    min-height: 20px;
    line-height: 20px;
    padding: 0;
  }
  
  .message-input::placeholder {
    color: #8696a0;
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