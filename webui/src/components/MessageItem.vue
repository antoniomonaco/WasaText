<!-- MessageItem.vue -->
<template>
  <div 
    ref="messageWrapper"
    class="message-wrapper" 
    :class="{
      'sent': isSentByMe,
      'received': !isSentByMe,
      'media-message': message.type === 'media'
    }"
  >
    <div 
      class="message"
      @contextmenu.prevent="openMessageMenu"
    >
      <!-- Nome mittente (solo per gruppi) -->
      <div 
        v-if="isGroup && !isSentByMe" 
        class="message-sender"
      >
        {{ message.sender.username }}
      </div>

      <div class="message-content">
        <!-- Messaggio citato -->
        <div 
          v-if="replyMessage" 
          class="quoted-message"
        >
          <div class="quoted-sender">
            {{ replyMessage.sender?.username || 'Unknown User' }}
          </div>
          <div class="quoted-content">
            <template v-if="replyMessage.type === 'media'">
              <img 
                :src="replyMessage.content" 
                alt="Quoted media"
                class="quoted-media"
              />
            </template>
            <template v-else>
              {{ replyMessage.content }}
            </template>
          </div>
        </div>

        <!-- Contenuto del messaggio -->
        <div 
          v-if="message.type === 'text'" 
          class="message-text"
        >
          {{ message.content }}
        </div>

        <!-- Media -->
        <div v-else-if="message.type === 'media'" class="media-container">
          <img 
            :src="message.content"
            :alt="'Image shared by ' + message.sender.username"
            class="message-media"
            @load="handleImageLoad"
            @error="handleImageError"
            @click="showFullImage"
          />
          <div v-if="imageError" class="image-error">
            Impossibile caricare l'immagine
          </div>
        </div>

        <!-- Badge commenti -->
        <div 
          v-if="hasComments" 
          class="comments-badge"
          @click="$emit('show-comments', message)"
        >
          <i class="fas fa-comment"></i>
          {{ message.comments.length }}
        </div>

        <!-- Metadata messaggio -->
        <div class="message-meta">
          <span class="message-time">
            {{ formatTime(message.timestamp) }}
          </span>
          <span 
            v-if="isSentByMe" 
            class="message-status"
          >
            <i 
              :class="[
                message.status === 'read' 
                  ? 'fas fa-check-double read' 
                  : 'fas fa-check'
              ]"
            ></i>
          </span>
        </div>
      </div>
    </div>

    <!-- Menu contestuale del messaggio -->
    <div 
      v-if="isMessageMenuOpen" 
      ref="messageMenu"
      class="message-menu"
      :class="{
        'menu-sent': isSentByMe,
        'menu-received': !isSentByMe
      }"
    >
      <div 
        class="menu-item" 
        @click="$emit('reply-message', message)"
      >
        <i class="fas fa-reply"></i>
        Rispondi
      </div>
      <div 
        class="menu-item" 
        @click="$emit('comment-message', message)"
      >
        <i class="fas fa-comment"></i>
        Commenta
      </div>
      <div 
        class="menu-item" 
        @click="$emit('show-comments', message)"
      >
        <i class="fas fa-comments"></i>
        Visualizza commenti
      </div>
      <div 
        class="menu-item" 
        @click="$emit('forward-message', message)"
      >
        <i class="fas fa-share"></i>
        Inoltra
      </div>
      <div 
        v-if="isSentByMe" 
        class="menu-item delete"
        @click="$emit('delete-message', message)"
      >
        <i class="fas fa-trash"></i>
        Elimina
      </div>
    </div>

    <!-- Modal immagine a schermo intero -->
    <div v-if="showFullImageModal" class="full-image-modal" @click="closeFullImage">
      <div class="full-image-container">
        <img :src="message.content" :alt="'Image shared by ' + message.sender.username" />
        <button class="close-button" @click="closeFullImage">
          <i class="fas fa-times"></i>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MessageItem',
  props: {
    message: {
      type: Object,
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
    isMessageMenuOpen: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      showFullImageModal: false,
      imageError: false,
      replyMessage: null
    }
  },
  computed: {
    isSentByMe() {
      return this.message.sender.id === this.currentUserId
    },
    hasComments() {
      return this.message.comments && this.message.comments.length > 0
    }
  },
  methods: {
    formatTime(timestamp) {
      if (!timestamp) return ''
      const date = new Date(timestamp)
      return date.toLocaleTimeString([], { 
        hour: '2-digit', 
        minute: '2-digit' ,
        timezone : 'UTC'
      })
    },
    openMessageMenu(event) {
      this.$emit('open-message-menu', this.message, event)
    },
    showFullImage() {
      this.showFullImageModal = true
      document.body.style.overflow = 'hidden'
    },
    closeFullImage() {
      this.showFullImageModal = false
      document.body.style.overflow = 'auto'
    },
    handleImageLoad(event) {
      this.imageError = false;
      console.log('Immagine caricata:', this.message);
      this.$emit('media-loaded');
    },
    handleImageError(event) {
      this.imageError = true;
      console.error('Errore nel caricamento dell\'immagine:', this.message.content);
    },
    async fetchReplyMessage() {
      if (this.message.replyTo) {
        try {
          const response = await this.$axios.get(
            `/conversations/${this.$parent.conversationID}/messages/${this.message.replyTo}`,
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('authToken')}`
              }
            }
          );
          this.replyMessage = response.data;
        } catch (error) {
          console.error('Error fetching reply message:', error);
        }
      }
    }
  },
  async created() {
    await this.fetchReplyMessage();
  }
}
</script>

<style scoped>
.message-wrapper {
  display: flex;
  position: relative;
  max-width: 65%;
  margin: 2px 0;
}

.message-wrapper.sent {
  margin-left: auto;
}

.message-wrapper.media-message {
  max-width: 85%;
}

.message {
  padding: 8px 10px;
  border-radius: 7.5px;
  position: relative;
  min-width: 100px;
  background-color: var(--bg-message);
  color: #e9edef;
}

.sent .message {
  --bg-message: #005c4b;
  margin-right: 4px;
}

.received .message {
  --bg-message: #202c33;
  margin-left: 4px;
}

.message-sender {
  color: #53bdeb;
  font-size: 13px;
  margin-bottom: 2px;
}

.message-content {
  position: relative;
}

.message-text {
  word-wrap: break-word;
  white-space: pre-wrap;
  font-size: 14.2px;
  line-height: 19px;
}

.media-container {
  position: relative;
  border-radius: 6px;
  overflow: hidden;
  cursor: pointer;
}

.message-media {
  max-width: 100%;
  max-height: 300px;
  display: block;
  object-fit: contain;
  border-radius: 6px;
  transition: transform 0.2s ease;
}

.message-media:hover {
  transform: scale(0.99);
}

.quoted-message {
  background-color: #ffffff1a;
  border-left: 4px solid #53bdeb;
  border-radius: 4px;
  margin-bottom: 8px;
  padding: 5px 10px;
}

.quoted-sender {
  color: #53bdeb;
  font-size: 13px;
  font-weight: 500;
}

.quoted-content {
  color: #e9edef;
  font-size: 13px;
  opacity: 0.7;
  overflow: hidden;
  text-overflow: ellipsis;
}

.quoted-media {
  max-width: 50px;
  max-height: 50px;
  border-radius: 4px;
  margin-top: 4px;
}

.comments-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: #8696a0;
  font-size: 12px;
  margin-top: 4px;
  cursor: pointer;
  padding: 2px 4px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.comments-badge:hover {
  background-color: rgba(134, 150, 160, 0.1);
}

.message-meta {
  display: flex;
  align-items: center;
  gap: 4px;
  float: right;
  margin-left: 8px;
  margin-top: 2px;
}

.message-time {
  color: #8696a0;
  font-size: 11px;
}

.message-status {
  color: #8696a0;
  font-size: 11px;
}

.message-status .read {
  color: #53bdeb;
}

/* Stili per i bordi dei messaggi */
.sent .message::before {
  content: "";
  position: absolute;
  top: 0;
  right: -8px;
  width: 8px;
  height: 13px;
  background: linear-gradient(
    to top left,
    transparent 50%,
    #005c4b 50%
  );
}

.received .message::before {
  content: "";
  position: absolute;
  top: 0;
  left: -8px;
  width: 8px;
  height: 13px;
  background: linear-gradient(
    to top right,
    transparent 50%,
    #202c33 50%
  );
}

.message-menu {
  position: absolute;
  background-color: #233138;
  border-radius: 6px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.4);
  z-index: 100;
  min-width: 200px;
  padding: 4px 0;
  top: 100%;
  margin-top: 8px;
}

.message-menu.menu-sent {
  right: 0;
}

.message-menu.menu-received {
  left: 0;
}

.menu-item {
  padding: 10px 16px;
  color: #e9edef;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
  transition: background-color 0.2s;
}

.menu-item:hover {
  background-color: #182229;
}

.menu-item i {
  width: 20px;
  text-align: center;
  font-size: 16px;
}

.menu-item.delete {
  color: #f15c6d;
}

.menu-item.delete i {
  color: #f15c6d;
}

/* Modal immagine a schermo intero */
.full-image-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.full-image-container {
  position: relative;
  max-width: 90vw;
  max-height: 90vh;
}

.full-image-container img {
  max-width: 100%;
  max-height: 90vh;
  object-fit: contain;
}

.close-button {
  position: absolute;
  top: -40px;
  right: 0;
  background: none;
  border: none;
  color: #fff;
  font-size: 24px;
  cursor: pointer;
  padding: 8px;
  border-radius: 50%;
  transition: background-color 0.2s;
}

.close-button:hover {
  background-color: rgba(255, 255, 255, 0.1);
}
</style>