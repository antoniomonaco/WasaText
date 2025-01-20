<!-- CommentsModal.vue -->
<template>
  <div class="comments-modal">
    <!-- Loading state -->
    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <span>Caricamento reazioni...</span>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="error-state">
      {{ error }}
    </div>

    <!-- Comments list -->
    <div class="comments-list" :class="{ 'scrollable': comments && comments.length > 4 }">
      <div v-if="!comments || comments.length === 0" class="no-comments">
        Nessuna reazione. Sii il primo a reagire!
      </div>
      
      <div v-else v-for="comment in comments" :key="comment.id" class="comment">
        <div class="comment-content-wrapper">
          <div class="comment-header">
            <span class="comment-sender">
              {{ comment.sender.username }}
            </span>
            <span class="comment-time">
              {{ formatTime(comment.timestamp) }}
            </span>
          </div>
          
          <div class="comment-content emoji">
            {{ comment.content }}
          </div>

          <button 
            v-if="canDelete(comment)"
            class="delete-comment-button"
            @click="$emit('delete-comment', comment)"
          >
            <i class="fas fa-times"></i>
          </button>
        </div>
      </div>
    </div>

    <!-- Emoji Picker -->
    <div class="reaction-container">
      <div 
        v-if="showEmojiPicker" 
        class="emoji-picker-wrapper"
      >
        <div class="emoji-picker-overlay" @click="showEmojiPicker = false"></div>
        <div class="emoji-picker-content">
          <EmojiPicker @select="selectEmoji" />
        </div>
      </div>

      <div class="new-reaction">
        <button 
          class="reaction-button"
          @click="toggleEmojiPicker"
        >
          <i class="far fa-smile"></i>
          Aggiungi reazione
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import EmojiPicker from './EmojiPicker.vue'

export default {
  name: 'CommentsModal',
  components: {
    EmojiPicker
  },
  props: {
    comments: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    },
    error: {
      type: String,
      default: ''
    },
    currentUserId: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      showEmojiPicker: false
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
    canDelete(comment) {
      return comment.sender.id === this.currentUserId
    },
    selectEmoji(emoji) {
      this.$emit('submit', emoji)
      this.showEmojiPicker = false
    },
    toggleEmojiPicker() {
      this.showEmojiPicker = !this.showEmojiPicker;
    }
  }
}
</script>

<style scoped>
.comments-modal {
  display: flex;
  flex-direction: column;
  height: 100%;
  max-height: 500px;
}

.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  color: #8696a0;
}

.loading-spinner {
  width: 30px;
  height: 30px;
  border: 3px solid #2a3942;
  border-top-color: #00a884;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 10px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.comments-list {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.comments-list.scrollable {
  overflow-y: auto;
}

.no-comments {
  text-align: center;
  color: #8696a0;
  padding: 20px;
  background-color: rgba(134, 150, 160, 0.05);
  border-radius: 6px;
}

.comment {
  margin-bottom: 12px;
  background-color: #202c33;
  border-radius: 8px;
  position: relative;
}

.comment-content-wrapper {
  padding: 12px;
  position: relative;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.comment-sender {
  color: #00a884;
  font-size: 14px;
  font-weight: 500;
}

.comment-time {
  color: #8696a0;
  font-size: 12px;
}

.comment-content {
  color: #e9edef;
  font-size: 14px;
  word-break: break-word;
}

.comment-content.emoji {
  font-size: 24px;
  text-align: center;
  padding: 8px 0;
}

.delete-comment-button {
  position: absolute;
  top: 8px;
  right: 8px;
  background: none;
  border: none;
  color: #f15c6d;
  cursor: pointer;
  padding: 4px;
  border-radius: 50%;
  opacity: 0;
  transition: all 0.2s;
}

.comment:hover .delete-comment-button {
  opacity: 1;
}

.delete-comment-button:hover {
  background-color: #2a3942;
}

.emoji-picker-container {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  margin-bottom: 8px;
  z-index: 1000;
}

.reaction-container {
  position: relative;
  width: 100%;
}

.emoji-picker-wrapper {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.emoji-picker-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
}

.emoji-picker-content {
  position: relative;
  z-index: 1001;
  background-color: #1f2c34;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
}

.new-reaction {
  padding: 16px;
  border-top: 1px solid #2a3942;
  position: relative;
  z-index: 1;
}

.reaction-button {
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
  font-size: 15px;
  transition: background-color 0.2s;
}

.reaction-button:hover {
  background-color: #374045;
}

.reaction-button i {
  font-size: 18px;
}
</style>