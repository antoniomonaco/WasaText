<!-- MessageItem.vue -->
<template>
  <div 
    ref="messageWrapper"
    class="message-wrapper" 
    :class="{
      'sent': isSentByMe,
      'received': !isSentByMe
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
        <!-- Messaggio citato (se presente) -->
        <div 
          v-if="message.replyTo" 
          class="quoted-message"
        >
          <div class="quoted-sender">
            {{ getReplyingSenderName(message.replyTo) }}
          </div>
          <div class="quoted-content">
            {{ getReplyMessageContent(message.replyTo) }}
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
        <img 
          v-else-if="message.type === 'media'"
          :src="message.content"
          :alt="message.type"
          class="message-media"
          @load="$emit('media-loaded')"
        />

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
        minute: '2-digit' 
      })
    },
    getReplyingSenderName(replyTo) {
      return replyTo.sender?.username || 'Mittente sconosciuto'
    },
    getReplyMessageContent(replyTo) {
      if (replyTo.type === 'text') {
        return replyTo.content
      }
      return 'Media'
    },
    openMessageMenu(event) {
      // Emetti un evento per aprire/chiudere il menu del messaggio
      this.$emit('open-message-menu', this.message, event)
    }
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

.message-media {
  max-width: 100%;
  border-radius: 4px;
  cursor: pointer;
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
  white-space: nowrap;
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
</style>